package client

import (
	"context"
	"log"
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/wylswz/native-discovery/api/v1"
	"github.com/wylswz/native-discovery/utils"
)

// PathFinderListOption inherates native options and region
type PathFinderListOption struct {
	client.ListOption
	Region string
}

type PathFinderV1 interface {
	Create(ctx context.Context, pathfinder *v1.PathFinder, opts client.CreateOption) error
	Delete(ctx context.Context, pathfinder *v1.PathFinder, opts client.DeleteOption) error
	List(ctx context.Context, pathfinderList *v1.PathFinderList, opts PathFinderListOption) error
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

// List path finder objects
// Besides native list options, you are also able to
// filter using regions
func (pfv1 PathFinderV1Impl) List(ctx context.Context, pathfinderList *v1.PathFinderList, opts PathFinderListOption) error {
	if opts.ListOption == nil {
		opts.ListOption = &client.ListOptions{}
	}
	err := pfv1.client.List(ctx, pathfinderList, opts)
	if err != nil {
		return err
	}
	filted := utils.Filter(
		pathfinderList.Items,
		func(l interface{}) bool { return l.(v1.PathFinder).Spec.Region == opts.Region },
		reflect.TypeOf(v1.PathFinder{}),
	)
	log.Println(filted)
	pathfinderList.Items = make([]v1.PathFinder, len(filted))
	for i := 0; i < len(filted); i++ {
		pathfinderList.Items[i] = filted[i].(v1.PathFinder)
	}
	return err
}

// NewPathFinderV1 Create a new pathfinder v1 api
func NewPathFinderV1(client client.Client, namespace string) PathFinderV1 {
	return PathFinderV1Impl{
		client:    client,
		namespace: namespace,
	}
}
