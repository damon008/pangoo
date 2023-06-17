package k8s

import (
	"context"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewNamespce() *Namespce {
	return &Namespce{ClientSet: pClientSet}
}

type Namespce struct {
	ClientSet *kubernetes.Clientset
}

func (n *Namespce) Exist(namespace string) bool {
	_, err := n.ClientSet.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	return !apierrors.IsNotFound(err)
}

func (n *Namespce) Create(namespace string, labels map[string]string) error {
	k8sNamespace := &v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}

	if labels != nil {
		k8sNamespace.Labels = labels
	}

	_, err := n.ClientSet.CoreV1().Namespaces().Create(context.TODO(), k8sNamespace, metav1.CreateOptions{})
	return err
}

func (n *Namespce) NotEXistAndCreate(nameSpace string) error {
	_, err := n.ClientSet.CoreV1().Namespaces().Get(context.TODO(), nameSpace, metav1.GetOptions{})
	if err == nil {
		return nil
	}
	k8sNamespace := &v1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: nameSpace,
		},
	}
	_, err = n.ClientSet.CoreV1().Namespaces().Create(context.TODO(), k8sNamespace, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
