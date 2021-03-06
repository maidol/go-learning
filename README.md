# go-learning
----

## [godoc.org](https://godoc.org/)

## [golang在线api文档](https://gowalker.org/)

## [Golang标准库](https://github.com/polaris1119/The-Golang-Standard-Library-by-Example)

## [Golang中文网](https://studygolang.com/)

## [golang.org](https://golang.org/pkg/)

## [GO语言圣经(中文版) ](https://docs.hacknode.org/) [web_doc](https://wizardforcel.gitbooks.io/gopl-zh/content/) [gitbook/pdf_download](https://legacy.gitbook.com/book/wizardforcel/gopl-zh/details)

## [GO命令详解](http://wiki.jikexueyuan.com/project/go-command-tutorial/0.14.html)

## [GO入门指南](http://wiki.jikexueyuan.com/project/the-way-to-go/directory.html)

## 查看文档
----

- ### 使用命令

```bash
go doc time # 包名
go doc time.Since # 包成员(函数, 变量, 类型 ...)
go doc time.Duration.Seconds # 包成员的方法
```
- ### 线上godoc.org/gowalker.org

- ### 使用本地godoc, localhost:8000

```bash
godoc -http :8000
```

- ### 查看包元数据 go list

```bash
go list -json time # 包名
```

## 开发
----

- ### 项目初始化

