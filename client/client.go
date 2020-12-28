package client

import (
	pathfinderv1 "github.com/wylswz/native-discovery/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

type XMClient interface {
	PathFinderV1(namespace string) PathFinderV1
}

type XMClientImpl struct {
	client client.Client
}

func (cl XMClientImpl) PathFinderV1(namespace string) PathFinderV1 {
	return NewPathFinderV1(cl.client, namespace)
}

func New(config *rest.Config) (XMClient, error) {
	client, err := client.New(config, client.Options{
		Scheme: scheme,
	})

	// +kubebuilder:scaffold:scheme
	if err != nil {
		return nil, err
	}
	return XMClientImpl{
		client: client,
	}, nil
}

func init() {

	_ = clientgoscheme.AddToScheme(scheme)
	_ = pathfinderv1.AddToScheme(scheme)
}
