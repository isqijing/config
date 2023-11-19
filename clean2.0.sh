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
rm -rf ${bin}/output2.0/*
