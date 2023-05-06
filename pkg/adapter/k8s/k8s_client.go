package k8s

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sync"
	vcclient "volcano.sh/apis/pkg/client/clientset/versioned"
)

var pClientSet *kubernetes.Clientset
var pClientSetOnce sync.Once

func NewK8sClient() *kubernetes.Clientset {
	return newK8sClientSet()
}

func newK8sClientSet() *kubernetes.Clientset {
	pClientSetOnce.Do(func() {
		for {
			config, err := rest.InClusterConfig()
			if err != nil {
				hlog.Error("newK8sClientSet InClusterConfig", err)
				continue
			}
			client, err := kubernetes.NewForConfig(config)
			if err != nil {
				hlog.Error("newK8sClientSet NewForConfig", err.Error())
			} else {
				pClientSet = client
				break
			}
		}
	})
	return pClientSet
}

// BuildConfig build kube config ,use rest.InClusterConfig
func BuildConfig(master, kubeconfig string) (*rest.Config, error) {
	if master != "" || kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags(master, kubeconfig)
	}
	return rest.InClusterConfig()
}

// CreateClients create kube client
func CreateClients(kConfig *rest.Config) kubernetes.Interface {
	kClient, err := kubernetes.NewForConfig(kConfig)
	if err != nil {
		panic(fmt.Errorf("Failed to create KubeClient: %v", err))
	}
	return kClient
}

func CreateVolcanoClients(kConfig *rest.Config) vcclient.Interface {
	vcClient, err := vcclient.NewForConfig(kConfig)
	if err != nil {
		panic(fmt.Errorf("Failed to create Volcano Client: %v", err))
	}
	return vcClient
}
