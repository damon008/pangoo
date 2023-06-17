package k8s

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type SecretPara struct {
	Name      string
	JobName   string
	JobUID    interface{}
	Namespace string
	Data      map[string][]byte
}

func NewSecret() *Secret {
	return &Secret{ClientSet: pClientSet}
}

type Secret struct {
	ClientSet *kubernetes.Clientset
}

func (sp *Secret) Create(secretPara *SecretPara) (*corev1.Secret, error) {
	trueVar := true
	refs := []metav1.OwnerReference{
		metav1.OwnerReference{
			APIVersion:         "batch.volcano.sh/v1alpha1",
			Kind:               "Job",
			Name:               secretPara.JobName,
			UID:                secretPara.JobUID.(types.UID),
			BlockOwnerDeletion: &trueVar,
			Controller:         &trueVar,
		},
	}

	secret, err := sp.ClientSet.CoreV1().Secrets(secretPara.Namespace).Create(context.TODO(), &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:            secretPara.Name,
			Namespace:       secretPara.Namespace,
			OwnerReferences: refs,
		},
		Data: secretPara.Data,
	}, metav1.CreateOptions{})
	if err != nil {
		hlog.Error("Create Secret Failed %v", err)
		return nil, err
	}
	return secret, nil
}

func (sp *Secret) Get(namespace, name string) (*corev1.Secret, error) {
	return sp.ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (sp *Secret) Delete(namespace, name string) error {
	return sp.ClientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (sp *Secret) CreateSecret(namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	return sp.ClientSet.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
}

func (sp *Secret) List(namespace string) (*corev1.SecretList, error) {
	secretList, err := sp.ClientSet.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	hlog.Info("secretList:%s", secretList)
	return secretList, nil
}

func (sp *Secret) Update(namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	return sp.ClientSet.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
}
