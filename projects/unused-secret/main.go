package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// SecretRef uniquely identifies a Secret by Namespace and Name.
type SecretRef struct {
	Name      string
	Namespace string
}

// String return the Secret reference in "Namespace/Name" format.
func (s SecretRef) String() string {
	return fmt.Sprintf("%s/%s", s.Namespace, s.Name)
}

func main() {
	clientset, err := newKubernetesClient()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	secrets, err := listAllSecrets(clientset)
	if err != nil {
		log.Fatalf("Failed to list Secrets: %v", err)
	}

	usedSecrets, err := findUsedSecrets(clientset)
	if err != nil {
		log.Fatalf("Failed to find used Secrets: %v", err)
	}

	for _, secret := range secrets {
		if !usedSecrets[secret] {
			fmt.Printf("Secret %s is unused.\n", secret)
		}
	}
}

// newKubernetesClient creates a Kubernetes clientset from kubeconfig file.
func newKubernetesClient() (*kubernetes.Clientset, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = *flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = *flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("building config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating clientset: %w", err)
	}

	return clientset, nil
}

// listAllSecrets returns all Secrets across all Namespaces.
func listAllSecrets(clientset *kubernetes.Clientset) ([]SecretRef, error) {
	secrets, err := clientset.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing secrets: %w", err)
	}

	result := make([]SecretRef, 0, len(secrets.Items))
	for _, secret := range secrets.Items {
		result = append(result, SecretRef{
			Name:      secret.Name,
			Namespace: secret.Namespace,
		})
	}

	return result, nil
}

// findUsedSecrets scans all workloads and returns Secrets that are referenced.
func findUsedSecrets(clientset *kubernetes.Clientset) (map[SecretRef]bool, error) {
	usedSecrets := make(map[SecretRef]bool)

	deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing deployments: %w", err)
	}
	for _, deploy := range deployments.Items {
		extractSecretsFromPodSpec(&deploy.Spec.Template.Spec, deploy.Namespace, usedSecrets)
	}

	statefulsets, err := clientset.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing statefulsets: %w", err)
	}
	for _, sts := range statefulsets.Items {
		extractSecretsFromPodSpec(&sts.Spec.Template.Spec, sts.Namespace, usedSecrets)
	}

	daemonsets, err := clientset.AppsV1().DaemonSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing daemonsets: %w", err)
	}
	for _, ds := range daemonsets.Items {
		extractSecretsFromPodSpec(&ds.Spec.Template.Spec, ds.Namespace, usedSecrets)
	}

	jobs, err := clientset.BatchV1().Jobs("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing jobs: %w", err)
	}
	for _, job := range jobs.Items {
		extractSecretsFromPodSpec(&job.Spec.Template.Spec, job.Namespace, usedSecrets)
	}

	cronjobs, err := clientset.BatchV1().CronJobs("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing cronjobs: %w", err)
	}
	for _, cj := range cronjobs.Items {
		extractSecretsFromPodSpec(&cj.Spec.JobTemplate.Spec.Template.Spec, cj.Namespace, usedSecrets)
	}

	ingresses, err := clientset.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing ingresses: %w", err)
	}
	for _, ing := range ingresses.Items {
		for _, tls := range ing.Spec.TLS {
			if tls.SecretName != "" {
				usedSecrets[SecretRef{Name: tls.SecretName, Namespace: ing.Namespace}] = true
			}
		}
	}

	return usedSecrets, nil
}

// extractSecretsFromPodSpec extracts all Secret references from a PodSpec.
func extractSecretsFromPodSpec(spec *corev1.PodSpec, namespace string, usedSecrets map[SecretRef]bool) {
	for _, ips := range spec.ImagePullSecrets {
		usedSecrets[SecretRef{Name: ips.Name, Namespace: namespace}] = true
	}

	for _, vol := range spec.Volumes {
		if vol.Secret != nil {
			usedSecrets[SecretRef{Name: vol.Secret.SecretName, Namespace: namespace}] = true
		}
	}

	allContainers := append(spec.Containers, spec.InitContainers...)
	for _, container := range allContainers {
		extractSecretsFromContainer(&container, namespace, usedSecrets)
	}
}

// extractSecretsFromContainer extracts Secret references from a container's environment configuration.
func extractSecretsFromContainer(container *corev1.Container, namespace string, usedSecrets map[SecretRef]bool) {
	for _, envFrom := range container.EnvFrom {
		if envFrom.SecretRef != nil {
			usedSecrets[SecretRef{Name: envFrom.SecretRef.Name, Namespace: namespace}] = true
		}
	}

	for _, env := range container.Env {
		if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil {
			usedSecrets[SecretRef{Name: env.ValueFrom.SecretKeyRef.Name, Namespace: namespace}] = true
		}
	}
}
