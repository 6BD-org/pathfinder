package client

import (
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	client, err := client.New(config, client.Options{})
	if err != nil {
		return nil, err
	}
	return XMClientImpl{
		client: client,
	}, nil
}
