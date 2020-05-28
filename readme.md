[toc]

### 微服务相关



### 客户端 client：

##### 服务限流

##### 服务熔断

##### 负载均衡

##### 服务发现

##### 连接池



### 服务端 server：

##### 超时管理

##### 请求路由

##### 连接管理

##### 服务注册



### 公共部分：

#### 服务治理

##### 分布式追踪

##### metrics采样

##### Grafana

#### 数据传输

##### 序列化

##### 反序列化

##### 粘包问题

​		

#### 基础组件

##### 日志组件

##### 配置组件




### 服务注册



### 日志管理

#### 熔断



### 分布式设计




##### go-redis.zip       -> brick_sdk/src/github.com/go-redis

##### redisclient.zip -> brick_sdk/src/commclient/redisclient





### 服务目录规范

- controller：存放服务的方法实现
- idl：存放本服务的idl定义
- main：服务的入口函数
- script：存放服务的脚本
- conf：存放服务的配置文件
- app/router：存放服务的路由
- app/config：存放服务的配置代码
- model：存放服务的实体代码
- generate：grpc生成的代码