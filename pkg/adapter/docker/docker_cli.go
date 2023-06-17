package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io/ioutil"
	"sync"

	//"github.com/docker/docker/api/types/container"
	//"github.com/docker/docker/api/types/mount"
	//"github.com/docker/docker/api/types/network"
	//"io/ioutil"

	"errors"
	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/golang/glog"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

const dockerSock = "unix:///var/run/docker.sock"

var (
	dCli *client.Client
	authStr string
	doOnce sync.Once
	instance atomic.Value
	//dCli *DockerAgent
)

func NewDockerAgent(username, password string) *DockerAgent {
	doOnce.Do(func() {
		if str,err := Auth(username, password); err ==nil {
			authStr = str
			instance.Store(&DockerAgent{
				Client:  dCli,
				AuthStr: str,
			})
		}
		instance.Store(&DockerAgent{
			Client:  dCli,
			AuthStr: "",
		})
	})
	return instance.Load().(*DockerAgent)
}

type DockerAgent struct {
	Client *client.Client
	AuthStr string
}

func NewClient() *client.Client {
	var err error
	doOnce.Do(func() {
		dCli, err = client.NewClientWithOpts(client.WithHost(dockerSock), client.WithAPIVersionNegotiation())
		if err != nil {
			hlog.Error("docker connect failed, %v", err)
			dCli = nil
		}
		//dinstance.Store(&dCli)
		/*dockerCli = cli
		dCli = &DockerClient{
			Client: dockerCli,
		}*/
	})
	return dCli
	//return instance.Load().(*client.Client)
}

func Auth(username, password string) (string, error) {
	authData := types.AuthConfig{
		Username: username,
		Password: password,
		//ServerAddress: "https://index.docker.io/v1/",
	}
	b, err := sonic.Marshal(authData)
	if err != nil {
		glog.Errorf("auth json %v marshal failed, %v", authData, err)
		return "", err
	}
	authStr = base64.StdEncoding.EncodeToString(b)
	return authStr, nil
}

func (cli *DockerAgent) ImageBuild(buildFilePath,tag string) {
	hlog.Info("buildFilePath: ", buildFilePath, "tag: ", tag)
	ctx := context.Background()
	//在 types.ImageBuildOptions 中，这两个属性的组合非常重要。Context 属性确定了构建上下文，
	//而 Dockerfile 属性则确定了在该上下文中使用哪个 Dockerfile 文件进行构建。通过这种方式，
	//可以轻松地在不同的上下文中使用不同的 Dockerfile 来构建不同的镜像。
	//
	//例如，如果想在当前目录下的 Dockerfile 文件中构建镜像，则可以将 Context 设置为 .，
	//而将 Dockerfile 设置为 Dockerfile。如果要在其他目录中构建，则可以将 Context 设置为该目录的路径，
	//而将 Dockerfile 设置为该目录中的 Dockerfile 文件的名称。
	//iob, err := util.ReadFile(buildFilePath)
	iob, err := os.Open(buildFilePath)

	if err != nil {
		hlog.Error(err)
		return
	}
	defer iob.Close()
	buildOptions := types.ImageBuildOptions{
		Tags:           []string{tag},
		Dockerfile:     "Dockerfile",//指定 Dockerfile 文件名
		Context:        iob, //设置上下文路径io.
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		NoCache:        true,
		SuppressOutput: false,
	}

	buildResp, err := dCli.ImageBuild(ctx, os.Stdin, buildOptions)
	if err != nil {
		hlog.Error(err)
	}
	// 读取构建输出
	defer buildResp.Body.Close()
	io.Copy(os.Stdout, buildResp.Body)

	// 读取构建输出
	defer buildResp.Body.Close()
	out, err := ioutil.ReadAll(buildResp.Body) // 必须读取完整个响应体才能关闭连接
	if err != nil {
		hlog.Error(err)
	}
	fmt.Println(string(out))

	// 读取构建输出
	imageBuildOutput := make([]byte, 1024)
	for {
		_, err := buildResp.Body.Read(imageBuildOutput)
		if err != nil {
			break
		}
		fmt.Print(string(imageBuildOutput))
	}

	// 读取构建输出
	scanner := bufio.NewScanner(buildResp.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// 推送镜像到 Docker 仓库
	//pushResp, err := cli.Client.ImagePush(ctx, tag, types.ImagePushOptions{RegistryAuth: auth})
	//if err != nil {
	//	hlog.Error(err)
	//}
	//defer pushResp.Close()
	//io.Copy(os.Stdout, pushResp)

	/*
	// Create a new container using the built image
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{"echo", "hello world"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/path/on/host",
				Target: "/path/in/container",
			},
		},
	}, &network.NetworkingConfig{}, "")
	if err != nil {
		hlog.Error(err)
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		hlog.Error(err)
	}

	// Wait for the container to finish running
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			hlog.Error(err)
		}
	case <-statusCh:
	}

	// Get stdout from container
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		hlog.Error(err)
	}
	defer out.Close()

	// Print stdout from container
	stdout := make([]byte, 1024)
	for {
		_, err := out.Read(stdout)
		if err != nil {
			break
		}
		fmt.Print(string(stdout))
	}*/
}

func (cli *DockerAgent) PullImage(repo string, auth string) error {
	if len(auth) == 0 {
		auth = authStr
	}
	resp, err := dCli.ImagePull(context.Background(), repo, types.ImagePullOptions{RegistryAuth: auth})
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

func (cli *DockerAgent) PushImage(repo string, auth string) error {
	if len(auth) == 0 {
		auth = authStr
	}
	resp, err := dCli.ImagePush(context.Background(), repo, types.ImagePushOptions{RegistryAuth: auth})
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

func (cli *DockerAgent) DeleteImage(repo string) error {
	_, err := cli.Client.ImageRemove(context.Background(), repo, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (cli *DockerAgent) GetImageInfo(image string) (types.ImageInspect, error) {
	response, _, err := cli.Client.ImageInspectWithRaw(context.Background(), image)
	if err != nil {
		return types.ImageInspect{}, err
	}
	return response, nil
}

func (cli *DockerAgent) ListContainers() ([]types.Container, error) {
	response, err := cli.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return []types.Container{}, err
	}
	return response, nil
}

// container: [podName]_[nameSpace]_[UID]
func (cli *DockerAgent) GetContainerLikeName(container string) (types.Container, error) {
	containers, err := cli.ListContainers()
	if err != nil {
		return types.Container{}, err
	}
	for _, c := range containers {
		//k8s_[containerName]_[podName]_[nameSpace]_[UID]_[index]
		if strings.Contains(c.Names[0], container) {
			// 排除pause容器
			hlog.Debug("the c.Names[0] is :",c.Names[0])
			hlog.Debug("the container is :",container)
			if !strings.Contains(c.Image,"/pause:"){
				return c, nil
			}
		}
	}
	return types.Container{}, errors.New("container not found")
}

func (cli *DockerAgent) GetContainerInfo(container string) (types.ContainerJSON, error) {
	response, err := cli.Client.ContainerInspect(context.Background(), container)
	if err != nil {
		return types.ContainerJSON{}, err
	}
	return response, nil
}

func (cli *DockerAgent) CommitImage(container string, image string) error {
	option := types.ContainerCommitOptions{
		Reference: image,
	}
	_, err := cli.Client.ContainerCommit(context.Background(), container, option)
	if err != nil {
		return err
	}
	return nil
}

func (cli *DockerAgent) GetContainerUpperDirSize(container string) (int64, error) {
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
