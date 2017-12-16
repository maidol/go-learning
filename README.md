# go-learning

## [Golang标准库](https://github.com/polaris1119/The-Golang-Standard-Library-by-Example)
## [Golang中文网](https://studygolang.com/)
## [golang.org](https://golang.org/pkg/)
## [golang在线api文档](https://gowalker.org/)
## [godoc.org](https://godoc.org/)

## 查看文档

### 使用命令
```bash
go doc time # 包名
go doc time.Since # 包成员(函数, 变量, 类型 ...)
go doc time.Duration.Seconds # 包成员的方法
```
### 线上godoc.org/gowalker.org

### 使用本地godoc, localhost:8000
```bash
godoc -http :8000
```

## 查看包元数据 go list
```bash
go list -json time # 包名
```

## 安装依赖 go get 为GOPATH管理代码

## 构建安装 go build / go install

## 直接代码启动 go run main.go

## 依赖管理工具 dep 为应用管理代码
[教程](https://studygolang.com/articles/10589)

### 设置环境变量

```bash
# 设置环境变量 使用vendor目录
GO15VENDOREXPERIMENT=1
```

### 安装

```bash
go get -u github.com/golang/dep/cmd/dep
```

### 初始化

```bash
# Gopkg.lock, Gopkg.toml, vender
dep init
```

### 其他依赖管理工具 [glide](https://my.oschina.net/u/553243/blog/1475626) [godep](https://studygolang.com/articles/4385)

## 测试 go test

## travis