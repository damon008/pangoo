package k8s

import (
	"context"
	"k8s.io/api/autoscaling/v2beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewHPA() *HPA {
	return &HPA{ClientSet: newK8sClientSet()}
}

type HPA struct {
	ClientSet *kubernetes.Clientset
}

func (hpa *HPA) Create(namespace string, horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler) (*v2beta2.HorizontalPodAutoscaler, error) {
	return hpa.ClientSet.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Create(context.TODO(), horizontalPodAutoscaler, metav1.CreateOptions{})
}

func (hpa *HPA) Update(namespace string, horizontalPodAutoscaler *v2beta2.HorizontalPodAutoscaler) (*v2beta2.HorizontalPodAutoscaler, error) {
	return hpa.ClientSet.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Update(context.TODO(), horizontalPodAutoscaler, metav1.UpdateOptions{})
}

func (hpa *HPA) Delete(namespace string, name string) error {
	return hpa.ClientSet.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (hpa *HPA) Get(namespace string, name string) (*v2beta2.HorizontalPodAutoscaler, error) {
	return hpa.ClientSet.AutoscalingV2beta2().HorizontalPodAutoscalers(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
