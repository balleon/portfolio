package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	// Kubernetes client configuration
	clientset := kubernetesClientAuth()

	// List Namespaces
	namespaces := listNamespaces(clientset)

	// Extract Secrets
	secretList := listSecrets(clientset, namespaces)

	// Identify used Secrets (imagePullSecrets, Volumes, Containers and initContainers)
	for _, i := range secretList {
		if !secretDeployments(clientset, i) && !secretStateFulsets(clientset, i) && !secretDaemonSets(clientset, i) {
			fmt.Printf("Secret %s is unused.\n", i)
		}
	}
}

func kubernetesClientAuth() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return clientset
}

func listNamespaces(clientset *kubernetes.Clientset) []string {
	namespaceClient, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	namespaceList := []string{}
	for i := range len(namespaceClient.Items) {
		namespaceList = append(namespaceList, namespaceClient.Items[i].Name)
	}

	return namespaceList
}

func listSecrets(clientset *kubernetes.Clientset, namespaces []string) []string {
	secretList := []string{}

	for i := range len(namespaces) {
		secretClient, err := clientset.CoreV1().Secrets(namespaces[i]).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for y := range len(secretClient.Items) {
			secretList = append(secretList, secretClient.Items[y].Name)
		}
	}

	return secretList
}

func secretDeployments(clientset *kubernetes.Clientset, secretName string) bool {
	var isUsed bool

	// Extract Deployments
	deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for i := range deployments.Items {
		// imagePullSecrets
		imagePullSecrets := false
		for j := range deployments.Items[i].Spec.Template.Spec.ImagePullSecrets {
			if deployments.Items[i].Spec.Template.Spec.ImagePullSecrets[j].Name == secretName {
				imagePullSecrets = true
			}
		}
		// Volumes
		Volumes := false
		for j := range deployments.Items[i].Spec.Template.Spec.Volumes {
			if deployments.Items[i].Spec.Template.Spec.Volumes[j].Secret != nil {
				if deployments.Items[i].Spec.Template.Spec.Volumes[j].Secret.SecretName == secretName {
					Volumes = true
				}
			}
		}
		// Containers using EnvFrom
		EnvFrom := false
		for j := range deployments.Items[i].Spec.Template.Spec.Containers {
			for k := range deployments.Items[i].Spec.Template.Spec.Containers[j].EnvFrom {
				if deployments.Items[i].Spec.Template.Spec.Containers[j].EnvFrom[k].SecretRef.Name == secretName {
					EnvFrom = true
				}
			}
		}
		// Containers using Env
		Env := false
		for j := range deployments.Items[i].Spec.Template.Spec.Containers {
			for k := range deployments.Items[i].Spec.Template.Spec.Containers[j].Env {
				if deployments.Items[i].Spec.Template.Spec.Containers[j].Env[k].ValueFrom.SecretKeyRef != nil {
					if deployments.Items[i].Spec.Template.Spec.Containers[j].Env[k].ValueFrom.SecretKeyRef.Name == secretName {
						Env = true
					}
				}
			}
		}
		// initContainers using EnvFrom
		initEnvFrom := false
		for j := range deployments.Items[i].Spec.Template.Spec.InitContainers {
			for k := range deployments.Items[i].Spec.Template.Spec.InitContainers[j].EnvFrom {
				if deployments.Items[i].Spec.Template.Spec.InitContainers[j].EnvFrom[k].SecretRef.Name == secretName {
					initEnvFrom = true
				}
			}
		}
		// initContainers using Env
		initEnv := false
		for j := range deployments.Items[i].Spec.Template.Spec.InitContainers {
			for k := range deployments.Items[i].Spec.Template.Spec.InitContainers[j].Env {
				if deployments.Items[i].Spec.Template.Spec.InitContainers[j].Env[k].ValueFrom.SecretKeyRef != nil {
					if deployments.Items[i].Spec.Template.Spec.InitContainers[j].Env[k].ValueFrom.SecretKeyRef.Name == secretName {
						initEnv = true
					}
				}
			}
		}

		// Determine Secret usage
		if imagePullSecrets || Volumes || EnvFrom || Env || initEnvFrom || initEnv {
			isUsed = true
		} else {
			isUsed = false
		}
	}

	return isUsed
}

