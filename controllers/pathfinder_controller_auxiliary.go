package controllers

import (
	"context"
	"fmt"
	"strings"

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

// GetDefaultPathFinderRegion get default region
func (r *PathFinderReconciler) GetDefaultPathFinderRegion(namespace string) (*v1.PathFinder, error) {
	p, err := r.GetPathFinderRegion(namespace, PathFinderDefaultRegion)
	return p, err
}

// ListServices lists all services under a namespace
func (r *PathFinderReconciler) ListServices(namespace string) *corev1.ServiceList {

	services := corev1.ServiceList{}
	r.Client.List(context.TODO(), &services, &client.ListOptions{})
	return &services
}

//BuildURLFromService build a domain name from a service
func BuildURLFromService(service *corev1.Service, port int32) string {
	return fmt.Sprintf("%s.%s.svc:%v", service.Namespace, service.Name, port)
}

// UpdatePathFinderWithService update pathfinder resource given verified service
func (r *PathFinderReconciler) UpdatePathFinderWithService(pf *v1.PathFinder, svc *corev1.Service) {
	serviceName, _ := svc.Annotations[PathFinderServiceRegistrationNameKey]
	if len(svc.Spec.Ports) == 0 {
		return
	}

	if len(svc.Spec.Ports) == 1 {
		r.updatePf(pf, svc, serviceName)
	} else {
		for _, p := range svc.Spec.Ports {
			newServiceName := strings.Join([]string{serviceName, p.Name}, "-")
			fmt.Println(pf, svc, newServiceName)
			r.updatePf(pf, svc, newServiceName)
		}
	}
}

func (r *PathFinderReconciler) updatePf(pf *v1.PathFinder, svc *corev1.Service, serviceName string) {
	existing, ok := pf.Spec.FindServiceEntry(serviceName)
	if !ok {
		r.Log.Info("Add service entry", "service", svc.Name, "namespace", svc.Namespace)
		pf.Spec.ServiceEntries = append(pf.Spec.ServiceEntries, v1.ServiceEntry{
			ServiceName: serviceName,
			ServiceHost: BuildURLFromService(svc, svc.Spec.Ports[0].Port),
		})
	} else {
		if existing.ServiceHost == BuildURLFromService(svc, svc.Spec.Ports[0].Port) {
			r.Log.Info("Unchanged", "service", svc.Name, "namespace", svc.Namespace)
		} else {
			r.Log.Info("Changed", "service", svc.Name, "namespace", svc.Namespace)
			existing.ServiceHost = BuildURLFromService(svc, svc.Spec.Ports[0].Port)
		}
	}
}
