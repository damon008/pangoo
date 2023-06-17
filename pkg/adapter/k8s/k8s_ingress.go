package k8s

import (
	"context"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewIngress() *Ingress {
	return &Ingress{ClientSet: pClientSet}
}

type Ingress struct {
	ClientSet *kubernetes.Clientset
}

func (ing *Ingress) Create(namespace string, ingress *v1beta1.Ingress) error {
	_, err := ing.ClientSet.ExtensionsV1beta1().Ingresses(namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	return err
}

func (ing *Ingress) Delete(namespace, name string) error {
	return ing.ClientSet.ExtensionsV1beta1().Ingresses(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (ing *Ingress) Get(namespace, name string) (*v1beta1.Ingress, error) {
	return ing.ClientSet.ExtensionsV1beta1().Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