func secretStateFulsets(clientset *kubernetes.Clientset, secretName string) bool {
	var isUsed bool

	// Extract StatefulSet
	statefulsets, err := clientset.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for i := range statefulsets.Items {
		// imagePullSecrets
		imagePullSecrets := false
		for j := range statefulsets.Items[i].Spec.Template.Spec.ImagePullSecrets {
			if statefulsets.Items[i].Spec.Template.Spec.ImagePullSecrets[j].Name == secretName {
				imagePullSecrets = true
			}
		}
		// Volumes
		Volumes := false
		for j := range statefulsets.Items[i].Spec.Template.Spec.Volumes {
			if statefulsets.Items[i].Spec.Template.Spec.Volumes[j].Secret != nil {
				if statefulsets.Items[i].Spec.Template.Spec.Volumes[j].Secret.SecretName == secretName {
					Volumes = true
				}
			}
		}
		// Containers using EnvFrom
		EnvFrom := false
		for j := range statefulsets.Items[i].Spec.Template.Spec.Containers {
			for k := range statefulsets.Items[i].Spec.Template.Spec.Containers[j].EnvFrom {
				if statefulsets.Items[i].Spec.Template.Spec.Containers[j].EnvFrom[k].SecretRef.Name == secretName {
					EnvFrom = true
				}
			}
		}
		// Containers using Env
		Env := false
		for j := range statefulsets.Items[i].Spec.Template.Spec.Containers {
			for k := range statefulsets.Items[i].Spec.Template.Spec.Containers[j].Env {
				if statefulsets.Items[i].Spec.Template.Spec.Containers[j].Env[k].ValueFrom.SecretKeyRef != nil {
					if statefulsets.Items[i].Spec.Template.Spec.Containers[j].Env[k].ValueFrom.SecretKeyRef.Name == secretName {
						Env = true
					}
				}
			}
		}
		// initContainers using EnvFrom
		initEnvFrom := false
		for j := range statefulsets.Items[i].Spec.Template.Spec.InitContainers {
			for k := range statefulsets.Items[i].Spec.Template.Spec.InitContainers[j].EnvFrom {
				if statefulsets.Items[i].Spec.Template.Spec.InitContainers[j].EnvFrom[k].SecretRef.Name == secretName {
					initEnvFrom = true
				}
			}
		}
		// initContainers using Env
		initEnv := false
		for j := range statefulsets.Items[i].Spec.Template.Spec.InitContainers {
			for k := range statefulsets.Items[i].Spec.Template.Spec.InitContainers[j].Env {
				if statefulsets.Items[i].Spec.Template.Spec.InitContainers[j].Env[k].ValueFrom.SecretKeyRef != nil {
					if statefulsets.Items[i].Spec.Template.Spec.InitContainers[j].Env[k].ValueFrom.SecretKeyRef.Name == secretName {
						initEnv = true
					}
				}
			}
		}

		// Determine Secret usage
		if imagePullSecrets || Volumes || EnvFrom || Env || initEnvFrom || initEnv {
			isUsed = true
		} else {
			isUsed = false
		}
	}

	return isUsed
}

func secretDaemonSets(clientset *kubernetes.Clientset, secretName string) bool {
	var isUsed bool

	// Extract DaemonSets
	daemonsets, err := clientset.AppsV1().DaemonSets("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for i := range daemonsets.Items {
		// imagePullSecrets
		imagePullSecrets := false
		for j := range daemonsets.Items[i].Spec.Template.Spec.ImagePullSecrets {
			if daemonsets.Items[i].Spec.Template.Spec.ImagePullSecrets[j].Name == secretName {
				imagePullSecrets = true
			}
		}
		// Volumes
		Volumes := false
		for j := range daemonsets.Items[i].Spec.Template.Spec.Volumes {
			if daemonsets.Items[i].Spec.Template.Spec.Volumes[j].Secret != nil {
				if daemonsets.Items[i].Spec.Template.Spec.Volumes[j].Secret.SecretName == secretName {
					Volumes = true
				}
			}
		}
		// Containers using EnvFrom
		EnvFrom := false
		for j := range daemonsets.Items[i].Spec.Template.Spec.Containers {
			for k := range daemonsets.Items[i].Spec.Template.Spec.Containers[j].EnvFrom {
				if daemonsets.Items[i].Spec.Template.Spec.Containers[j].EnvFrom[k].SecretRef.Name == secretName {
					EnvFrom = true
				}
			}
		}
		// Containers using Env
		Env := false
		for j := range daemonsets.Items[i].Spec.Template.Spec.Containers {
			for k := range daemonsets.Items[i].Spec.Template.Spec.Containers[j].Env {
				if daemonsets.Items[i].Spec.Template.Spec.Containers[j].Env[k].ValueFrom != nil && daemonsets.Items[i].Spec.Template.Spec.Containers[j].Env[k].ValueFrom.SecretKeyRef != nil {
					if daemonsets.Items[i].Spec.Template.Spec.Containers[j].Env[k].ValueFrom.SecretKeyRef.Name == secretName {
						Env = true
					}
				}
			}
		}
		// initContainers using EnvFrom
		initEnvFrom := false
		for j := range daemonsets.Items[i].Spec.Template.Spec.InitContainers {
			for k := range daemonsets.Items[i].Spec.Template.Spec.InitContainers[j].EnvFrom {
				if daemonsets.Items[i].Spec.Template.Spec.InitContainers[j].EnvFrom[k].SecretRef.Name == secretName {
					initEnvFrom = true
				}
			}
		}
		// initContainers using Env
		initEnv := false
		for j := range daemonsets.Items[i].Spec.Template.Spec.InitContainers {
			for k := range daemonsets.Items[i].Spec.Template.Spec.InitContainers[j].Env {
				if daemonsets.Items[i].Spec.Template.Spec.InitContainers[j].Env[k].ValueFrom.SecretKeyRef != nil {
					if daemonsets.Items[i].Spec.Template.Spec.InitContainers[j].Env[k].ValueFrom.SecretKeyRef.Name == secretName {
						initEnv = true
					}
				}
			}
		}

		// Determine Secret usage
		if imagePullSecrets || Volumes || EnvFrom || Env || initEnvFrom || initEnv {
			isUsed = true
		} else {
			isUsed = false
		}
	}

	return isUsed
}
