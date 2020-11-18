package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ServicesDiscoverer struct {
	Client *kubernetes.Clientset
}

// ListServices List all services in a namespace
func (s *ServicesDiscoverer) ListServices(namespace string) ([]v1.Service, error) {
	svcLst, err := s.Client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	return svcLst.Items, err
}

// NewDiscoverer Create a new service discoverer using kubeconfig
func NewDiscoverer(kubeconfig string) *ServicesDiscoverer {
	var clientset *kubernetes.Clientset
	config, err := rest.InClusterConfig()
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	s := ServicesDiscoverer{
		Client: clientset,
	}
	return &s
}
