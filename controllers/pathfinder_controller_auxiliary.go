package controllers

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "github.com/wylswz/native-discovery/api/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *PathFinderReconciler) GetPathFinderRegion(namespace string, region string) (*v1.PathFinder, error) {
	pl := v1.PathFinderList{}
	if err := r.List(context.TODO(), &pl, client.InNamespace(namespace)); err != nil {
		return nil, err
	}
	if len(pl.Items) > 1 {
		return nil, errors.Errorf("Dup")
	}
	if len(pl.Items) < 1 {
		return nil, errors.Errorf("Not found")
	}
	for _, pf := range pl.Items {
		if pf.Spec.Region == region {
			return &pf, nil
		}
	}
	return nil, errors.Errorf("Not found")

}

func (r *PathFinderReconciler) GetDefaultPathFinderRegion(namespace string) *v1.PathFinder {
	p, err := r.GetPathFinderRegion(namespace, PathFinderDefaultRegion)
	if err != nil {
		r.Log.Error(err, "Error Getting default region")
	}
	return p
}

func (r *PathFinderReconciler) ListServices(namespace string) *corev1.ServiceList {

	services := corev1.ServiceList{}
	r.Client.List(context.TODO(), &services, &client.ListOptions{})
	return &services
}

func BuildUrlFromService(service *corev1.Service) string {
	return fmt.Sprintf("%s.%s.svc", service.Namespace, service.Name)
}
