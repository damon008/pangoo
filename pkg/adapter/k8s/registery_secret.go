package k8s

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	RegistrySecretName     = "registry-secret"
	PrivateRegistryProject = "private"
	PublicRegistryProject  = "public"
)

func IsSupportPrivateRegistryProject() bool {
	return true//beego.AppConfig.DefaultString("registry::type", PublicRegistryProject) == PrivateRegistryProject
}

func CreateRegistrySecret(namespace string) error {
	secret, err := NewSecret().Get(namespace, RegistrySecretName)
	if err != nil && !errors.IsNotFound(err) {
		hlog.Error("query register secret fail, error=%s", err.Error())
		return err
	}

	if err == nil && secret != nil {
		hlog.Info("registry secret %s exist", RegistrySecretName)
		return nil
	} else {
		hlog.Info("registry secret %s not exist,create it", RegistrySecretName)
	}
	registryName := ""//beego.AppConfig.String("registry::name")
	//userName := ""//beego.AppConfig.String("registry::userName")
	//password := ""//beego.AppConfig.String("registry::password")
	auth := ""//util.Base64Encode(fmt.Sprintf("%s:%s", userName, password))
	value := fmt.Sprintf(`{"auths":{"%s":{"auth":"%s"}}}`, registryName, string(auth))
	data := map[string][]byte{
		v1.DockerConfigJsonKey: []byte(value),
	}
	secretPara := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: RegistrySecretName,
		},
		Data: data,
		Type: v1.SecretTypeDockerConfigJson,
	}

	_, err = NewSecret().CreateSecret(namespace, secretPara)
	if err != nil {
		hlog.Error("Create Secret Failed %v", err)
	}

	hlog.Info("create %s success", RegistrySecretName)
	return err
}

func ImagePullSecrets() []v1.LocalObjectReference {
	if IsSupportPrivateRegistryProject() {
		return []v1.LocalObjectReference{
			v1.LocalObjectReference{
				Name: RegistrySecretName,
			},
		}
	} else {
		return []v1.LocalObjectReference{}
	}
}
