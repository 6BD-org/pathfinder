package client

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/wylswz/native-discovery/api/v1"
)

type PathFinderV1 interface {
	Create(ctx context.Context, pathfinder *v1.PathFinder, opts client.CreateOption) error
	Delete(ctx context.Context, pathfinder *v1.PathFinder, opts client.DeleteOption) error
}

type PathFinderV1Impl struct {
	client    client.Client
	namespace string
}

// Create create a new pathfinder
func (pfv1 PathFinderV1Impl) Create(ctx context.Context, pathfinder *v1.PathFinder, opts client.CreateOption) error {
	return pfv1.client.Create(ctx, pathfinder, opts)
}

// Delete a pathfinder from namespace
func (pfv1 PathFinderV1Impl) Delete(ctx context.Context, pathfinder *v1.PathFinder, opts client.DeleteOption) error {
	return pfv1.client.Delete(ctx, pathfinder, opts)

}

func NewPathFinderV1(client client.Client, namespace string) PathFinderV1 {
	return PathFinderV1Impl{
		client:    client,
		namespace: namespace,
	}
}
