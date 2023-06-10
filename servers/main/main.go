package main

import (
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

	var pathProjectProto string
	pathProjectProto = "config/proto/output/proto" // path which your new project'proto
	flag.StringVar(&pathProjectProto, "path_project_proto", pathProjectProto, "path which your new project'proto")
	flag.StringVar(&pathProjectProto, "ppp", pathProjectProto, "path which your new project'proto")
	flag.StringVar(&pathProjectProto, "p", pathProjectProto, "path which your new project'proto")
	var name string
	flag.StringVar(&name, "name", nameProto, "${name}.proto")
	flag.StringVar(&name, "n", nameProto, "${name}.proto")

	flag.Parse()

	nameProto = name
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

	pathProto := "proto/" + nameProto
	err = os.Mkdir(pathProto, 0600)
	if err != nil {
		log.Fatalln(err)
	}

	fileProto, err := os.Create(pathProto + "/" + nameProto + ".proto")
	if err != nil {
		log.Fatalln(err)
	}
	defer fileProto.Close()
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

	log.Println("NOT FINISHED, please execute command below:")
	log.Println("protoc -I=\".\" --go_out=./proto/output ./proto/" + nameProto + "/*.proto")
	log.Println("protoc -I=\".\" --go-grpc_out=./proto/output ./proto/" + nameProto + "/*proto")
	log.Println()

	// 生成f_main_config.go
	byteTemplate, err := os.ReadFile("template.txt")
	if err != nil {
		log.Fatalln(err)
	}
	//log.Println(string(byteTemplate))
	pathOutputMain := "output/your_main_" + nameProto
	err = os.Mkdir(pathOutputMain, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	fileModule, err := os.Create(pathOutputMain + "/main_" + nameProto + ".go")
	if err != nil {
		log.Fatalln(err)
	}
	defer fileModule.Close()

	template := string(byteTemplate)
	template = strings.ReplaceAll(template, "__PATH_PROJECT_PROTO__", pathProjectProto)
	template = strings.ReplaceAll(template, "__NAME_PROTO__", nameProto)
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

	read_by_file := `	file, err := os.ReadFile("config/` + nameProto + `_config.json")
	if err != nil {
		detail.Errors = append(detail.Errors, err)
	}
	var data Data
	err = json.Unmarshal(file, &data)
	if err != nil {
		detail.Errors = append(detail.Errors, err)
	}
	detail.Metadata.Change = append(detail.Metadata.Change, "file", string(file))
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
	_, err = fileModule.Write([]byte(template))
	if err != nil {
		log.Fatalln(err)
	}
	fileConfig, err := os.Create("config/" + nameProto + "_config.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer fileConfig.Close()
	fileConfig.WriteString(string(open))

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
