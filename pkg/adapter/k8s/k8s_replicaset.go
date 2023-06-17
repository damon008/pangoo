package k8s

import (
	"context"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewReplicaSet() *ReplicaSet {
	return &ReplicaSet{ClientSet: pClientSet}
}

type ReplicaSet struct {
	ClientSet *kubernetes.Clientset
}

func (dep *ReplicaSet) Create(namespace string, ReplicaSet *v1.ReplicaSet) (*v1.ReplicaSet, error) {
	return dep.ClientSet.AppsV1().ReplicaSets(namespace).Create(context.TODO(), ReplicaSet, metav1.CreateOptions{})
}

func (dep *ReplicaSet) Update(namespace string, ReplicaSet *v1.ReplicaSet) (*v1.ReplicaSet, error) {
	return dep.ClientSet.AppsV1().ReplicaSets(namespace).Update(context.TODO(), ReplicaSet, metav1.UpdateOptions{})
}

func (dep *ReplicaSet) Delete(namespace string, name string) error {
	return dep.ClientSet.AppsV1().ReplicaSets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (dep *ReplicaSet) Get(namespace string, name string) (*v1.ReplicaSet, error) {
	return dep.ClientSet.AppsV1().ReplicaSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (dep *ReplicaSet) List(namespace string) (*v1.ReplicaSetList, error) {
	return dep.ClientSet.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
}