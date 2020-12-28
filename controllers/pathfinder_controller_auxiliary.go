package controllers

import (
	"context"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	v1 "github.com/wylswz/native-discovery/api/v1"
	"github.com/wylswz/native-discovery/messages"
	"github.com/wylswz/native-discovery/utils"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *PathFinderReconciler) GetPathFinderRegion(namespace string, region string) (*v1.PathFinder, error) {
	pl := v1.PathFinderList{}
	if err := r.List(context.TODO(), &pl, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	filted := utils.Filter(
		pl.Items,
		func(pf interface{}) bool { return pf.(v1.PathFinder).Spec.Region == region },
		reflect.TypeOf(v1.PathFinder{}),
	)

	if len(filted) > 1 {
		return nil, errors.Errorf(messages.DuplicatedRegion, namespace, region)
	}
	if len(filted) < 1 {
		return nil, errors.Errorf(messages.RegionNotFound, namespace, region)
	}
	for _, pf := range pl.Items {
		if pf.Spec.Region == region {
			return &pf, nil
		}
	}
	return nil, errors.Errorf(messages.RegionNotFound, namespace, region)

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

	if len(svc.Spec.Ports) == 0 {
		return
	}

	for _, p := range svc.Spec.Ports {
		r.updatePf(pf, svc, p)
	}

}

func (r *PathFinderReconciler) updatePf(pf *v1.PathFinder, svc *corev1.Service, port corev1.ServicePort) {
	serviceName, _ := svc.Annotations[PathFinderServiceRegistrationNameKey]
	serviceName = formatServiceName(serviceName, port.Name)
	existing, ok := pf.Status.FindServiceEntry(serviceName)
	if !ok {
		r.Log.Info("Add service entry", "service", svc.Name, "namespace", svc.Namespace)
		pf.Status.ServiceEntries = append(pf.Status.ServiceEntries, v1.ServiceEntry{
			ServiceName: serviceName,
			ServiceHost: BuildURLFromService(svc, port.Port),
			Payload: v1.Payload{
				KeyValPairs: []v1.PayloadKeyValPair{},
			},
		})
	} else {
		// TODO: update and generate payload
		if existing.ServiceHost == BuildURLFromService(svc, port.Port) {
			r.Log.Info("Unchanged", "service", svc.Name, "namespace", svc.Namespace)
		} else {
			r.Log.Info("Changed", "service", svc.Name, "namespace", svc.Namespace)
			existing.ServiceHost = BuildURLFromService(svc, port.Port)
			existing.Payload.KeyValPairs = []v1.PayloadKeyValPair{}
		}
	}
}

func formatServiceName(service string, portName string) string {
	return fmt.Sprintf("%s/%s", service, portName)
}
