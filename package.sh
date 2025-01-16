#!/bin/bash

# 编译linux版本
go build -o tcpmk_node -ldflags -w main.go

# 排除打包文件
exclude="--exclude=*.gz --exclude=data/config/*.toml --exclude=.git --exclude=.gitignore --exclude=.github --exclude=api --exclude=./config --exclude=middleware --exclude=model --exclude=router --exclude=utils --exclude=go* --exclude=main.go --exclude=*.sh"


tar -zcvf tcpmk_node.tar.gz ${exclude} .

curl --location --request POST 'https://soft.xiaoz.org/api/upload' \
--header 'X-Token: web-f4f9fff2c257ae82c5e1d55040714929' \
--form 'path=/UniBin/tcpmk_node' \
--form "file=@/data/apps/tcpmk_node/tcpmk_node.tar.gz" \
--form 'type=public'
