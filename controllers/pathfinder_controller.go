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

	v1 "github.com/6BD-org/pathfinder/api/v1"
	"github.com/6BD-org/pathfinder/consts"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// PathFinderAnnotationKey should present if a service need pathfinder discovery feature
	PathFinderAnnotationKey              = "XM-PathFinder-Service"
	PathFinderRegionKey                  = "XM-PathFinder-Region"
	PathFinderServiceRegistrationNameKey = "XM-PathFinder-ServiceName"

	// PathFinderActivated indicates that this service is ready for discovery
	PathFinderActivated = "Activated"
	// PathFinderDeactiveted indicates that this service is hidden from discovery
	PathFinderDeactiveted   = "Deactivated"
	PathFinderDefaultRegion = "DEFAULT"
)

// PathFinderReconciler reconciles a PathFinder object
type PathFinderReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=xmbsmdsj.com,resources=pathfinders,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch
// +kubebuilder:rbac:groups=xmbsmdsj.com,resources=pathfinders/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete

// Reconcile is the main logic of interpreting PathFinder CRDs
func (r *PathFinderReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("pathfinder", req.NamespacedName)

	svcMap := make(map[string][]corev1.Service)

	serviceList := r.ListServices(req.Namespace)

	for _, svc := range serviceList.Items {
		enabled := verify(&svc)
		if !enabled {
			continue
		}
		region, _ := svcRegion(svc)
		_, ok := svcMap[region]
		if !ok {
			svcMap[region] = make([]corev1.Service, 0)
		}
		svcMap[region] = append(svcMap[region], svc)

	}

	for region := range svcMap {
		svcs, ok := svcMap[region]
		pathFinderRegion, err := r.GetPathFinderRegion(req.Namespace, region)
		oldPathFinderRegion := pathFinderRegion.DeepCopy()

		if err != nil {
			r.Log.Error(err, consts.ERR_GET_PATHFINDER_REGION, "msg", err.Error())
			continue
		}
		if ok {
			r.RebuildPathfinderRegion(pathFinderRegion, svcs)
			if r.shouldUpdate(oldPathFinderRegion, pathFinderRegion) {
				err := r.Update(context.TODO(), pathFinderRegion)
				if err != nil {
					r.Log.Error(
						errors.Errorf(consts.ERR_UPDATE_FAIL),
						consts.ERR_UPDATE_FAIL,
						"msg", err.Error(),
					)
				}
			}
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
