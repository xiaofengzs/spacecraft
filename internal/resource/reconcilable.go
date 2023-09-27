package resource

import (
	spacecraftv1 "github.com/xiaofengzs/spacecraft/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconcilable interface {
	MakeResource(carrier *spacecraftv1.Spacecraft) client.Object
	MakeDefaultResource() client.Object
	Validate(carrier *spacecraftv1.Spacecraft) error
}
