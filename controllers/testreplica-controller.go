package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TestReplicaController struct {
	client.Client
	Log logr.Logger
}

func (r *TestReplicaController) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	log := r.Log
	log.V(1).Info("reconciling replica set")
	ctx := context.Background()

	found := true
	var rs appsv1.ReplicaSet
	if err := r.Get(ctx, req.NamespacedName, &rs); err != nil {
		if !apierrors.IsNotFound(err) {
			log.Error(err, "unable to get replica")
			return ctrl.Result{}, err
		}

		found = false
	}

	if found {
		log.Info(fmt.Sprintf("replica set %v", rs.Name))
	} else {
		log.Info(fmt.Sprintf("replica set not found %v", rs.Name))
	}

	return ctrl.Result{}, nil
}
