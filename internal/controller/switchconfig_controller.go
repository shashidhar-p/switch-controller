/*
Copyright 2025.

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

package controller

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"golang.org/x/crypto/ssh"
	vulcanv1 "vulcan/switch-controller/api/v1"
)

// SwitchConfigReconciler reconciles a SwitchConfig object
type SwitchConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=vulcan.vulcan,resources=switchconfigs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=vulcan.vulcan,resources=switchconfigs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=vulcan.vulcan,resources=switchconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SwitchConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *SwitchConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var sc vulcanv1.SwitchConfig
	if err := r.Get(ctx, req.NamespacedName, &sc); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	spec := sc.Spec
	addr := fmt.Sprintf("%s:22", spec.SwitchIP)

	sshConfig := &ssh.ClientConfig{
		User: spec.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(spec.SSHPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		log.Error(err, "SSH dial failed", "addr", addr)
		sc.Status.Phase = "Failed"
		sc.Status.Message = fmt.Sprintf("SSH dial error: %v", err)
		_ = r.Status().Update(ctx, &sc)
		return ctrl.Result{}, err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Error(err, "SSH session failed")
		sc.Status.Phase = "Failed"
		sc.Status.Message = fmt.Sprintf("SSH session error: %v", err)
		_ = r.Status().Update(ctx, &sc)
		return ctrl.Result{}, err
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	cmd := spec.Config
	if err := session.Run(cmd); err != nil {
		log.Error(err, "Failed to run config", "stderr", stderr.String())
		sc.Status.Phase = "Failed"
		sc.Status.Message = fmt.Sprintf("Config apply error: %v, %s", err, stderr.String())
		_ = r.Status().Update(ctx, &sc)
		return ctrl.Result{}, err
	}

	log.Info("Config applied successfully", "stdout", stdout.String())
	sc.Status.Phase = "Successful"
	sc.Status.Message = "Configuration applied successfully"
	if err := r.Status().Update(ctx, &sc); err != nil {
		log.Error(err, "Failed to update Status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwitchConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vulcanv1.SwitchConfig{}).
		Named("switchconfig").
		Complete(r)
}
