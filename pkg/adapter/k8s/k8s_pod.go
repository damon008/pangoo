package k8s

import (
	"context"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func NewPodService() *PodService {
	return &PodService{ClientSet: pClientSet}
}

type Container = v1.Container
type Pod = v1.Pod
type ContainerPort = v1.ContainerPort
type Protocol = v1.Protocol
type ResourceRequirements = v1.ResourceRequirements
type ResourceList = v1.ResourceList
type VolumeMount = v1.VolumeMount

const (
	// PodRunning Pod Phase Running
	PodRunning v1.PodPhase = "Running"
	// PodFailed Pod Phase Failed
	PodFailed v1.PodPhase = "Failed"
)

const (
	// ProtocolTCP is the TCP protocol.
	ProtocolTCP Protocol = "TCP"
	// ProtocolUDP is the UDP protocol.
	ProtocolUDP Protocol = "UDP"
	// ProtocolSCTP is the SCTP protocol.
	ProtocolSCTP Protocol = "SCTP"
)

type PodService struct {
	ClientSet *kubernetes.Clientset
}

func (ps *PodService) Create(pod *Pod) (*Pod, error) {
	return ps.ClientSet.CoreV1().Pods(pod.ObjectMeta.Namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
}

func (ps *PodService) Delete(namespace, name string) error {
	return ps.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (ps *PodService) Get(namespace, name string) (*Pod, error) {
	return ps.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
