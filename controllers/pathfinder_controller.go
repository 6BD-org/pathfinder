/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"log"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/6BD-org/pathfinder/api/v1"
	"github.com/6BD-org/pathfinder/messages"
	"github.com/6BD-org/pathfinder/utils"
)

const (
	// PathFinderAnnotationKey should present if a service need pathfinder discovery feature
	PathFinderAnnotationKey = "XM-PathFinder-Service"

	// PathFinderActivated indicates that this service is ready for discovery
	PathFinderActivated = "Activated"

	// PathFinderDeactiveted indicates that this service is hidden from discovery
	PathFinderDeactiveted = "Deactivated"

	PathFinderRegionKey = "XM-PathFinder-Region"

	PathFinderDefaultRegion = "DEFAULT"

	PathFinderServiceRegistrationNameKey = "XM-PathFinder-ServiceName"
)

// PathFinderReconciler reconciles a PathFinder object
type PathFinderReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=pathfinder.xmbsmdsj.com,resources=pathfinders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=pathfinder.xmbsmdsj.com,resources=pathfinders/status,verbs=get;update;patch

// Reconcile is the main logic of interpreting PathFinder CRDs
func (r *PathFinderReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("pathfinder", req.NamespacedName)

	serviceList := r.ListServices(req.Namespace)
	for _, s := range serviceList.Items {
		annotations := s.Annotations
		pathFinderState, ok := annotations[PathFinderAnnotationKey]
		if ok && pathFinderState == PathFinderActivated {
			region, ok := annotations[PathFinderRegionKey]
			if !ok {
				r.Log.Info("Region unspecified, will use default", "service", s.Name, "namespace", s.Namespace)
				region = PathFinderDefaultRegion
			}
			_, ok = annotations[PathFinderServiceRegistrationNameKey]
			if !ok {
				r.Log.Error(nil, messages.ServiceNameUnspecified)
			}
			pathFinderRegion, err := r.GetPathFinderRegion(req.Namespace, region)
			if err != nil {
				r.Log.Error(err, "Unable to find region")
			} else {

				// Do registrations
				r.UpdatePathFinderWithService(pathFinderRegion, &s)
				r.Log.Info("Updating pathfinder")
				err = r.Update(context.TODO(), pathFinderRegion)

				utils.CheckErr(err, "Error updating", r.Log)
			}

		}

	}

	log.Println("Start cleaning up")

	svcMap := make(map[string][]corev1.Service)
	pfMap := make(map[string]v1.PathFinder)

	pfs, err := r.ListPathFinders(req.Namespace)
	if err != nil {
		log.Println("Error: ", err)
	}
	for _, pf := range pfs.Items {
		pfMap[pf.Spec.Region] = pf
	}

	for _, svc := range serviceList.Items {
		region := svcRegion(svc)
		_, ok := svcMap[region]
		if !ok {
			svcMap[region] = make([]corev1.Service, 0)
		}
		svcMap[svcRegion(svc)] = append(svcMap[svcRegion(svc)], svc)
	}

	for region := range pfMap {
		pf, _ := pfMap[region]
		svcs, ok := svcMap[region]
		if ok {
			r.CleanUpServices(&pf, svcs)
		} else {
			log.Println("No service found for region: ", region)
		}

	}
	return ctrl.Result{}, nil
}

func (r *PathFinderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	err := ctrl.NewControllerManagedBy(mgr).For(&v1.PathFinder{}).Complete(r)
	if err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Service{}).
		Complete(r)
}
