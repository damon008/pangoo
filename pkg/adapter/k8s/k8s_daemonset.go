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
	dinstance atomic.Value
	donce     sync.Once
)

func NewDaemonSet() *DaemonSet {
	donce.Do(func() {
		instance.Store(&DaemonSet{ClientSet: pClientSet})
	})
	return dinstance.Load().(*DaemonSet)//&Deployment{ClientSet: pClientSet}//newK8sClientSet()
}

type DaemonSet struct {
	ClientSet *kubernetes.Clientset
}

func (dep *DaemonSet) Create(namespace string, DaemonSet *v1.DaemonSet) (*v1.DaemonSet, error) {
	return dep.ClientSet.AppsV1().DaemonSets(namespace).Create(context.TODO(), DaemonSet, metav1.CreateOptions{})
}

func (dep *DaemonSet) Update(namespace string, DaemonSet *v1.DaemonSet) (*v1.DaemonSet, error) {
	return dep.ClientSet.AppsV1().DaemonSets(namespace).Update(context.TODO(), DaemonSet, metav1.UpdateOptions{})
}

func (dep *DaemonSet) Delete(namespace string, name string) error {
	return dep.ClientSet.AppsV1().DaemonSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (dep *DaemonSet) Get(namespace string, name string) (*v1.DaemonSet, error) {
	return dep.ClientSet.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (dep *DaemonSet) List(namespace string) (*v1.DaemonSetList, error) {
	return dep.ClientSet.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
}