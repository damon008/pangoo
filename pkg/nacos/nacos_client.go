package nacos

import (
	//"github.com/hertz-contrib/registry/nacos/common"
	//"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"primus/pkg/util"

	//"github.com/natefinch/lumberjack"
	"os"
)

func NewNacosConfig(ip string, port int32, logPath string) (naming_client.INamingClient, error) {
	sc := []constant.ServerConfig{
		//*constant.NewServerConfig(os.Getenv("serverAddr"), uint64(NacosPort())),
		*constant.NewServerConfig(ip, uint64(port)),
	}

	cc := constant.ClientConfig{
		NamespaceId:         os.Getenv("namespace"),
		RegionId:            "cn-hangzhou",
		CustomLogger:        util.NewHLogger("/data/nacos/logs/"),
		NotLoadCacheAtStart: true,
		//LogDir:              "/data/nacos/logs",
		CacheDir: "/data/nacos/cache",
		LogLevel: "warn",
		//LogRollingConfig: &lumberjack.Logger{},
		//Username:            "your-name",
		//Password:            "your-password",
	}
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

/*func NewCustomNacosLogger() logger.Logger {

}*/

func InitNacos(ip string, port uint64) (naming_client.INamingClient, error) {
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ip, port),
	}

	cc := constant.ClientConfig{
		NamespaceId:         "public",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/data/nacos/log",
		CacheDir:            "/data/nacos/cache",
		LogLevel:            "info",
		//Username:            "your-name",
		//Password:            "your-password",
	}

	cli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
		return nil, err
	}
	return cli, err
}

/*func NacosPort() int64 {
	portText := os.Getenv(NACOS_ENV_PORT)
	if len(portText) == 0 {
		return NACOS_DEFAULT_PORT
	}
	port, err := strconv.ParseInt(portText, 10, 64)
	if err != nil {
		klog.Errorf("ParseInt failed,err:%s", err.Error())
		return NACOS_DEFAULT_PORT
	}
	return port
}

// NacosAddr Get Nacos addr from environment variables
func NacosAddr() string {
	addr := os.Getenv(NACOS_ENV_SERVER_ADDR)
	if len(addr) == 0 {
		return NACOS_DEFAULT_SERVER_ADDR
	}
	return addr
}*/
