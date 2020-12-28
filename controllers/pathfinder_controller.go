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

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1 "github.com/wylswz/native-discovery/api/v1"
	"github.com/wylswz/native-discovery/messages"
	"github.com/wylswz/native-discovery/utils"
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
	return ctrl.Result{}, nil
}

func (r *PathFinderReconciler) CheckDefaultRegion(req ctrl.Request) {
	ns := req.Namespace
	_, err := r.GetDefaultPathFinderRegion(ns)
	if err != nil {
		// Default region does not exists, create one
		pf := v1.PathFinder{}
		pf.Spec.Region = PathFinderDefaultRegion
		pf.ObjectMeta.Name = "pathfinder-default"
		pf.ObjectMeta.Namespace = req.Namespace
		r.Client.Create(context.TODO(), &pf, &client.CreateOptions{})
	}
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
