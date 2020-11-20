package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ClientSet is a clientset to manipulate with built-in resources
func ClientSet(config *rest.Config) *kubernetes.Clientset {
	var clientset *kubernetes.Clientset
	clientset = kubernetes.NewForConfigOrDie(config)
	return clientset
}
