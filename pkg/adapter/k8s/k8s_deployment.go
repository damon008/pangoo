package k8s

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewDeployment() *Deployment {
	return &Deployment{ClientSet: newK8sClientSet()}
}

type Deployment struct {
	ClientSet *kubernetes.Clientset
}

func (dep *Deployment) Create(namespace string, deployment *v1.Deployment) (*v1.Deployment, error) {
	return dep.ClientSet.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
}

func (dep *Deployment) Update(namespace string, deployment *v1.Deployment) (*v1.Deployment, error) {
	return dep.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
}

func (dep *Deployment) Delete(namespace string, name string) error {
	return dep.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (dep *Deployment) Get(namespace string, name string) (*v1.Deployment, error) {
	return dep.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (dep *Deployment) List(namespace string) (*v1.DeploymentList, error) {
	return dep.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
}