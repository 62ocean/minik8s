package kubectl

import (
	"encoding/json"
	"errors"
	"fmt"
	"k8s/object"
	"k8s/pkg/global"
	"k8s/pkg/util/HTTPClient"
	"k8s/pkg/util/parseYaml"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

var APIClient = HTTPClient.CreateHTTPClient(global.ServerHost)

func CmdExec() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create an object based on .yaml file",
				Subcommands: []*cli.Command{
					{
						Name:  "dns",
						Usage: "create a dns based on a dns.yaml",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "f",
								Usage:    "the path of the configuration file of a dns",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							fmt.Println("create: ", c.String("f"))
							filePath := c.String("f")
							newPod := parseYaml.ParseYaml[object.Dns](filePath)
							// id, _ := uuid.NewUUID()
							// newPod.Metadata.Uid = id.String()
							client := HTTPClient.CreateHTTPClient(global.ServerHost)
							dnsJson, _ := json.Marshal(newPod)
							fmt.Println(newPod.Metadata.Name)
							client.Post("/dns/create", dnsJson)
							return nil
						},
					},
					{
						Name:  "pod",
						Usage: "create a pod",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "f",
								Usage:    "the path of the configuration file of a pod",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							filePath := c.String("f")
							log.Println("create pod: ", c.String("f"))
							newPod := parseYaml.ParseYaml[object.Pod](filePath)
							podJson, _ := json.Marshal(newPod)
							log.Println(newPod)
							APIClient.Post("/pods/create", podJson)
							return nil
						},
					},
					{
						Name:  "service",
						Usage: "create a service",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "f",
								Usage:    "the path of the configuration file of a service",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							filePath := c.String("f")
							log.Println("create service: ", c.String("f"))
							newService := parseYaml.ParseYaml[object.Service](filePath)
							serviceJson, _ := json.Marshal(newService)
							log.Println(newService)
							APIClient.Post("/services/create", serviceJson)
							return nil
						},
					},
					{
						Name:  "RS",
						Usage: "get the running information of a replicaset",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "f",
								Usage:    "the path of the configuration file of a replicaset",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							filePath := c.String("f")
							log.Println("create RS: ", c.String("f"))
							newRS := parseYaml.ParseYaml[object.ReplicaSet](filePath)
							rsJson, _ := json.Marshal(newRS)
							log.Println(newRS)
							APIClient.Post("/replicasets/create", rsJson)
							return nil
						},
					},
					{
						Name:  "HPA",
						Usage: "get the running information of a HPA",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "f",
								Usage:    "the path of the configuration file of a HPA",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							filePath := c.String("f")
							log.Println("create HPA: ", c.String("f"))
							newHPA := parseYaml.ParseYaml[object.Hpa](filePath)
							HPAJson, _ := json.Marshal(newHPA)
							log.Println(newHPA)
							APIClient.Post("/hpas/create", HPAJson)
							return nil
						},
					},
					{
						Name:  "GPUJob",
						Usage: "get the running information of a service",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "f",
								Usage:    "the path of the configuration file of a pod",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							filePath := c.String("f")
							log.Println("create GPUJob: ", c.String("f"))
							// job存入apiserver
							job := parseYaml.ParseYaml[object.GPUJob](filePath)
							job.Status = object.PENDING
							jobInfo, _ := json.Marshal(job)
							APIClient.Post("/gpuJobs/create", jobInfo)

							// 构造pod 存入apiserver
							port := object.ContainerPort{Port: 8080}
							container := object.Container{
								Name:  "commit_" + "GPUJob_" + job.Metadata.Name,
								Image: "saltfishy/gpu_server:v8",
								Ports: []object.ContainerPort{
									port,
								},
								Command: []string{
									"./main ",
								},
								Args: []string{
									job.Metadata.Name,
								},
								// TODO 此处写入kubectl时需要修改为参数指定的文件路径
								CopyFile: job.Spec.Program,
								CopyDst:  "/apps",
							}
							newPod := object.Pod{
								ApiVersion: "v1",
								Kind:       "Pod",
								Metadata: object.Metadata{
									Name: "GPUJob_" + job.Metadata.Name,
									Labels: object.Labels{
										App: "GPU",
										Env: "prod",
									},
								},
								Spec: object.PodSpec{
									Containers: []object.Container{
										container,
									},
								},
							}
							podInfo, _ := json.Marshal(newPod)
							APIClient.Post("/pods/create", podInfo)
							return nil
						},
					},
				},
			},
			{
				// TODO
				Name:  "delete",
				Usage: "delete an object based on name",
				Subcommands: []*cli.Command{
					{
						Name:  "pod",
						Usage: "delete a pod",
						Action: func(c *cli.Context) error {
							if c.NArg() != 1 {
								return errors.New("the pod name must be specified")
							}
							name := c.Args().First()
							nameReq, _ := json.Marshal(name)
							APIClient.Post("/pods/remove", nameReq)
							return nil
						},
					},
					{
						Name:  "service",
						Usage: "delete a service",
						Action: func(c *cli.Context) error {
							if c.NArg() != 1 {
								return errors.New("the service name must be specified")
							}
							name := c.Args().First()
							nameReq, _ := json.Marshal(name)
							APIClient.Post("/services/remove", nameReq)
							return nil
						},
					},
					{
						Name:  "RS",
						Usage: "get the running information of a replicaset",
						Action: func(c *cli.Context) error {
							filePath := c.String("f")
							log.Println("delete RS: ", c.String("f"))
							newRS := parseYaml.ParseYaml[object.ReplicaSet](filePath)
							rsJson, _ := json.Marshal(newRS)
							log.Println(newRS)
							APIClient.Post("/replicasets/delete", rsJson)
							return nil
						},
					},
					{
						Name:  "HPA",
						Usage: "get the running information of a HPA",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "f",
								Usage:    "the path of the configuration file of a HPA",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							filePath := c.String("f")
							log.Println("delete HPA: ", c.String("f"))
							newHPA := parseYaml.ParseYaml[object.Hpa](filePath)
							HPAJson, _ := json.Marshal(newHPA)
							log.Println(newHPA)
							APIClient.Post("/hpas/delete", HPAJson)
							return nil
						},
					},
				},
			},
			{
				Name:  "get",
				Usage: "get the running information of a pod or a service",
				Subcommands: []*cli.Command{
					{
						Name:  "pod",
						Usage: "get the running information of a pod",
						Action: func(c *cli.Context) error {
							if c.NArg() != 1 {
								return errors.New("the pod name must be specified")
							}
							name := c.Args().First()
							podInfo := APIClient.Get("/pods/get/" + name)
							podStorage := object.PodStorage{}
							_ = json.Unmarshal([]byte(podInfo), &podStorage)
							fmt.Println("NAME\t\t\tSATUS\t\t\tAGE")
							createTime := podStorage.Config.Metadata.CreationTimestamp
							newtime := time.Now()
							d := newtime.Sub(createTime)
							fmt.Printf("%s\t\t\t%s\t\t\t%s\n", name, podStorage.Status.ToString(), d.Truncate(time.Second).String())
							return nil
						},
					},
					{
						Name:  "pods",
						Usage: "get the running information of all pod",
						Action: func(c *cli.Context) error {
							response := APIClient.Get("/pods/getAll")
							var podList map[string]string
							_ = json.Unmarshal([]byte(response), &podList)
							fmt.Println("NAME\t\t\tSATUS\t\t\tAGE")
							for _, val := range podList {
								podStorage := object.PodStorage{}
								_ = json.Unmarshal([]byte(val), &podStorage)
								createTime := podStorage.Config.Metadata.CreationTimestamp
								newtime := time.Now()
								d := newtime.Sub(createTime)
								fmt.Printf("%s\t\t\t%s\t\t\t%s\n", podStorage.Config.Metadata.Name, podStorage.Status.ToString(), d.Truncate(time.Second).String())

							}
							return nil
						},
					},
					{
						Name:  "GPUJob",
						Usage: "get the running information of a GPUJob",
						Action: func(c *cli.Context) error {
							if c.NArg() != 1 {
								return errors.New("the job name must be specified")
							}
							name := c.Args().First()
							jobInfo := APIClient.Get("/gpuJobs/get/" + name)
							job := object.GPUJob{}
							_ = json.Unmarshal([]byte(jobInfo), &job)
							fmt.Println("NAME\t\t\tSATUS\t\t\tAGE")
							createTime := job.Metadata.CreationTimestamp
							newtime := time.Now()
							d := newtime.Sub(createTime)
							fmt.Printf("%s\t\t\t%s\t\t\t%s\n", name, job.Status.ToString(), d.Truncate(time.Second).String())
							if job.Status == 3 {
								fmt.Printf("OUTPUT: \n")
								fmt.Println(job.Output)
							}
							return nil
						},
					},
					{
						Name:  "services",
						Usage: "get the running information of all services",
						Action: func(c *cli.Context) error {
							response := APIClient.Get("/services/getAll")
							var serviceList map[string]string
							_ = json.Unmarshal([]byte(response), &serviceList)
							fmt.Println("NAME\t\t\tCLUSTERIP\t\t\tLABEL")
							for _, val := range serviceList {
								service := object.Service{}
								_ = json.Unmarshal([]byte(val), &service)
								label := fmt.Sprint("app:%s env:%s", service.Metadata.Labels.App, service.Metadata.Labels.Env)
								fmt.Printf("%s\t\t\t%s\t\t\t%s\n", service.Metadata.Name, service.Spec.ClusterIP, label)

							}
							return nil
						},
					},
				},
			},

			{
				Name:  "describe",
				Usage: "get the detailed running information of a pod or a service",
				Subcommands: []*cli.Command{
					{
						Name:  "pod",
						Usage: "get the running information of a pod",
						Action: func(c *cli.Context) error {
							if c.NArg() != 1 {
								return errors.New("the pod name must be specified")
							}
							podName := c.Args().First()
							podInfo := APIClient.Get("/pods/get/" + podName)
							podStorage := object.PodStorage{}
							_ = json.Unmarshal([]byte(podInfo), &podStorage)
							fmt.Println(podStorage)
							yamlData, err := yaml.Marshal(podStorage)
							if err != nil {
								fmt.Println("转换为 YAML 失败:", err)
								return nil
							}
							fmt.Println(string(yamlData))
							return nil
						},
					},
					{
						Name:  "service",
						Usage: "get the running information of a service",
						Action: func(c *cli.Context) error {
							if c.NArg() != 1 {
								return errors.New("the service name must be specified")
							}
							fmt.Println("service information: ", c.Args().First())
							serviceName := c.Args().First()

							getMsg, _ := json.Marshal(serviceName)
							resp := APIClient.Post("/services/get", getMsg)
							// fmt.Printf("Kubectl get service: %s\n", resp)
							service := object.Service{}
							var svStr string
							json.Unmarshal([]byte(resp), &svStr)
							json.Unmarshal([]byte(svStr), &service)

							yamlData, err := yaml.Marshal(service)
							if err != nil {
								fmt.Println("转换为 YAML 失败:", err)
								return nil
							}
							fmt.Println(string(yamlData))

							// fmt.Printf("Name: %s\n", service.Metadata.Name)
							// fmt.Printf("Type: %s\n", service.Spec.Type)
							// fmt.Printf("Selector:\n")
							// fmt.Printf("\tapp: %s\n", service.Spec.Selector.App)
							// fmt.Printf("\tenv: %s\n", service.Spec.Selector.Env)
							// fmt.Printf("ClusterIP: %s\n", service.Spec.ClusterIP)
							return nil
						},
					},
					{
						Name:  "dns",
						Usage: "get the information of dns and path",
						Action: func(c *cli.Context) error {
							resp := APIClient.Get("/dns/get")
							dns := object.Dns{}
							json.Unmarshal([]byte(resp), &dns)
							fmt.Println(dns)
							yamlData, err := yaml.Marshal(dns)
							if err != nil {
								fmt.Println("转换为 YAML 失败:", err)
								return nil
							}
							fmt.Println(string(yamlData))
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
