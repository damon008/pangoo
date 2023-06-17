package k8s

import (
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

//默认是对当前服务所在的集群
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

// ClientOptions used to build kube rest config.
type ClientOptions struct {
	Master     string
	KubeConfig string
	QPS        float32
	Burst      int
}

// BuildConfig build kube config ,use rest.InClusterConfig
//需要对某个集群进行处理，就需要知道是哪个集群
// BuildConfig builds kube rest config with the given options.
func BuildConfig(opt ClientOptions) (*rest.Config, error) {
	var cfg *rest.Config
	var err error

	master := opt.Master
	kubeconfig := opt.KubeConfig
	cfg, err = clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		return nil, err
	}
	cfg.QPS = opt.QPS
	cfg.Burst = opt.Burst

	return cfg, nil
}


// 根据上面的config，生成对于自定义集群的client CreateClients create kube client
func CreateClients(kConfig *rest.Config) kubernetes.Interface {
	kClient, err := kubernetes.NewForConfig(kConfig)
	if err != nil {
		hlog.Errorf("Failed to create KubeClient: %v", err)
	}
	return kClient
}

func CreateClientsByCluster(master, kubeconfig string) kubernetes.Interface {
	var config *rest.Config
	var err error
	if master != "" || kubeconfig != "" {
		config,_ = clientcmd.BuildConfigFromFlags(master, kubeconfig)
	} else {
		config,_ = rest.InClusterConfig()
	}
	if (config == nil) {
		return nil
	}
	//config.TLSClientConfig.Insecure = true
	pClientSet,err = kubernetes.NewForConfig(config)
	if err != nil {
		hlog.Errorf("Failed to create KubeClient: %v", err)
	}
	return pClientSet
}

// 根据上面的config，生成对于自定义集群的volcano client
func CreateVolcanoClients(kConfig *rest.Config) vcclient.Interface {
	vcClient, err := vcclient.NewForConfig(kConfig)
	if err != nil {
		hlog.Errorf("Failed to create Volcano Client: %v", err)
	}
	return vcClient
}
