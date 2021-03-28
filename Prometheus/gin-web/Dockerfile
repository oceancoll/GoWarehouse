FROM golang
COPY . /$GOPATH/src/promeServer/
WORKDIR /$GOPATH/src/promeServer/
#设置环境变量，开启go module和设置下载代理
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
#会在当前目录生成一个go.mod文件用于包管理
# RUN go mod init
#增加缺失的包，移除没用的包
RUN go mod tidy
RUN go build server.go
EXPOSE 4433:4433
CMD ["go","run","server.go"]