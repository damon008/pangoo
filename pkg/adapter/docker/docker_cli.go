package docker

import (
	"context"
	"encoding/base64"

	"errors"
	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/golang/glog"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const dockerSock = "unix:///var/run/docker.sock"

type DockerClient struct {
	Client *client.Client
}

func NewClient() *DockerClient {
	cli, err := client.NewClientWithOpts(client.WithHost(dockerSock), client.WithAPIVersionNegotiation())
	if err != nil {
		glog.Errorf("docker connect failed, %v", err)
	}
	return &DockerClient{
		Client: cli,
	}
}

func (cli *DockerClient) Auth(username, password string) (string, error) {
	authData := types.AuthConfig{
		Username: username,
		Password: password,
	}
	b, err := sonic.Marshal(authData)
	if err != nil {
		glog.Errorf("auth json %v marshal failed, %v", authData, err)
		return "", err
	}
	encodeStr := base64.StdEncoding.EncodeToString(b)
	return encodeStr, nil
}

func (cli DockerClient) PullImage(image string, auth string) error {
	resp, err := cli.Client.ImagePull(context.Background(), image, types.ImagePullOptions{RegistryAuth: auth})
	if resp != nil {
		defer resp.Close()
		_, err := io.Copy(os.Stdout, resp)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (cli DockerClient) PushImage(image string, auth string) error {
	resp, err := cli.Client.ImagePush(context.Background(), image, types.ImagePushOptions{RegistryAuth: auth})
	if resp != nil {
		defer resp.Close()
		_, err := io.Copy(os.Stdout, resp)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (cli DockerClient) DeleteImage(image string) error {
	_, err := cli.Client.ImageRemove(context.Background(), image, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (cli DockerClient) GetImageInfo(image string) (types.ImageInspect, error) {
	response, _, err := cli.Client.ImageInspectWithRaw(context.Background(), image)
	if err != nil {
		return types.ImageInspect{}, err
	}
	return response, nil
}

func (cli DockerClient) ListContainers() ([]types.Container, error) {
	response, err := cli.Client.ContainerList(context.Background(), types.ContainerListOptions{Quiet: false})
	if err != nil {
		return []types.Container{}, err
	}
	return response, nil
}

// container: [podName]_[nameSpace]_[UID]
func (cli DockerClient) GetContainerLikeName(container string) (types.Container, error) {
	containers, err := cli.ListContainers()
	if err != nil {
		return types.Container{}, err
	}
	for _, c := range containers {
		//k8s_[containerName]_[podName]_[nameSpace]_[UID]_[index]
		if strings.Contains(c.Names[0], container) {
			// 排除pause容器
			logs.Debug("the c.Names[0] is :",c.Names[0])
			logs.Debug("the container is :",container)
			if !strings.Contains(c.Image,"/pause:"){
				return c, nil
			}
		}
	}
	return types.Container{}, errors.New("container not found")
}

func (cli DockerClient) GetContainerInfo(container string) (types.ContainerJSON, error) {
	response, err := cli.Client.ContainerInspect(context.Background(), container)
	if err != nil {
		return types.ContainerJSON{}, err
	}
	return response, nil
}

func (cli DockerClient) CommitImage(container string, image string) error {
	option := types.ContainerCommitOptions{
		Reference: image,
	}
	_, err := cli.Client.ContainerCommit(context.Background(), container, option)
	if err != nil {
		return err
	}
	return nil
}

func (cli DockerClient) GetContainerUpperDirSize(container string) (int64, error) {
	response, err := cli.GetContainerInfo(container)
	if err != nil {
		return 0, err
	}
	size, err := dirSizeB(response.GraphDriver.Data["UpperDir"])
	if err != nil {
		return 0, err
	}
	return size, nil
}

//getFileSize get file size by path(B)
func dirSizeB(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
