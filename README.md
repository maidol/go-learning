# go-learning
----

## [godoc.org](https://godoc.org/)

## [golang在线api文档](https://gowalker.org/)

## [Golang标准库](https://github.com/polaris1119/The-Golang-Standard-Library-by-Example)

## [Golang中文网](https://studygolang.com/)

## [golang.org](https://golang.org/pkg/)

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

## 开发
----

- ### 查看包元数据 go list

```bash
go list -json time # 包名
```

- ### 安装依赖

```bash
go get # 为GOPATH管理代码
```

- ### 构建/安装

```bash
go build
go install
```

- ### 直接代码启动

```bash
go run main.go
```

- ### 依赖管理工具, dep 为应用管理代码, [dep教程](https://studygolang.com/articles/10589)


>- 设置环境变量

```bash
# 设置环境变量 使用vendor目录
GO15VENDOREXPERIMENT=1
```

>- 安装

```bash
go get -u github.com/golang/dep/cmd/dep
```

>- 初始化

```bash
# Gopkg.lock, Gopkg.toml, vender
dep init
```

>- 其他依赖管理工具 [glide](https://my.oschina.net/u/553243/blog/1475626) [godep](https://studygolang.com/articles/4385)

### 测试 go test

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

### 基准测试(衡量/优化性能)

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

>- 剖析net/http
```bash
$ go test -run=NONE -bench=ClientServerParallelTLS64 \
    -cpuprofile=cpu.log net/http
$ go tool pprof -text -nodecount=10 ./http.test cpu.log
```

### 示例函数
根据示例函数的后缀名部分，godoc这个web文档服务器会将示例函数关联到某个具体函数或包本身，因此ExampleIsPalindrome示例函数将是IsPalindrome函数文档的一部分，Example示例函数将是包文档的一部分

## 持续集成, travis, ci
----

## 记录
---------
- ### [fmt字符串格式化占位符](https://studygolang.com/articles/2644)