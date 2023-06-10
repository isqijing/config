package main

import (
	pb "config/proto/output/proto/sample_750648506"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}



type Data struct {
	IsSpace	string `json:"is_space"`
	IsLocation	string `json:"is_location"`
	IsCar	string `json:"is_car"`

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
	log.Println("finally config: ", conf.Data)
	//for _, v := range conf.Metadata.Change {
	//	log.Println(v)
	//}

	wg.Wait()
}


type Config struct {
	Metadata Metadata `json:"metadata"`
	Data     Data     `json:"data"`
	Errors   []error  `json:"errors"`
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


func (detail *Config) ReadByInput() {

	var isLocation string
	flag.StringVar(&isLocation, "is_location", detail.Data.IsLocation, "usage unimplement")
	detail.Metadata.Change = append(detail.Metadata.Change, "input", isLocation)
	var isCar string
	flag.StringVar(&isCar, "is_car", detail.Data.IsCar, "usage unimplement")
	detail.Metadata.Change = append(detail.Metadata.Change, "input", isCar)
	var isSpace string
	flag.StringVar(&isSpace, "is_space", detail.Data.IsSpace, "usage unimplement")
	detail.Metadata.Change = append(detail.Metadata.Change, "input", isSpace)
	detail.Data.IsLocation = isLocation
	detail.Data.IsCar = isCar
	detail.Data.IsSpace = isSpace



}

func (detail *Config) ReadByFile() {

	file, err := os.ReadFile("config/sample_750648506_config.json")
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

detail.Metadata.Change = append(detail.Metadata.Change, "env")
	isLocation := os.Getenv("is_location")
	if isLocation != "" {
		detail.Metadata.Change = append(detail.Metadata.Change, "is_location", detail.Data.IsLocation, isLocation)
		detail.Data.IsLocation = isLocation
	}
	isCar := os.Getenv("is_car")
	if isCar != "" {
		detail.Metadata.Change = append(detail.Metadata.Change, "is_car", detail.Data.IsCar, isCar)
		detail.Data.IsCar = isCar
	}
	isSpace := os.Getenv("is_space")
	if isSpace != "" {
		detail.Metadata.Change = append(detail.Metadata.Change, "is_space", detail.Data.IsSpace, isSpace)
		detail.Data.IsSpace = isSpace
	}


}

func (detail *Config) ReadByGrpc() {

	credentials := insecure.NewCredentials()
	conn, _ := grpc.Dial("127.0.0.1:2023", grpc.WithTransportCredentials(credentials), grpc.WithBlock())
	defer conn.Close()

	client := pb.NewConfigClient(conn)
	resp, err := client.ReadConfig(context.Background(), &pb.RequestReadConfig{
		FromId:       "sample_750648506__from_id",
		FromNickname: "sample_750648506__from_nickname",
		Content:      "sample_750648506__content",
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

detail.Data.IsSpace = "true"
detail.Data.IsLocation = "false"
detail.Data.IsCar = "false"


}
