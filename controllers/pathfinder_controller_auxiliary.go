package controllers

import (
	"context"
	"fmt"
	"reflect"

	v1 "github.com/6BD-org/pathfinder/api/v1"
	"github.com/6BD-org/pathfinder/common"
	"github.com/6BD-org/pathfinder/consts"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GetPathFinderRegion Find pathfinder in a specified region
func (r *PathFinderReconciler) GetPathFinderRegion(namespace string, region string) (*v1.PathFinder, error) {
	return common.GetPathFinderRegion(r.Client, namespace, region)
}

// ListPathFinders List all pathfinders in a namespace
func (r *PathFinderReconciler) ListPathFinders(namespace string) (*v1.PathFinderList, error) {
	return common.ListPathFinders(r.Client, namespace)
}

// GetDefaultPathFinderRegion get default region
func (r *PathFinderReconciler) GetDefaultPathFinderRegion(namespace string) (*v1.PathFinder, error) {
	p, err := r.GetPathFinderRegion(namespace, PathFinderDefaultRegion)
	return p, err
}

// ListServices lists all services under a namespace
func (r *PathFinderReconciler) ListServices(namespace string) *corev1.ServiceList {

	services := corev1.ServiceList{}
	r.Client.List(context.TODO(), &services, client.InNamespace(namespace))
	return &services
}

func (r *PathFinderReconciler) shouldUpdate(oldPf *v1.PathFinder, pf *v1.PathFinder) bool {
	return (!reflect.DeepEqual(pf.Spec, oldPf.Spec)) ||
		(!reflect.DeepEqual(pf.Status, oldPf.Status))
}

// RebuildPathfinderRegion Rebuild pathfinder from services from that region
func (r *PathFinderReconciler) RebuildPathfinderRegion(pf *v1.PathFinder, svcs []corev1.Service) error {
	svcEntries := make([]v1.ServiceEntry, 0)
	for _, svc := range svcs {
		region, ok := svcRegion(svc)
		if !ok {
			r.Log.Info(consts.WARN_REGION_UNSPECIFIED)
			region = "DEFAULT"
		}
		if region != pf.Spec.Region {
			r.Log.Info(consts.WARN_REGION_INCONSISTENT, "namespace", svc.Namespace, "svc", svc.Name)
		} else {
			for _, p := range svc.Spec.Ports {
				name, _ := svcRegistractionName(svc)
				entry := v1.ServiceEntry{
					ServiceName: formatServiceName(name, p.Name),
					ServiceHost: buildURLFromService(svc, p.Port),
					Payload: v1.Payload{
						KeyValPairs: make([]v1.PayloadKeyValPair, 0),
					},
				}
				svcEntries = append(svcEntries, entry)

			}
		}

	}
	pf.Status.ServiceEntries = svcEntries
	return nil
}

//BuildURLFromService build a domain name from a service
func buildURLFromService(service corev1.Service, port int32) string {
	return fmt.Sprintf("%s.%s.svc:%v", service.Name, service.Namespace, port)
}

func formatServiceName(service string, portName string) string {
	if len(portName) == 0 {
		return service
	}
	return fmt.Sprintf("%s/%s", service, portName)
}

func svcRegion(svc corev1.Service) (string, bool) {
	k, ok := svc.Annotations[PathFinderRegionKey]
	return k, ok
}

func svcRegistractionName(svc corev1.Service) (string, bool) {
	n, ok := svc.Annotations[PathFinderServiceRegistrationNameKey]
	return n, ok
}

func svcPathFinderEnabled(svc corev1.Service) bool {
	p, ok := svc.Annotations[PathFinderAnnotationKey]
	if !ok {
		return false
	} else {
		return p == PathFinderActivated
	}
}

func verify(svc *corev1.Service) bool {
	// Verify enabled
	enabled := svcPathFinderEnabled(*svc)
	if !enabled {
		return false
	}

	_, ok := svcRegistractionName(*svc)
	if !ok {
		return false
	}

	_, ok = svcRegion(*svc)
	if !ok {
		svc.Annotations[PathFinderRegionKey] = "DEFAULT"
	}

	return true
}
