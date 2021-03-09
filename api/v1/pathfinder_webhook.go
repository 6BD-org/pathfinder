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

package v1

import (
	"context"

	"github.com/6BD-org/pathfinder/consts"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var pathfinderlog = logf.Log.WithName("pathfinder-resource")
var k8sClient client.Client
var hookScheme = runtime.NewScheme()

func (r *PathFinder) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-pathfinder-xmbsmdsj-com-v1-pathfinder,mutating=true,failurePolicy=fail,groups=xmbsmdsj.com,resources=pathfinders,verbs=create;update,versions=v1,name=mpathfinder.kb.io

var _ webhook.Defaulter = &PathFinder{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *PathFinder) Default() {
	pathfinderlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-pathfinder-xmbsmdsj-com-v1-pathfinder,mutating=false,failurePolicy=fail,groups=xmbsmdsj.com,resources=pathfinders,versions=v1,name=vpathfinder.kb.io

var _ webhook.Validator = &PathFinder{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *PathFinder) ValidateCreate() error {
	pathfinderlog.Info("validate create", "name", r.Name)
	// TODO(user): fill in your validation logic upon object creation.
	err := r.CheckDuplication()
	if err != nil {
		return err
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *PathFinder) ValidateUpdate(old runtime.Object) error {
	// TODO Optimize duplication check
	oldPf := PathFinder{}
	err := hookScheme.Convert(&old, &oldPf, context.Background())
	if err != nil {
		return err
	}

	if oldPf.Spec.Region != r.Spec.Region {
		return r.CheckDuplication()
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *PathFinder) ValidateDelete() error {
	pathfinderlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.

	return nil
}

// CheckDuplication Check duplicated pathfinder in same namespace and same region
func (r *PathFinder) CheckDuplication() error {
	pfl := PathFinderList{}
	err := getClient().List(context.TODO(), &pfl, client.InNamespace(r.Namespace))
	if err != nil {
		return err
	}
	for _, pf := range pfl.Items {
		if pf.Spec.Region == r.Spec.Region {
			return errors.Errorf("Duplicated pathfinder in same region")
		}
	}
	return nil
}

func getClient() client.Client {
	var err error
	if k8sClient == nil {
		k8sClient, err = client.New(ctrl.GetConfigOrDie(), client.Options{
			Scheme: hookScheme,
		})
		if err != nil {
			panic(err)
		}
		return k8sClient
	} else {
		return k8sClient
	}

}

func init() {
	var err error

	_ = clientgoscheme.AddToScheme(hookScheme)
	_ = AddToScheme(hookScheme)

	if err != nil {
		pathfinderlog.Error(err, consts.ERR_WEBHOOK_INIT_FAIL)
	}
}
