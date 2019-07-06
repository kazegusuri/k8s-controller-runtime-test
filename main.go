package main

import (
	"fmt"
	"os"

	"github.com/kazegusuri/k8s-controller-runtime-test/controllers"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

func main() {
	ctrl.SetLogger(zap.Logger(true))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	err = ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.ReplicaSet{}).
		Watches(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(obj handler.MapObject) []ctrl.Request {
				owners := obj.Meta.GetOwnerReferences()
				for _, owner := range owners {
					if owner.Kind != "ReplicaSet" {
						continue
					}
					ctrl.Log.Info(fmt.Sprintf("pod %v/%v", obj.Meta.GetNamespace(), obj.Meta.GetName()))
					return []ctrl.Request{
						{
							NamespacedName: types.NamespacedName{
								Namespace: obj.Meta.GetNamespace(),
								Name:      owner.Name,
							},
						},
					}
				}
				return nil
			}),
		}).
		Complete((&controllers.TestReplicaController{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("testreplica"),
		}))
	if err != nil {
		setupLog.Error(err, "unable to create controller")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
