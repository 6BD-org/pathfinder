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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pathfinderv1 "github.com/wylswz/native-discovery/api/v1"
	"github.com/wylswz/native-discovery/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	clientSet := k8s.ClientSet(config)
	serviceList, err := clientSet.CoreV1().Services(req.Namespace).List(metav1.ListOptions{})
	for _, s := range serviceList.Items {
		annotations := s.Annotations
		pathFinderState, ok := annotations[PathFinderAnnotationKey]
		if ok && pathFinderState == PathFinderActivated {

		}
	}
	return ctrl.Result{}, nil
}

func (r *PathFinderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pathfinderv1.PathFinder{}).
		Complete(r)
}
