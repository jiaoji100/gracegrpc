## 功能 
```
平滑重启golang gRPC server，重启过程不丢流量。类似facebook的gracehttp
```

## 使用
#### 平滑重启进程
```shell script
kill -SIGUSR2  pid 
```
#### 正常杀进程
```shell script
kill -9 pid 
kill -SIGINT pid
kill -SIGTERM pid 
```

## 测试
#### 启动server
```shell script
cd example/server
go build 
./server-demo
```
#### 启动client
```shell script
cd example/client
go build 
./client-demo
```

#### 平滑重启server
```shell script
kill -SIGUSR2 pid
```