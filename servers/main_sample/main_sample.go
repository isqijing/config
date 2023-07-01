package main

import (
	"config/middlewares"
	pb "config/proto/output/proto/config"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

type Config struct {
	Metadata Metadata `json:"metadata"`
	Data     Data     `json:"data"`
	Errors   []error  `json:"errors"`
}

type Data struct {
	IsCompany string `json:"is_company"` // 在公司
	IsHome    string `json:"is_home"`    // 在家
}
type Metadata struct {
	Change Change `json:"change"`
}
type Change []string

type ConfigServer struct {
	pb.UnimplementedConfigServer
}

func (s *ConfigServer) ReadConfig(ctx context.Context, request *pb.RequestReadConfig) (*pb.ReplyReadConfig, error) {
	//log.Println(request.FromId)
	//log.Println(request.FromNickname)
	//log.Println(request.Content)
	return &pb.ReplyReadConfig{
		Code:    200,
		Content: "{\"is_home\":\"grpc:true\",\"is_company\":\"grpc:false\"}",
		Count:   1,
		Msg: &pb.Msg{
			Success: "success",
			Fail:    "",
		},
	}, nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		server := grpc.NewServer()
		pb.RegisterConfigServer(server, &ConfigServer{})
		lis, _ := net.Listen("tcp", "127.0.0.1:2023")
		go func() {
			time.Sleep(time.Second * 2)
			wg.Done()
		}()
		server.Serve(lis)
	}()
	time.Sleep(time.Second * 1)
	conf := Config{}
	conf.ReadByDefault()
	conf.ReadByGrpc()
	conf.ReadByEnv()
	conf.ReadByFile()
	conf.ReadByInput()

	if len(conf.Errors) > 0 {
		for k, v := range conf.Errors {
			log.Println("error", k, ":", v.Error())
		}
		panic("ERROR: PLEASE fix above errors")
	}
	log.Println("config: ", conf)
	//for _, v := range conf.Metadata.Change {
	//	log.Println(v)
	//}
	wg.Add(1)
	go func() {
		defer wg.Done()
		router := gin.Default()
		gin.ForceConsoleColor()
		router.Use(middlewares.Cors())
		router.LoadHTMLGlob("webserver/templates/*")
		router.StaticFile("/favicon.ico", "webserver/static/ico/favicon.ico")

		routerConfig := router.Group("/")
		{
			routerConfig.GET("/config", func(c *gin.Context) {
				c.HTML(http.StatusOK, "ui.tmpl", gin.H{
					"title": "七镜配置",
				})
			})
			routerConfig.POST("/config/update", func(c *gin.Context) {
				var data Data
				err := c.BindJSON(&data)
				if err != nil {
					log.Println(err)
					c.JSON(http.StatusOK, err)
					return
				}
				conf.Data = data
				byteConfig, err := json.Marshal(conf.Data)
				if err != nil {
					c.JSON(http.StatusOK, err)
					return
				}
				err = os.WriteFile("config/config_config.json", byteConfig, 0600)
				if err != nil {
					c.JSON(http.StatusOK, err)
					return
				}

				c.JSON(http.StatusOK, map[string]string{"msg": "配置更新完毕，已写入配置文件，重启应用即可生效"})

			})
		}
		router.Run(":8081")
	}()
	wg.Wait()
}

func (detail *Config) ReadByUi() {

}

func (detail *Config) ReadByInput() {

	var isHome string
	flag.StringVar(&isHome, "is_home", detail.Data.IsHome, "location is home? true|false")
	var isCompany string
	flag.StringVar(&isCompany, "is_company", detail.Data.IsCompany, "location is company? true|false")
	flag.Parse()
	detail.Metadata.Change = append(detail.Metadata.Change, "input", isHome, isCompany)
	detail.Data.IsHome = isHome
	detail.Data.IsCompany = isCompany
}

func (detail *Config) ReadByFile() {
	file, err := os.ReadFile("config/config_config.json")
	if err != nil {
		detail.Errors = append(detail.Errors, err)
	}
	var data Data
	err = json.Unmarshal(file, &data)
	if err != nil {
		detail.Errors = append(detail.Errors, err)
	}
	detail.Metadata.Change = append(detail.Metadata.Change, "file", string(file))
}

func (detail *Config) ReadByEnv() {
	isHome := os.Getenv("is_home")
	detail.Metadata.Change = append(detail.Metadata.Change, "env")
	if isHome != "" {
		detail.Metadata.Change = append(detail.Metadata.Change, "is_home", detail.Data.IsHome, isHome)
		detail.Data.IsHome = isHome
	}
	isCompany := os.Getenv("is_company")
	if isCompany != "" {
		detail.Metadata.Change = append(detail.Metadata.Change, "is_company", detail.Data.IsCompany, isCompany)
		detail.Data.IsCompany = isCompany
	}

}

func (detail *Config) ReadByGrpc() {
	// default server host:port is localhost:2018
	credentials := insecure.NewCredentials()
	conn, _ := grpc.Dial("127.0.0.1:2023", grpc.WithTransportCredentials(credentials), grpc.WithBlock())
	defer conn.Close()

	client := pb.NewConfigClient(conn)
	resp, err := client.ReadConfig(context.Background(), &pb.RequestReadConfig{
		FromId:       "config",
		FromNickname: "config",
		Content:      "config",
	})

	if err != nil {
		detail.Errors = append(detail.Errors, err)
	}
	if resp.Code != 200 {
		detail.Errors = append(detail.Errors, errors.New(resp.Msg.Fail))
	}
	var data Data
	err = json.Unmarshal([]byte(resp.Content), &data)
	if err != nil {
		detail.Errors = append(detail.Errors, err)
	}
	detail.Metadata.Change = append(detail.Metadata.Change, "grpc", resp.Content)
}

func (detail *Config) ReadByDefault() {
	detail.Data.IsHome = "true"
	detail.Data.IsCompany = "false"
}
