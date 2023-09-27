package resource

import (
	"errors"
	spacecraftv1 "github.com/xiaofengzs/spacecraft/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type service struct {
}

func NewService() Reconcilable {
	return &service{}
}

func (s service) MakeResource(carrier *spacecraftv1.Spacecraft) client.Object {
	selector := map[string]string{
		"app": carrier.Name,
	}
	svc := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "v1",
			APIVersion: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      carrier.Name,
			Namespace: carrier.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(carrier, spacecraftv1.GroupVersion.WithKind("Spacecraft")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: selector,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       carrier.Spec.ServiceConfig.Port,
					TargetPort: intstr.IntOrString{IntVal: carrier.Spec.ServiceConfig.TargetPort},
				},
			},
		},
	}

	if carrier.Spec.ServiceConfig.Type != "" && carrier.Spec.ServiceConfig.NodePort != 0 {
		svc.Spec.Type = corev1.ServiceTypeNodePort
		svc.Spec.Ports = []corev1.ServicePort{
			{
				Port:       carrier.Spec.ServiceConfig.Port,
				TargetPort: intstr.IntOrString{IntVal: carrier.Spec.ServiceConfig.TargetPort},
				NodePort:   carrier.Spec.ServiceConfig.NodePort,
			},
		}
	}

	return svc
}

func (s service) MakeDefaultResource() client.Object {
	return &corev1.Service{}
}

func (s service) Validate(carrier *spacecraftv1.Spacecraft) error {
	if (carrier.Spec.ServiceConfig.Type != "" && carrier.Spec.ServiceConfig.NodePort == 0) ||
		carrier.Spec.ServiceConfig.Type == "" && carrier.Spec.ServiceConfig.NodePort != 0 {
		return errors.New("when you enabler service, should set service config type and nodePort")
	}
	return nil
}
