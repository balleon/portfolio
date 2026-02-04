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
		return nil, fmt.Errorf("Building config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Creating clientset: %w", err)
	}

	return clientset, nil
}

// listAllSecrets returns all Secrets across all Namespaces.
func listAllSecrets(clientset *kubernetes.Clientset) ([]SecretRef, error) {
	secrets, err := clientset.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Listing Secrets: %w", err)
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

	allDeployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Listing Deployments: %w", err)
	}
	for _, deployment := range allDeployments.Items {
		extractSecretsFromPodSpec(&deployment.Spec.Template.Spec, deployment.Namespace, usedSecrets)
	}

	allStatefulSets, err := clientset.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Listing StatefulSets: %w", err)
	}
	for _, statefulset := range allStatefulSets.Items {
		extractSecretsFromPodSpec(&statefulset.Spec.Template.Spec, statefulset.Namespace, usedSecrets)
	}

	allDaemonSets, err := clientset.AppsV1().DaemonSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Listing DaemonSets: %w", err)
	}
	for _, daemonset := range allDaemonSets.Items {
		extractSecretsFromPodSpec(&daemonset.Spec.Template.Spec, daemonset.Namespace, usedSecrets)
	}

	allJobs, err := clientset.BatchV1().Jobs("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Listing Jobs: %w", err)
	}
	for _, job := range allJobs.Items {
		extractSecretsFromPodSpec(&job.Spec.Template.Spec, job.Namespace, usedSecrets)
	}

	allCronJobs, err := clientset.BatchV1().CronJobs("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Listing CronJobs: %w", err)
	}
	for _, cronjob := range allCronJobs.Items {
		extractSecretsFromPodSpec(&cronjob.Spec.JobTemplate.Spec.Template.Spec, cronjob.Namespace, usedSecrets)
	}

	allIngresses, err := clientset.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Listing Ingresses: %w", err)
	}
	for _, ingress := range allIngresses.Items {
		for _, tls := range ingress.Spec.TLS {
			if tls.SecretName != "" {
				usedSecrets[SecretRef{Name: tls.SecretName, Namespace: ingress.Namespace}] = true
			}
		}
	}

	return usedSecrets, nil
}

// extractSecretsFromPodSpec extracts all Secret references from a PodSpec.
func extractSecretsFromPodSpec(spec *corev1.PodSpec, namespace string, usedSecrets map[SecretRef]bool) {
	for _, imagepullsecret := range spec.ImagePullSecrets {
		usedSecrets[SecretRef{Name: imagepullsecret.Name, Namespace: namespace}] = true
	}

	for _, volume := range spec.Volumes {
		if volume.Secret != nil {
			usedSecrets[SecretRef{Name: volume.Secret.SecretName, Namespace: namespace}] = true
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
