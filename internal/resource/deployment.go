package resource

import (
	spacecraftv1 "github.com/xiaofengzs/spacecraft/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type deployment struct {
}

func NewDeployment() Reconcilable {
	return &deployment{}
}

func (d deployment) MakeResource(carrier *spacecraftv1.Spacecraft) client.Object {
	labels := map[string]string{
		"app": carrier.Name,
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "apps/corev1",
			APIVersion: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      carrier.Name,
			Namespace: carrier.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(carrier, spacecraftv1.GroupVersion.WithKind("Spacecraft")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &carrier.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:  carrier.Name,
							Image: carrier.Spec.Image,
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									ContainerPort: carrier.Spec.Port,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d deployment) MakeDefaultResource() client.Object {
	return &appsv1.Deployment{}
}

func (d deployment) Validate(carrier *spacecraftv1.Spacecraft) error {
	return nil
}
