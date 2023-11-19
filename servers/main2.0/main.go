package main

import (
	"config/utils/copy_dir"
	"config/utils/welcome"
	"encoding/json"
	"flag"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	flags := log.Flags()
	flags &^= log.Ldate | log.Ltime
	log.SetFlags(flags)
}
func main() {

	welcome.Welcome()

	rand.Seed(time.Now().UnixNano())
	nameProto := "sample_" + strconv.Itoa(rand.Intn(1000000000))

	var pathProject string

	// 项目目录名称
	pathProject = "config"
	flag.StringVar(&pathProject, "path_project", pathProject, "project name")
	flag.StringVar(&pathProject, "pp", pathProject, "project name")

	// 组件目录名称
	var pathModules string
	pathModules = "output2.0/modules"
	flag.StringVar(&pathModules, "path_modules", pathModules, "modules name")
	flag.StringVar(&pathModules, "pm", pathModules, "modules name")

	// config导入之后项目的别名
	var pathConfigAlias string
	pathConfigAlias = "my_config"
	flag.StringVar(&pathConfigAlias, "path_config", pathConfigAlias, "config alias name")
	flag.StringVar(&pathConfigAlias, "pc", pathConfigAlias, "config alias name")

	// proto 目录
	var pathProtoTmp string
	pathProtoTmp = "proto"

	flag.StringVar(&nameProto, "name_proto", nameProto, "name which you like. final example: ${name}.proto")
	flag.StringVar(&nameProto, "np", nameProto, "name which you like. final example: ${name}.proto")

	var pathWebserver string
	pathWebserver = "sample_" + strconv.Itoa(rand.Intn(1000000000))
	flag.StringVar(&pathWebserver, "path_webserver", pathWebserver, "webserver'path which you like. final example: qijing_config; The directory of top level MUST be your project name")
	flag.StringVar(&pathWebserver, "pw", pathWebserver, "webserver'path which you like. final example: qijing_config; The directory of top level MUST be your project name")

	var pathOutput string
	flag.StringVar(&pathOutput, "path_output", "", "release path")
	flag.StringVar(&pathOutput, "po", "", "release path")

	flag.Parse()

	pathRootWebserver := pathModules + "/" + pathConfigAlias + "/webserver"

	open, err := os.ReadFile("dynamic.json")
	if err != nil {
		log.Fatalln(err)
	}
	var dynamic map[string]string
	err = json.Unmarshal(open, &dynamic)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Println(dynamic)

	log.Println(pathModules + "/" + pathConfigAlias + "/" + pathProtoTmp)
	err = os.MkdirAll(pathModules+"/"+pathConfigAlias+"/"+pathProtoTmp, 0600)
	if err != nil {
		log.Fatalln(err)
	}

	pathProto := pathModules + "/" + pathConfigAlias + "/" + pathProtoTmp + "/" + nameProto
	err = os.MkdirAll(pathProto, 0600)
	if err != nil {
		log.Fatalln(err)
	}

	pathConfig := pathModules + "/" + pathConfigAlias + "/" + "config"

	err = os.MkdirAll(pathConfig, 0600)
	if err != nil {
		log.Fatalln(err)
	}

	pathFunc := "index"

	err = os.MkdirAll(pathModules+"/"+pathConfigAlias+"/"+pathFunc, 0600)
	if err != nil {
		log.Fatalln(err)
	}

	fileProto, err := os.Create(pathProto + "/" + nameProto + ".proto")
	if err != nil {
		log.Fatalln(err)
	}
	//defer fileProto.Close()
	_, err = fileProto.WriteString(`syntax = "proto3";

option go_package = "proto/` + nameProto + `";

service Config {
  rpc ReadConfig (RequestReadConfig) returns (ReplyReadConfig) {};
}

message RequestReadConfig {
  string from_id = 1;
  string from_nickname = 2;
  string content = 3;
}

message ReplyReadConfig {
  uint32 code = 1;
  string content = 2;
  uint64 count = 3;
  Msg msg = 4;
}

message Msg {
  string success = 1;
  string fail = 2;
}`)
	if err != nil {
		return
	}
	fileProto.Close()

	// 生成f_main_config.go
	byteTemplate, err := os.ReadFile("template2.0.txt")
	if err != nil {
		log.Fatalln(err)
	}
	//log.Println(string(byteTemplate))

	pathOutputIndex := pathModules + "/" + pathConfigAlias + "/" + pathFunc
	//err = os.MkdirAll(pathOutputMain, 0600)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	fileIndex, err := os.Create(pathOutputIndex + "/index_" + nameProto + ".go")
	if err != nil {
		log.Fatalln(err)
	}

	pathMainQuickStart := pathModules + "/" + pathConfigAlias
	//err = os.MkdirAll(pathOutputMain, 0600)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	fileMainQuickStart, err := os.Create(pathMainQuickStart + "/main_quick_start.go")
	if err != nil {
		log.Fatalln(err)
	}

	template := string(byteTemplate)
	template = strings.ReplaceAll(template, "__PATH_CONFIG_ALIAS__", pathFunc)
	template = strings.ReplaceAll(template, "__PATH_PROJECT_PROTO__", pathProject+"/"+pathModules+"/"+pathConfigAlias+"/"+pathProtoTmp)
	template = strings.ReplaceAll(template, "__NAME_PROTO__", nameProto)
	template = strings.ReplaceAll(template, "__PATH_WEBSERVER__", pathRootWebserver+"/"+pathWebserver)
	struct_data := ""
	for k, _ := range dynamic {
		struct_data += "\t" + underscoreToCamel(k) + "\tstring `json:\"" + k + "\"`\n"

	}
	template = strings.ReplaceAll(template, "__STRUCT_DATA__", struct_data)

	read_by_input := ""
	var read_by_input_assign []string
	for k, _ := range dynamic {
		read_by_input += "\tvar " + underscoreToLowerCamel(k) + " string\n"
		read_by_input += "\tflag.StringVar(&" + underscoreToLowerCamel(k) + ", \"" + k + "\", detail.Data." + underscoreToCamel(k) + ", \"" + "usage unimplement" + "\")\n"
		read_by_input += "\tflag.Parse()\n"

		read_by_input += "\tdetail.Metadata.Change = append(detail.Metadata.Change, \"input\", " + underscoreToLowerCamel(k) + ")\n"
		read_by_input_assign = append(read_by_input_assign, "\tdetail.Data."+underscoreToCamel(k)+" = "+underscoreToLowerCamel(k)+"\n")
	}
	read_by_input += strings.Join(read_by_input_assign, "") + "\n"
	template = strings.ReplaceAll(template, "__READ_BY_INPUT__", read_by_input)

	read_by_file := `	var reqReadByFile ReqReadByFile
	switch interfaceReqReadByFile.(type) {
	case ReqReadByFile:
		reqReadByFile = interfaceReqReadByFile.(ReqReadByFile)
	default:
		log.Fatalln(reqReadByFile)
		return
	}
	file, err := os.ReadFile(reqReadByFile.PathConfig)
	if err != nil {
		detail.Errors = append(detail.Errors, err)
	}
	var data Data
	err = json.Unmarshal(file, &data)
	if err != nil {
		detail.Errors = append(detail.Errors, err)
	}
	detail.Metadata.Change = append(detail.Metadata.Change, "file", string(file))
	detail.Data = data
`
	template = strings.ReplaceAll(template, "__READ_BY_FILE__", read_by_file)

	read_by_env := "detail.Metadata.Change = append(detail.Metadata.Change, \"env\")\n"

	for k, _ := range dynamic {
		read_by_env += `	` + underscoreToLowerCamel(k) + ` := os.Getenv("` + k + `")
	if ` + underscoreToLowerCamel(k) + ` != "" {
		detail.Metadata.Change = append(detail.Metadata.Change, "` + k + `", detail.Data.` + underscoreToCamel(k) + `, ` + underscoreToLowerCamel(k) + `)
		detail.Data.` + underscoreToCamel(k) + ` = ` + underscoreToLowerCamel(k) + `
	}
`
	}
	template = strings.ReplaceAll(template, "__READ_BY_ENV__", read_by_env)

	read_by_grpc := `	credentials := insecure.NewCredentials()
	conn, _ := grpc.Dial("127.0.0.1:2023", grpc.WithTransportCredentials(credentials), grpc.WithBlock())
	defer conn.Close()

	client := pb.NewConfigClient(conn)
	resp, err := client.ReadConfig(context.Background(), &pb.RequestReadConfig{
		FromId:       "` + nameProto + `__from_id",
		FromNickname: "` + nameProto + `__from_nickname",
		Content:      "` + nameProto + `__content",
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
	detail.Metadata.Change = append(detail.Metadata.Change, "grpc", resp.Content)`

	template = strings.ReplaceAll(template, "__READ_BY_GRPC__", read_by_grpc)

	read_by_default := ""
	for k, v := range dynamic {
		read_by_default += "detail.Data." + underscoreToCamel(k) + " = \"" + v + "\"\n"
	}

	template = strings.ReplaceAll(template, "__READ_BY_DEFAULT__", read_by_default)

	read_by_ui := pathConfig + "/" + nameProto + "_config.json"
	template = strings.ReplaceAll(template, "__READ_BY_UI__", read_by_ui)

	_, err = fileIndex.Write([]byte(template))
	if err != nil {
		log.Fatalln(err)
	}
	fileIndex.Close()
	_, err = fileMainQuickStart.Write([]byte(`package main

import (
	"` + pathProject + `/` + pathModules + `/` + pathConfigAlias + `/` + pathFunc + `"
	"log"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup
	wg.Add(1)

    // Test modify DefaultConfig
	// defaultConfig := index.DefaultReqReadConfig
	// defaultConfig.PathConfig = "E:\\personal\\golang\\config\\tmp\\40339_config.json"
	// err, conf := index.ReadConfig(defaultConfig)


	err, conf := index.ReadConfig(nil)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {

		for i := 9; i >= 0; i-- {
			time.Sleep(time.Second * 1)
			log.Printf("program will be termainated in %ds later", i)
		}

		wg.Done()
	}()
	log.Println(conf)
	wg.Wait()
}

`))
	if err != nil {
		log.Fatalln(err)
	}
	fileMainQuickStart.Close()

	fileConfig, err := os.Create(pathConfig + "/" + nameProto + "_config.json")
	if err != nil {
		log.Fatalln(err)
	}
	//defer fileConfig.Close()

	fileConfig.WriteString(string(open))
	fileConfig.Close()

	_, err = copy_dir.CopyDir("webserver/sample", pathRootWebserver)
	if err != nil {
		return
	}
	err = os.Rename(pathRootWebserver+"/sample", pathRootWebserver+"/"+pathWebserver)
	if err != nil {
		log.Fatalln(err)
	}
	if pathOutput != "" {
		//_, err = copy_dir.CopyDir(pathModules, pathOutput)
		//if err != nil {
		//	log.Fatalln(err)
		//} else {
		//	time.Sleep(time.Second * 5)
		//	err := os.RemoveAll(pathModules)
		//	if err != nil {
		//		log.Fatalln(err)
		//	}
		//}

		err = os.MkdirAll(pathOutput, 0600)
		if err != nil {
			log.Fatalln(err)
		}
		splitPaths := strings.Split(strings.ReplaceAll(pathModules, "\\", "/"), "/")
		log.Println(pathModules)
		log.Println(pathOutput + "/" + splitPaths[len(splitPaths)-1])
		err = os.Rename(pathModules, pathOutput+"/"+splitPaths[len(splitPaths)-1])
		if err != nil {
			log.Fatalln(err)
		}
		pathModules = pathOutput + "/" + splitPaths[len(splitPaths)-1]
	}
	log.Println("NOT FINISHED, please execute command below:")
	log.Println("protoc -I=\".\" --go_out=./" + pathModules + "/" + pathConfigAlias + " ./" + pathModules + "/" + pathConfigAlias + "/proto/" + nameProto + "/*.proto")
	log.Println("protoc -I=\".\" --go-grpc_out=./" + pathModules + "/" + pathConfigAlias + " ./" + pathModules + "/" + pathConfigAlias + "/proto/" + nameProto + "/*proto")
	log.Println()

}

func underscoreToCamel(s string) string {
	// split the string by underscore
	parts := strings.Split(s, "_")

	// capitalize the first letter of each part
	for i, part := range parts {
		parts[i] = cases.Title(language.English).String(part)

	}

	// join the parts and return
	return strings.Join(parts, "")
}

func underscoreToLowerCamel(s string) string {
	// split the string by underscore
	parts := strings.Split(s, "_")

	// capitalize the first letter of each part except the first one
	for i := 1; i < len(parts); i++ {
		parts[i] = strings.Title(parts[i])
	}

	// join the parts and return
	return strings.Join(parts, "")
}
