/*
Copyright 2019 GDS.

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

package metadata

import (
	"context"
	"reflect"

	verifyv1beta1 "github.com/alphagov/verify-metadata-controller/pkg/apis/verify/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

// Add creates a new Metadata Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMetadata{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("metadata-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Metadata
	err = c.Watch(&source.Kind{Type: &verifyv1beta1.Metadata{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to created Deployments
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &verifyv1beta1.Metadata{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to created ConfigMaps
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &verifyv1beta1.Metadata{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to created Services
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &verifyv1beta1.Metadata{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMetadata{}

// ReconcileMetadata reconciles a Metadata object
type ReconcileMetadata struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Metadata object and makes changes based on the state read
// and what is in the Metadata.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=,resources=configmaps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=,resources=services/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=verify.gov.uk,resources=metadata,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=verify.gov.uk,resources=metadata/status,verbs=get;update;patch
func (r *ReconcileMetadata) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Metadata instance
	instance := &verifyv1beta1.Metadata{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// TODO: shell out to ruby generate_metadata
	// xmlBlob := generateMetadata(instance.Spec) // os.Exec("the ryby script")

	// TODO: ask HSM to sign the xmlBlob
	// signedXmlBlob := signMetadata(xmlBlob, os.Getenv("PCKS11_ADDRESS"), os.Getenv("PCKS11_PIN"), etc) // os.Exec("pkcs11-tool ??? xmlBlob")

	// TODO: generate ConfigMap containing signedXmlBlob
	metadataConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		BinaryData: map[string][]byte{
			"metadata.xml": []byte("<xml>" + instance.Spec.EntityID + "</xml>"),
		},
	}
	if err := controllerutil.SetControllerReference(instance, metadataConfigMap, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	// Find or craete metadataConfigMap
	foundConfigMap := &corev1.ConfigMap{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: metadataConfigMap.Name, Namespace: metadataConfigMap.Namespace}, foundConfigMap)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating ConfigMap", "namespace", metadataConfigMap.Namespace, "name", metadataConfigMap.Name)
		err = r.Create(context.TODO(), metadataConfigMap)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}
	// Update the found object and write the result back if there are any changes
	if !reflect.DeepEqual(metadataConfigMap.BinaryData, foundConfigMap.BinaryData) {
		foundConfigMap.BinaryData = metadataConfigMap.BinaryData
		log.Info("Updating ConfigMap", "namespace", metadataConfigMap.Namespace, "name", metadataConfigMap.Name)
		err = r.Update(context.TODO(), foundConfigMap)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	metadataLabels := map[string]string{
		"deployment": instance.Name + "-deployment",
	}

	metadataDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: metadataLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: metadataLabels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/usr/share/nginx/html",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{Name: metadataConfigMap.Name},
								},
							},
						},
					},
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(instance, metadataDeployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	// Find or craete metadataDeployment
	foundDeployment := &appsv1.Deployment{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: metadataDeployment.Name, Namespace: metadataDeployment.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Deployment", "namespace", metadataDeployment.Namespace, "name", metadataDeployment.Name)
		err = r.Create(context.TODO(), metadataDeployment)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Update the found object and write the result back if there are any changes
	if !reflect.DeepEqual(metadataDeployment.Spec, foundDeployment.Spec) {
		foundDeployment.Spec = metadataDeployment.Spec
		log.Info("Updating Deployment", "namespace", metadataDeployment.Namespace, "name", metadataDeployment.Name)
		err = r.Update(context.TODO(), foundDeployment)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	metadataService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: metadataLabels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
		},
	}
	if err := controllerutil.SetControllerReference(instance, metadataService, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Find or craete metadataService
	foundService := &corev1.Service{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: metadataService.Name, Namespace: metadataService.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Service", "namespace", metadataService.Namespace, "name", metadataService.Name)
		err = r.Create(context.TODO(), metadataService)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}
	// Update the found object and write the result back if there are any changes
	if !reflect.DeepEqual(metadataService.Spec, foundService.Spec) {
		foundService.Spec = metadataService.Spec
		log.Info("Updating Service", "namespace", metadataService.Namespace, "name", metadataService.Name)
		err = r.Update(context.TODO(), foundService)
		if err != nil {
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}