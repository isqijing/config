#!/bin/bash
#Program:
#   clean tmp generate
#History:
#2023/08/27 junfenghe.cloud@qq.com  version:0.0.1 msg:init

path=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export path

bin=$(dirname ${BASHS_SOURCE[0]-${0}})
bin=$(cd ${bin}; pwd)
echo ${bin}

rm -rf ${bin}/output/servers/*

echo "rm -rf ${bin}/output/servers/*[finished.]"
## `config` is important
path_proto_output_proto=${bin}/proto/output/proto
ls ${path_proto_output_proto} | grep -v "config" | xargs -I {} rm -rf  ${path_proto_output_proto}/{}

echo "ls ${path_proto_output_proto} | grep -v "config" | xargs -I {} rm -rf  ${path_proto_output_proto}/{} [finished.]"

path_proto=${bin}/proto
ls ${path_proto}  | grep -vE "^config$|^output$" | xargs -I {} rm -rf ${path_proto}/{}

echo "ls ${path_proto}  | grep -vE "^config$\|^output$" | xargs -I {} rm -rf ${path_proto}/{} [finished.]"



