package k8s

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
	"sync/atomic"
)

var (
	//这些操作都是原子性的，因此可以安全地在多个goroutine中使用。
	//使用sync.Once和atomic.Value结合实现单例模式，可以保证高并发情况下的高可用性。
	instance atomic.Value
	once     sync.Once
)

func NewDeployment() *Deployment {
	once.Do(func() {
		instance.Store(&Deployment{ClientSet: pClientSet})
	})
	return instance.Load().(*Deployment)//&Deployment{ClientSet: pClientSet}//newK8sClientSet()
}

type Deployment struct {
	ClientSet *kubernetes.Clientset
}

func (dep *Deployment) Create(namespace string, deployment *v1.Deployment) (*v1.Deployment, error) {
	//pClientSet.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
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