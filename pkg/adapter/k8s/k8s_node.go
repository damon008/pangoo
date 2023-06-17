package k8s

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
)

type NodeService struct {
	ClientSet *kubernetes.Clientset
}

type Node struct {
	*v1.Node
	Request v1.ResourceList
}

type NodeList struct {
	metav1.TypeMeta
	metav1.ListMeta
	Items []Node
}

func NewNodeService() *NodeService {
	return &NodeService{ClientSet: pClientSet}
}

func PodRequests(pod *v1.Pod) (reqs map[v1.ResourceName]resource.Quantity, err error) {
	reqs = map[v1.ResourceName]resource.Quantity{}
	for _, container := range pod.Spec.Containers {
		for name, quantity := range container.Resources.Requests {
			if value, ok := reqs[name]; !ok {
				reqs[name] = quantity
			} else {
				value.Add(quantity)
				reqs[name] = value
			}
		}
	}
	// init containers define the minimum of any resource
	for _, container := range pod.Spec.InitContainers {
		for name, quantity := range container.Resources.Requests {
			value, ok := reqs[name]
			if !ok {
				reqs[name] = quantity
				continue
			}
			if quantity.Cmp(value) > 0 {
				reqs[name] = quantity
			}
		}
	}
	return
}

func getPodsTotalRequests(podList []v1.Pod) (reqs map[v1.ResourceName]resource.Quantity, err error) {
	reqs = map[v1.ResourceName]resource.Quantity{}
	for _, pod := range podList {
		podReqs, err := PodRequests(&pod)
		if err != nil {
			return nil, err
		}
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = podReqValue
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
	}
	return
}

func (ns *NodeService) clusterRequests() (map[string]v1.ResourceList, error) {
	selector := "spec.nodeName!=%22%22" +
		",status.phase!=" + string(v1.PodSucceeded) + ",status.phase!=" + string(v1.PodFailed)
	fieldSelector, err := fields.ParseSelector(selector)
	if err != nil {
		hlog.Error(err)
		return nil, err
	}

	nonTerminatedPodsList, err := ns.ClientSet.CoreV1().Pods("").List(context.TODO(),
		metav1.ListOptions{FieldSelector: fieldSelector.String()})
	if err != nil {
		hlog.Error(err)
		return nil, err
	}

	nodeNonTerminatedPodsList := map[string][]v1.Pod{}
	for _, pod := range nonTerminatedPodsList.Items {
		if _, ok := nodeNonTerminatedPodsList[pod.Spec.NodeName]; !ok {
			nodeNonTerminatedPodsList[pod.Spec.NodeName] = []v1.Pod{pod}
		} else {
			nodeNonTerminatedPodsList[pod.Spec.NodeName] = append(nodeNonTerminatedPodsList[pod.Spec.NodeName], pod)
		}
	}

	requestResourceList := map[string]v1.ResourceList{}
	for nodeName, nodeNonTerminatedPods := range nodeNonTerminatedPodsList {
		reqs, err := getPodsTotalRequests(nodeNonTerminatedPods)
		if err != nil {
			hlog.Error(err)
			return nil, err
		}
		requestResourceList[nodeName] = reqs
	}

	return requestResourceList, nil
}

func (ns *NodeService) unionResourceList(resourceListSrc v1.ResourceList, resourceListDest v1.ResourceList) {
	for podReqName, podReqValue := range resourceListSrc {
		if value, ok := resourceListDest[podReqName]; !ok {
			resourceListDest[podReqName] = podReqValue
		} else {
			value.Add(podReqValue)
			resourceListDest[podReqName] = value
		}
	}
}

func (ns *NodeService) requests(nodeName string) (v1.ResourceList, error) {
	selector := "spec.nodeName=" + nodeName +
		",status.phase!=" + string(v1.PodSucceeded) + ",status.phase!=" + string(v1.PodFailed)
	fieldSelector, err := fields.ParseSelector(selector)
	if err != nil {
		hlog.Error(err)
		return nil, err
	}

	nonTerminatedPodsList, err := ns.ClientSet.CoreV1().Pods("").List(context.TODO(),
		metav1.ListOptions{FieldSelector: fieldSelector.String()})
	if err != nil {
		hlog.Error(err)
		return nil, err
	}

	reqs, err := getPodsTotalRequests(nonTerminatedPodsList.Items)
	if err != nil {
		hlog.Error(err)
		return nil, err
	}

	return reqs, nil
}

func (ns *NodeService) Get(name string) (*Node, error) {
	nodeWithAvail := &Node{}
	node, err := ns.ClientSet.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	nodeWithAvail.Node = node
	nodeWithAvail.Request, err = ns.requests(name)

	return nodeWithAvail, err
}

func (ns *NodeService) List() (*NodeList, error) {
	nodes, err := ns.ClientSet.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		hlog.Error(err)
		return nil, err
	}

	nodeList := &NodeList{TypeMeta: nodes.TypeMeta, ListMeta: nodes.ListMeta}
	for index, node := range nodes.Items {
		req, err := ns.clusterRequests()
		if err != nil {
			hlog.Error(err)
			return nil, err
		}

		nodeList.Items = append(nodeList.Items, Node{Node: &nodes.Items[index], Request: req[node.Name]})
	}

	return nodeList, nil
}
