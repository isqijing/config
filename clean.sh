#!/bin/bash
#Program:
#   clean tmp generate
#History:
#2023/08/27 junfenghe.cloud@qq.com  version:0.0.1 msg:init;
#2023/11/19 junfenghe.cloud@qq.com	version: 0.0.2	msg:add comment;


path=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export path

bin=$(dirname ${BASHS_SOURCE[0]-${0}})
bin=$(cd ${bin}; pwd)
echo ${bin}


## 1. delete all files in `output/servers` directory based on current directory.
rm -rf ${bin}/output/servers/*

echo "rm -rf ${bin}/output/servers/*[finished.]"


## `config` is important
path_proto_output_proto=${bin}/proto/output/proto
ls ${path_proto_output_proto} | grep -v "config" | xargs -I {} rm -rf  ${path_proto_output_proto}/{}

echo "ls ${path_proto_output_proto} | grep -v "config" | xargs -I {} rm -rf  ${path_proto_output_proto}/{} [finished.]"

path_proto=${bin}/proto
ls ${path_proto}  | grep -vE "^config$|^output$" | xargs -I {} rm -rf ${path_proto}/{}

echo "ls ${path_proto}  | grep -vE "^config\$\|^output\$" | xargs -I {} rm -rf ${path_proto}/{} [finished.]"

path_config=${bin}/config

ls ${path_config} | grep -vE "^config_config.json$" | xargs -I {} rm -rf ${path_config}/{}

echo "ls ${path_config} | grep -vE "^config_config.json\$" | xargs -I {} rm -rf ${path_config}/{} [finished.]"

path_webserver=${bin}/webserver

ls ${path_webserver} | grep -vE "^sample$" | xargs -I {} rm -rf ${path_webserver}/{}

echo "ls ${path_webserver} | grep -vE "^sample$" | xargs -I {} rm -rf ${path_webserver}/{} [finished.]"

path_output_webserver=${bin}/output/webserver

rm -rf ${path_output_webserver}/*

echo "rm -rf ${path_output_webserver}/* [finished.]"
