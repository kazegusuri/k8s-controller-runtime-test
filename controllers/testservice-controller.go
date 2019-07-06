package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TestServiceController struct {
	client.Client
	Log logr.Logger
}

func (r *TestServiceController) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	log := r.Log
	log.V(1).Info("reconciling service")
	ctx := context.Background()

	found := true
	var svc corev1.Service
	if err := r.Get(ctx, req.NamespacedName, &svc); err != nil {
		if !apierrors.IsNotFound(err) {
			log.Error(err, "unable to get service")
			return ctrl.Result{}, err
		}

		found = false
	}
	if found {
		if err := r.updateService(ctx, svc); err != nil {
			return ctrl.Result{}, err
		}
	} else {
		log.Info(fmt.Sprintf("service not found %v", svc.Name))
	}

	return ctrl.Result{}, nil
}

func (r *TestServiceController) updateService(ctx context.Context, svc corev1.Service) error {
	objKey := types.NamespacedName{Namespace: svc.Namespace, Name: svc.Name}
	var endpoints corev1.Endpoints
	if err := r.Get(ctx, objKey, &endpoints); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}

		return nil
	}

	for _, subsets := range endpoints.Subsets {
		ips := make([]string, len(subsets.Addresses))
		for i, addr := range subsets.Addresses {
			ips[i] = addr.IP
		}

		r.Log.Info(fmt.Sprintf("[%v] %s \n", objKey, strings.Join(ips, ", ")))
	}

	return nil
}
