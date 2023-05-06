package k8s

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

type ConfigMapPara struct {
	Name      string
	JobName   string
	JobUID    interface{}
	Namespace string
	Data      map[string]string
}

func NewConfigMap() *ConfigMap {
	return &ConfigMap{ClientSet: newK8sClientSet()}
}

type ConfigMap struct {
	ClientSet *kubernetes.Clientset
}

func (cm *ConfigMap) Create(configMapPara *ConfigMapPara) (*corev1.ConfigMap, error) {
	trueVar := true
	refs := []metav1.OwnerReference{
		metav1.OwnerReference{
			APIVersion:         "batch.volcano.sh/v1alpha1",
			Kind:               "Job",
			Name:               configMapPara.JobName,
			UID:                configMapPara.JobUID.(types.UID),
			BlockOwnerDeletion: &trueVar,
			Controller:         &trueVar,
		},
	}

	configMap, err := cm.ClientSet.CoreV1().ConfigMaps(configMapPara.Namespace).Create(context.TODO(), &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:            configMapPara.Name,
			Namespace:       configMapPara.Namespace,
			OwnerReferences: refs,
		},
		Data: configMapPara.Data,
	}, metav1.CreateOptions{})
	if err != nil {
		hlog.Error("Create ConfigMap Failed %v", err)
		return nil, err
	}
	return configMap, nil
}

func (cm *ConfigMap) CreateHorovodConfigMapIfEnableHpa(namespace string) error{
	configMapPara := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "horovod-host-discover",
			Namespace: namespace,
		},
		Data: map[string]string{
			"discover_hosts.sh": "for hostFile in $(ls /etc/volcano/ | grep .host);do\n      for line in $(cat /etc/volcano/\"$hostFile\");do\n          gpuNum=$(ssh -q -o StrictHostKeyChecking=no $line \"ls /dev/nvidia* | grep \"[0-9]$\" | wc -l\")\n          if [ -n \"$gpuNum\" ];then\n              echo \"$line:$gpuNum\"\n          fi\n      done\n    done\n",
		},
	}
	_, err := cm.Get(namespace, configMapPara.Name)
	if err == nil {
		hlog.Info("horovod configMap exsit")
		return nil
	}
	if !errors.IsNotFound(err) {
		return err
	}

	if _, err := cm.ClientSet.CoreV1().ConfigMaps(namespace).Create(context.TODO(), configMapPara, metav1.CreateOptions{}); err != nil {
		return err
	}
	return nil
}


func (cm *ConfigMap) Delete(configMapPara *ConfigMapPara) error {
	err := cm.ClientSet.CoreV1().ConfigMaps(configMapPara.Namespace).Delete(context.TODO(), configMapPara.Name, metav1.DeleteOptions{})
	if err != nil {
		hlog.Error("Delete Configmap On K8s Failed %v", err)
		return err
	}
	return nil
}

func (cm *ConfigMap) Get(namespace, name string) (*corev1.ConfigMap, error) {
	configMap, err := cm.ClientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		hlog.Error("get Configmap On K8s Failed %v", err)
		return nil, err
	}
	return configMap, nil
}
func (cm *ConfigMap) Update(nameapace string, configMap *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	configMap, err := cm.ClientSet.CoreV1().ConfigMaps(nameapace).Update(context.TODO(), configMap, metav1.UpdateOptions{})
	if err != nil {
		hlog.Error("Update ConfigMap Failed %v", err)
		return nil, err
	}
	return configMap, nil
}