依赖管理工具, dep 为应用管理代码, [dep教程](http://www.mamicode.com/info-detail-1947553.html)[github: dep](https://github.com/golang/dep)

>- 设置环境变量

```bash
# 设置环境变量 使用vendor目录
GO15VENDOREXPERIMENT=1
```

>- 安装

```bash
go get -u github.com/golang/dep/cmd/dep
```

>- 初始化项目

```bash
mkdir sample
cd sample
# Gopkg.lock, Gopkg.toml, vender
dep init
dep init -v #
```

>- 安装依赖

```bash
dep ensure # 
dep ensure -update
dep ensure -add github.com/pkg/errors
```

>- 其他依赖管理工具 [glide](https://github.com/Masterminds/glide) [godep](github.com/tools/godep) [govendor](https://github.com/kardianos/govendor)

- ### 监控代码变化自动重启 [gin](https://github.com/codegangsta/gin) [fresh](https://github.com/pilu/fresh)

>- 安装

```bash
go get github.com/codegangsta/gin
```

>- 运行

```bash
gin -p 3000 -a 8080 # -p 代理端口, -a web程序端口, 访问 localhost:3000 => localhost:8080
```

- ### 使用go get安装依赖

```bash
go get # 为GOPATH管理代码
```

- ### 构建/安装

```bash
go build
go install
```

- ### 代码启动

```bash
go run main.go
```

- ### 测试 go test

```txt
示例代码目录
go-learning/word1
go-learning/echo
```

>- 运行测试

```bash
go test
go test -run=${regx} # ${regx} 指定测试符合正则的测试函数
```

>- 运行测试并统计覆盖率

```bash
go test -cover
```

>- 运行测试并生成覆盖率报告

```bash
go test -coverprofile=c.out # 只统计代码是否被运行过
go test -coverprofile=c.out -covermode=count # -covermode=count 统计代码的运行权重
```

>- 查看覆盖率报告

```bash
go tool cover -html=c.out
```

- ### 基准测试(衡量/优化性能)

不应该过度纠结于细节的优化，应该说约97%的场景：过早的优化是万恶之源。仅当关键代码已经被确认的前提下才会进行优化

```txt
示例代码目录
go-learning/word1
```

>- 运行

```bash
go test -bench=. 
go test -bench=. -benchmem # -benchmem 内存分配情况, 频繁的内存分配会影响性能
```

>- 性能剖析

```bash
$ go test -cpuprofile=cpu.out
$ go test -blockprofile=block.out
$ go test -memprofile=mem.out
```

>- 剖析net/http, [参考](https://www.cnblogs.com/ghj1976/p/5473693.html)

```bash
$ go test -run=NONE -bench=ClientServerParallelTLS64 \
    -cpuprofile=cpu.log net/http
$ go tool pprof -text -nodecount=10 ./http.test cpu.log
```

- ### 示例函数

根据示例函数的后缀名部分，godoc这个web文档服务器会将示例函数关联到某个具体函数或包本身，因此ExampleIsPalindrome示例函数将是IsPalindrome函数文档的一部分，Example示例函数将是包文档的一部分

## 持续集成, travis, ci
----

## 记录
----
- ### go get代理问题 [参考0](https://studygolang.com/articles/9490?fr=sidebar) [参考1](https://stackoverflow.com/questions/10383299/how-do-i-configure-go-to-use-a-proxy) [参考2](https://stackoverflow.com/questions/128035/how-do-i-pull-from-a-git-repository-through-an-http-proxy/3406766#3406766)
```txt
go get 需要http代理, shadowsocks使用的是socks5代理, 需要添加http代理在shadowsocks的前端, 而shadowsocks作为二级代理。go get 下载package时, 第一步go get先根据包名获取真正的代码下载地址, 再使用版本控制软件下载代码, 最后go安装。这里的http代理涉及到两个, 一个是go get使用的代理, 另一个是版本控制软件使用的代理, 需分别设置
```
>- 安装并设置http代理[cow](https://github.com/cyfdecyf/cow/)
>- 设置版本控制软件的http代理[git mercurial svn](https://github.com/golang/go/wiki/GoGetProxyConfig) [相关问题](https://stackoverflow.com/questions/128035/how-do-i-pull-from-a-git-repository-through-an-http-proxy/3406766#3406766)
>- 设置go get的http代理 [相关问题](https://stackoverflow.com/questions/10383299/how-do-i-configure-go-to-use-a-proxy)
>- ps: go 1.12版本之后支持go modules和GOPROXY特性，只要开启go modules就可以使用GOPROXY特性，设置GOPROXY=https://goproxy.io 或 GOPROXY=https://athens.azurefd.net 环境变量也可解决代理问题

```bash
# for linux
#You can set these environment variables in your bash_profile, but if you want to limit their usage to go, you can run it like this:
$ http_proxy=127.0.0.1:7777 go get -u -v code.google.com/p/go.crypto/bcrypt
# for linux
#If that's what you always want, set this alias to avoid typing proxy part every time:
$ alias go='http_proxy=127.0.0.1:7777 go'
```

>- [设置http_proxy/https_proxy环境变量](http://nanxiao.me/en/set-proxy-when-executing-go-get-command/) go get使用的代理

```shell
# powershell
$env:http_proxy="http://127.0.0.1:7777"
$env:https_proxy="http://127.0.0.1:7777"
# cmd
set http_proxy="http://127.0.0.1:7777"
set https_proxy="http://127.0.0.1:7777"
# bash
export http_proxy="http://127.0.0.1:7777"
export https_proxy="http://127.0.0.1:7777"
```

- ### vscode插件
Beautify/.jshintrc Genertor/Docker/Document This/EditorConfig/ESLint/ftp-sync/Git History/Git Merger/Git Project Manager/Go/HTML Snippets/jshint/mdeiawiki/nginx.conf/TODO Highlight/TODO Parser/vscode-icons/markdownlint

- ### [fmt字符串格式化占位符](https://studygolang.com/articles/2644)

- ### ab压测工具

- ### [构建最小的Go程序镜像](http://time-track.cn/build-minimal-go-image.html) [解决错误no such file](https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/)
- ### [go-swagger](https://github.com/go-swagger/go-swagger)

- ### [排序](https://www.cnblogs.com/wangchaowei/p/7802811.html)

- ### [golang交叉编译](https://studygolang.com/articles/6002)

使用Go来构建微服务的一个优点是，它会被编译成二进制包，这样的话，它就不需要框架或者运行依赖，这样非常有利，因为正如前面所说Alpine是一个非常轻量级的分发版，并不是所有C语言依赖库都有安装，所以Go的动态库依赖很可能也没有。所幸的是有专门的方法去禁用了cgo依赖，可以把应用通过链接的方式编译，我们只需要这样告诉编译器去重新构建我们的所有应用包就可以了：

```bash
$ CGO_ENABLED=0 go build -a -installsuffix cgo .
```

我们更详细说一下上边这个命令的细节： 
CGO_ENABLED=0      是一个编译标志，会让构建系统忽略cgo并且静态链接所有依赖； 
-a                 会强制重新编译，即使所有包都是由最新代码编译的； 
-installsuffix cgo 会为新编译的包目录添加一个后缀，这样可以把编译的输出与默认的路径分离。

- ### [Go语言TCP Socket编程](https://studygolang.com/articles/5372)
- ### [文件操作](https://gocn.io/article/40)
- ### [深入解析go](https://tiancaiamao.gitbooks.io/go-internals/content/zh/08.1.html)
