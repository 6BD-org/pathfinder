package client

import (
	"context"
	"fmt"
	"reflect"

	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/6BD-org/pathfinder/api/v1"
	"github.com/6BD-org/pathfinder/utils"
)

// PathFinderListOption inherates native options and region
type PathFinderListOption struct {
	client.ListOptions
	Region string
}

// PathFinderKey uses namespace + region as a key
type PathFinderKey struct {
	Namespace string
	Region    string
}

// PathFinderV1 is api interface for pathfinder v1
type PathFinderV1 interface {
	Create(ctx context.Context, pathfinder *v1.PathFinder, opts ...client.CreateOption) error
	Delete(ctx context.Context, pathfinder *v1.PathFinder, opts ...client.DeleteOption) error
	Get(ctx context.Context, name string, pathfinder *v1.PathFinder) error
	GetByRegion(ctx context.Context, region string, pathfinder *v1.PathFinder) error
	List(ctx context.Context, pathfinderList *v1.PathFinderList, opts PathFinderListOption) error
	Update(ctx context.Context, pathfinder *v1.PathFinder, opts ...client.UpdateOption) error
}

// PathFinderV1Impl Not for ploymorphism, but for parameter overview
type PathFinderV1Impl struct {
	client    client.Client
	namespace string
}

// Create create a new pathfinder
func (pfv1 PathFinderV1Impl) Create(ctx context.Context, pathfinder *v1.PathFinder, opts ...client.CreateOption) error {
	return pfv1.client.Create(ctx, pathfinder, opts...)
}

// Delete a pathfinder from namespace
func (pfv1 PathFinderV1Impl) Delete(ctx context.Context, pathfinder *v1.PathFinder, opts ...client.DeleteOption) error {
	return pfv1.client.Delete(ctx, pathfinder, opts...)

}

// Get pathfinder object using namespaced name
func (pfv1 PathFinderV1Impl) Get(ctx context.Context, name string, pathfinder *v1.PathFinder) error {
	key := client.ObjectKey{
		Namespace: pfv1.namespace,
		Name:      name,
	}
	err := pfv1.client.Get(ctx, key, pathfinder)
	if err != nil {
		return err
	}
	return nil
}

// GetByRegion get pathfinder by namespaced-region
func (pfv1 PathFinderV1Impl) GetByRegion(ctx context.Context, region string, pathfinder *v1.PathFinder) error {
	pfl := v1.PathFinderList{}
	err := pfv1.List(ctx, &pfl, PathFinderListOption{
		Region: region,
	})
	if err != nil {
		return err
	}
	if len(pfl.Items) == 0 {
		return fmt.Errorf("Pathfinder with region %s not found in namespace %s", region, pfv1.namespace)
	}
	pfl.Items[0].DeepCopyInto(pathfinder)
	return nil
}

// List path finder objects
// Besides native list options, you are also able to
// filter using regions
func (pfv1 PathFinderV1Impl) List(ctx context.Context, pathfinderList *v1.PathFinderList, opts PathFinderListOption) error {

	nsOpt := &client.ListOptions{Namespace: pfv1.namespace}
	nsOpt.ApplyToList(&opts.ListOptions)

	err := pfv1.client.List(ctx, pathfinderList, &opts.ListOptions)
	if err != nil {
		return err
	}
	if len(opts.Region) > 0 {
		filted := utils.Filter(
			pathfinderList.Items,
			func(l interface{}) bool { return l.(v1.PathFinder).Spec.Region == opts.Region },
			reflect.TypeOf(v1.PathFinder{}),
		)
		pathfinderList.Items = make([]v1.PathFinder, len(filted))
		for i := 0; i < len(filted); i++ {
			pathfinderList.Items[i] = filted[i].(v1.PathFinder)
		}
	}

	return err
}

// Update a pathfinder
func (pfv1 PathFinderV1Impl) Update(ctx context.Context, pathfinder *v1.PathFinder, opts ...client.UpdateOption) error {
	return pfv1.client.Update(ctx, pathfinder, opts...)
}

// NewPathFinderV1 Create a new pathfinder v1 api
func NewPathFinderV1(client client.Client, namespace string) PathFinderV1 {
	return PathFinderV1Impl{
		client:    client,
		namespace: namespace,
	}
}
