package main

import (
	pb "config/proto/output/proto/config"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
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
		pathWebserver := "sample"
		router := gin.Default()
		gin.ForceConsoleColor()
		router.Use(Cors())
		router.LoadHTMLGlob("webserver/" + pathWebserver + "/templates/*")
		router.StaticFile("/favicon.ico", "webserver/"+pathWebserver+"/static/ico/favicon.ico")

		routerConfig := router.Group("/")
		{
			routerConfig.GET("/config", func(c *gin.Context) {
				byteConfData, err := json.Marshal(conf.Data)
				if err != nil {
					log.Fatalln(err)
				}
				var mapConfData map[string]string
				err = json.Unmarshal(byteConfData, &mapConfData)
				if err != nil {
					log.Fatalln(err)
				}
				c.HTML(http.StatusOK, "ui.tmpl", gin.H{
					"dynamic": string(byteConfData),
					"table":   mapConfData,
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
		log.Println("wait a minute, you can access http://localhost:8081/config")
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method               //请求方法
		origin := c.Request.Header.Get("Origin") //请求头部
		var headerKeys []string                  // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")                                       // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE") //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar") // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")                                                                                                                                                           // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")                                                                                                                                                  //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")                                                                                                                                                              // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next() //  处理请求
	}
}
