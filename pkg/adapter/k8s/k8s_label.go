package k8s

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

func NewLabel() *Label {
	return &Label{ClientSet: pClientSet}
}

type Label struct {
	ClientSet *kubernetes.Clientset
}

func (lb *Label) Bind(name string, data []byte) error {
	//_, err := lb.ClientSet.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	_, err := lb.ClientSet.CoreV1().Nodes().Patch(context.TODO(), name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
	if err != nil {
		hlog.Error(err)
		return err
	}
	return nil
}

func (lb *Label) UnBind(name string, data []byte) error {
	_, err := lb.ClientSet.CoreV1().Nodes().Patch(context.TODO(), name, types.StrategicMergePatchType, data, metav1.PatchOptions{})
	if err != nil {
		hlog.Error(err)
		return err
	}
	return nil
}
