package conf

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"pangoo/pkg/filewatcher"
	"path/filepath"
	"sync"
)

var EnvConfig AdmissionConfiguration

type DockerConfig struct {
	Path string `yaml:"path"`
	Repo string `yaml:"repo"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DockerfilePath string `yaml:"dockerfilePath"`
}

type JenkinsConfig struct {
	BaseUrl string `yaml:"baseUrl"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type AdmissionConfiguration struct {
	sync.Mutex
	Config Config
}

type Config struct {
	SwagConfig SwagConfig `yaml:"swag"`
	GitConfig GitConfig `yaml:"git"`
	K8sConfig K8sConfig `yaml:"k8s"`
	DockerConfig DockerConfig `yaml:"docker"`
	JenkinsConfig JenkinsConfig `yaml:"jenkins"`
}

type SwagConfig struct {
	Host string `yaml:"host"`
}

type GitConfig struct {
	ProjectUrl string `yaml:"project_url"`
	Token string `yaml:"token"`
	BranchUrl string `yaml:"branch_url"`
}

type K8sConfig struct {
	KubeConfig string `yaml:"kubeConfig"`
}

type CacheCleanConfig struct {
	CacheCleanFrequencyMinutes int    `yaml:"cleanFrequencyMinutes"`
	JobRetentionTime           int    `yaml:"jobRetention"`
	PushMQCronExpression       string `yaml:"pushCron"`
}

// DBConfig Mysql Configuration
type DBConfig struct {
	Host     string `yaml:"host"`     //地址
	Port     int    `yaml:"port"`     //端口
	Name     string `yaml:"username"` //用户
	Pass     string `yaml:"password"` //密码
	DBName   string `yaml:"database"` //库名
	Charset  string `yaml:"charset"`  //编码
	Timezone string `yaml:"timezone"` //时区
	MaxIdle  int    `yaml:"maxidle"`  //最大空间连接
	MaxOpen  int    `yaml:"maxopen"`  //最大连接数
}

// MinioConfig minio configuration
type MinioConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey"`
	UseSSL          bool   `yaml:"ssl"`
}

type RocketmqConfig struct {
	Endpoint       []string `yaml:"endpoints"`
	JobStatusGroup string   `yaml:"jobGroup"`
}

type KubeConfig struct {
	APIServerAddress string `yaml:"apiserver"`
	KubeConfig       string `yaml:"kubeconfig"`
}

type ElasticConfig struct {
	ElasticProportion float64 `yaml:"elasticProportion"` // 资源预留比例
	ExpandThreshold   float64 `yaml:"expandThreshold"`
	MedinaThreshold   float64 `yaml:"medinaThreshold"`  // TODO 中间阈值: (收缩阈值+扩张阈值)/2
	ScalingThreshold  float64 `yaml:"scalingThreshold"` // 收缩阈值
	ScalingPolicy     string  `yaml:"scalingPolicy"`    // fair/priority
	ExpandPolicy      string  `yaml:"expandPolicy"`     // fair/priority
	HpaCheckMinutes   int     `yaml:"hpaCheckMinutes"`
}

type ClusterResource struct {
	ResourceManagerUrl string `yaml:"resourceManagerUrl"`
}

func (c *Config) ToYaml(file string) {
	marshal, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
		hlog.Errorf("Serialize Config Failed %v", err)
	}
	err = ioutil.WriteFile(file, marshal, 0644)
	if err != nil {
		panic(err)
		hlog.Errorf("Save To File Failed %v", err)
	}
}
func (c *Config) PrintYaml() string {
	marshal, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
		hlog.Errorf("Serialize Config Failed %v", err)
	}
	return string(marshal)
}

func ParseFromYAML(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		hlog.Errorf("File Read Failed %v", err)
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		hlog.Errorf("Parse Config Yaml Failed %v", err)
		return nil, err
	}
	return c, nil
}

/*func (c *Config) GenDataBaseSource() string {
	name, err := base64.StdEncoding.DecodeString(c.DBConfig.Name)
	if err != nil {
		logs.Error("DBConfig.Name decoding failed!")
	}
	pass, err := base64.StdEncoding.DecodeString(c.DBConfig.Pass)
	if err != nil {
		logs.Error("DBConfig.Pass decoding failed!")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local&time_zone=%s",
		name,
		pass,
		c.DBConfig.Host,
		c.DBConfig.Port,
		c.DBConfig.DBName,
		c.DBConfig.Charset,
		url.QueryEscape(c.DBConfig.Timezone),
	)
	return dsn
}*/






func WatchAdmissionConf(path string, stopCh chan os.Signal) {
	dirPath := filepath.Dir(path)
	fileWatcher, err := filewatcher.NewFileWatcher(dirPath)
	if err != nil {
		hlog.Errorf("failed to create filewatcher for %s: %v", path, err)
		return
	}

	eventCh := fileWatcher.Events()
	errCh := fileWatcher.Errors()
	for {
		select {
		case event, ok := <-eventCh:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				hlog.Debugf("watch %s event: %v", dirPath, event)
				go LoadAdmissionConf(path)
				//go readFileByEnv(path)
			}
		case err, ok := <-errCh:
			if !ok {
				return
			}
			hlog.Errorf("watch %s error: %v", path, err)
		case <-stopCh:
			return
		}
	}
}


func LoadAdmissionConf(path string) {
	if path == "" {
		return
	}
	configBytes, err := os.ReadFile(path)
	if err != nil {
		hlog.Errorf("read admission file failed, err=%v", err)
		return
	}

	data := Config{}
	if err := yaml.Unmarshal(configBytes, &data); err != nil {
		hlog.Errorf("Unmarshal admission file failed, err=%v", err)
		return
	}
	EnvConfig.Lock()
	EnvConfig.Config = data
	EnvConfig.Unlock()
	hlog.Debugf("enfconfig: %v", EnvConfig.Config)
	if &EnvConfig.Config == nil {
		hlog.Errorf("load config failed")
	}
}
