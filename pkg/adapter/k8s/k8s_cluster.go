package k8s

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"primus/pkg/constants"
	"strings"
)

type Cluster struct {
	Nodes *NodeList
}

func (c *Cluster) Load() error {
	nodeList, err := NewNodeService().List()
	if err != nil {
		hlog.Error(err)
		return err
	}

	c.Nodes = nodeList
	return nil
}

func (c *Cluster) GetCardWithVendor(cardType string) string {
	for _, node := range c.Nodes.Items {
		for resourceType, v := range node.Labels {
			if strings.HasPrefix(resourceType, constants.LabelNodeResourceType+".") && v != "false" && strings.HasSuffix(resourceType, "/"+cardType) {
				return strings.Replace(resourceType, constants.LabelNodeResourceType+".", "", -1)
			}
		}
	}

	return ""
}

func (c *Cluster) Exist(cardType string) bool {
	return c.GetCardWithVendor(cardType) != ""
}

func (c *Cluster) GetNodeAccDeviceType(node *Node) (string, bool) {
	for resourceType, v := range node.Labels {
		if strings.HasPrefix(resourceType, constants.LabelNodeResourceType+".") && strings.Contains(resourceType, "/") && v != "false" {
			resourceTypeWithVendor := strings.Replace(resourceType, constants.LabelNodeResourceType+".", "", -1)
			return strings.Split(resourceTypeWithVendor, "/")[1], true
		}
	}
	return "", false
}
