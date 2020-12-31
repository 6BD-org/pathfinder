package controllers

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	v1 "github.com/6BD-org/pathfinder/api/v1"
	"github.com/6BD-org/pathfinder/messages"
	"github.com/6BD-org/pathfinder/utils"
	"github.com/pkg/errors"
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

func (r *PathFinderReconciler) ListPathFinders(namespace string) (*v1.PathFinderList, error) {
	pl := v1.PathFinderList{}
	if err := r.List(context.TODO(), &pl, client.InNamespace(namespace)); err != nil {
		return nil, err
	}
	return &pl, nil
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

// CleanUpServices Clean up deleted services
// Service entries are cleaned up if:
// 1. The service has been deleted, that is for service entry svc-a/admin-server and svc-a/debug, if
// 		svc-a has been removed from k8s, then both of them are removed
// 2. The port with corresponding name has been removed or renamed.
func (r *PathFinderReconciler) CleanUpServices(pf *v1.PathFinder, svcs []corev1.Service) {
	cleanUpPorts(pf, svcs)
}

// UpdatePathFinderWithService update pathfinder resource given verified service
// Service entries are updated if:
// 1. The address of service entry is modified
// 2. New service entries are found
func (r *PathFinderReconciler) UpdatePathFinderWithService(pf *v1.PathFinder, svc *corev1.Service) {

	if len(svc.Spec.Ports) == 0 {
		return
	}

	for _, p := range svc.Spec.Ports {
		r.updatePf(pf, svc, p)
	}

}

func (r *PathFinderReconciler) updatePf(pf *v1.PathFinder, svc *corev1.Service, port corev1.ServicePort) {
	serviceName := svcRegistractionName(*svc)
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

// cleanUp port entries that are nolonger in service
func cleanUpPorts(pf *v1.PathFinder, svcs []corev1.Service) {

	entries := utils.Filter(
		pf.Status.ServiceEntries,
		func(e interface{}) bool {
			return !toRemove(e.(v1.ServiceEntry), svcs)
		},
		reflect.TypeOf(v1.ServiceEntry{}),
	)

	pf.Status.ServiceEntries = make([]v1.ServiceEntry, len(entries))
	for i, v := range entries {
		pf.Status.ServiceEntries[i] = v.(v1.ServiceEntry)
	}

}

// TODO: Improve lookup efficiency
func toRemove(entry v1.ServiceEntry, svcs []corev1.Service) bool {
	var svc corev1.Service
	svcName, portName := deformatServiceName(entry.ServiceName)

	serviceFound := false
	for _, svc = range svcs {
		if svc.Name == svcName {
			serviceFound = true
			break
		}
	}

	if !serviceFound {
		return true
	}

	res := true

	for _, p := range svc.Spec.Ports {
		if portName == p.Name {
			res = false
		}
	}
	// Remove if from same service and port not found
	return res

}

func deformatServiceName(entryServiceName string) (string, string) {
	sp := strings.Split(entryServiceName, "/")
	if len(sp) == 0 {
		return "", ""
	}
	if len(sp) == 1 {
		return sp[0], ""
	}
	return sp[0], sp[1]

}

func formatServiceName(service string, portName string) string {
	if len(portName) == 0 {
		return service
	}
	return fmt.Sprintf("%s/%s", service, portName)
}

func svcRegion(svc corev1.Service) string {
	return svc.Annotations[PathFinderRegionKey]
}

func svcRegistractionName(svc corev1.Service) string {
	return svc.Annotations[PathFinderServiceRegistrationNameKey]
}
