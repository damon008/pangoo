package k8s

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type K8sServicePara = v1.Service

func NewService() *Service {
	return &Service{ClientSet: newK8sClientSet()}
}

type Service struct {
	ClientSet *kubernetes.Clientset
}

func (svc *Service) Create(namespce string, service *v1.Service) error {
	_, err := svc.ClientSet.CoreV1().Services(namespce).Create(context.TODO(), service, metav1.CreateOptions{})
	return err
}

func (svc *Service) Delete(namespace, name string) error {
	return svc.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (svc *Service) Get(namespace, name string) (*v1.Service, error) {
	return svc.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
