package license

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"primus/pkg/adapter/k8s"
)

func GetClusterNodes() (int, error) {
	nodeService := k8s.NewNodeService()
	list, err := nodeService.List()
	if err != nil {
		hlog.Error("get the node list failed: ", err)
		return 0, err
	}
	return len(list.Items), nil
}
