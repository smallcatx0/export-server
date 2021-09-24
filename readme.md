# 通用导出服务

golang实现的高性能、企业级的通用数据导出服务

- 异步任务队列
-  资源复用协程池
-  golang有序并发爬虫
- 分页写入excel

## 目录

- [背景介绍](#背景介绍)
- [快速安装](#快速安装)
  - [配置文件](#配置文件)
  - [数据库结构](#数据库结构)
- [快速接入](#快速接入)
- [设计与架构](#设计与架构)
- [接口文档](#接口文档)



## 背景介绍

各种管理后台均会遇到将查询结果导出成excel的需求（并上传云存储）的需求。在微服务下。企业后台管理系统几十上百个项目。每个项目都去造一遍轮子，对接一遍云存储。即是浪费时间，也容易导致云存储的key分布与各个项目中，难以维护难以规范。

因此，企业导出服务应运而生，各业务仅需要按照[分页查询接口](#http数据源示例)规范，提供一个数据查询接口，只需关注数据，即可非常简单的对接到导出服务，以便捷、高效的方法完成数据导出的需求。

## 快速安装

要求:

- go  >= 1.14
- redis >=4.0
- mysql >= 5.6

编译:

```shell
# 下载源码
git clone git@gitee.com:smallcatx0/export-server.git
# 编译
cd export-server
go build -o ./export-serv .
# 运行
./export-serv -help
#-----output--------
#  -config string
#        配置文件地址 (default "conf/app.json")
#  -help
#        帮助
```

### 配置文件

```yaml
# 基础配置
app_name: export-server #服务唯一名
debug: true
env: dev
http_port: 80
describe: 通用导出服务
# redis 配置
redis:
  open: true
  addr: redis-serv:6379
  db: 0
  pwd: ''
  pool_size: 30
  max_reties: 3
  idle_timeout: 1000
# mysql 配置
mysql:
  dsn: 'root:123123@tcp(mysql-serv:3306)/export?charset=utf8&parseTime=True&loc=Local'
  debug: true
  maxIdleConns: 10
  maxOpenConns: 100
  connMaxLifetime: 6
# 日志配置
log:
  type: file #可选 file stdout
  level: DEBUG 
  path: "./logs/log.log"
  max_size: 32
  max_age: 30
  max_backups: 300

# 临时存储文件设置
storage:
  outexcel_tmp: "/tmp/outexcel" # 生成excel存放的临时目录 任务进行中的数据 任务结束后会自动删除
  source_raw: "/tmp/params"     # 直接数据源 的参数存储临时目录

http_req_conn: 5     # http数据源导出，爬取默认并发数 请求参数中可覆盖该值
excel_maxlines: 5000 # excel文件最大行数（数据量超过会生成多个文件）

# 消费者数量
taskPool:
  httpWorker: 1 # http数据源 worker 数量
  rowWorker: 1  # 直接数据源 worker 数量

# 导出结果存储介质
exp_storage:
  # 选择存储通道
  channel: local # 选择存储介质:支持本地存储、阿里云oss
  # 本地存储
  local:
    path: "/tmp/down"
    down_url: http://127.0.0.1  # 需自行起一个nginx做文件下载服务
  # 阿里云oss
  alioss:
    bucket: heykui
    endpoint_up: oss-cn-beijing.aliyuncs.com
    key: '*****'
    secret: '*******'
    endpoint_down: https://heykui.oss-cn-beijing.aliyuncs.com
  # TODO: ftp
```

### 数据库结构
数据表定义详见`database/export_log.sql` [点击跳转](https://gitee.com/smallcatx0/export-server/blob/master/database/export_log.sql)


## 快速接入

数据量较小(<5000行) 可直接将数据携带至请求体中使用[post数据源导出](#post数据源导出)

接口概要(curl **bash**):

```bash
curl --location --request POST 'http://127.0.0.1:8080/v1/export/raw' \
--header 'Content-Type: application/json' \
--data-raw '{
    "timestamp": 1625106207,
    "ext_type": "XLSX",
    "title": "导出任务的标题",
    "user_id": "导出者用户ID",
    "callback": "http://you-serv/export-notice",
	"source_raw": "[{\"name\":\"姓名\",\"age\":\"年龄\"},{\"name\":\"李雷\",\"age\":17}]"
}'
```

如果数据量较大则可 提供一个[分页查询接口](#http数据源示例) 由导出服务来爬取数据 使用[http数据源导出](#http数据源导出)

接口概要:

```bash
curl --location --request POST 'http://127.0.0.1:8080/v1/export/http' \
--header 'Content-Type: application/json' \
--data-raw '{
    "timestamp": 1625106204,
    "ext_type": "XLSX",
    "title": "导出任务的标题",
    "user_id": "导出者用户ID",
    "callback": "http://you-serv/export-notice",
    "source_http": {
        "url": "http://127.0.0.1:8080/demo/page",
        "method": "GET",
        "header": {},
        "param": {
            "limit": 20,
            "page": 1
        }
    }
}'
```



## 设计与架构



![时序图](https://github.com/smallcatx0/export-server/blob/main/Doc/%E6%97%B6%E5%BA%8F%E5%9B%BE.png?raw=true)



## 接口文档

接口文档详细见[接口文档](https://github.com/smallcatx0/export-server/blob/main/Doc/接口文档.md)

### http数据源示例


http数据源的标准格式，导出服务会并发来爬取该接口的数据并写入excel

- 接口路径: `GET` `/demo/page`

- 请求参数: 

  Query: 
  
  | 参数名 | 必须 | 示例 | 描述                 |
  | ------ | ---- | ---- | -------------------- |
  | page   | 否   | 1    | 分页参数（当前页码） |
  | limit  | 否   | 10   | 分页参数（每页大小） |

- 返回数据:

  Body: 
  
  ```json
  {
      "errcode": 0,
      "msg": "",
      "data": {
          "data": [
              {
                  "Age": "年龄",
                  "ID": "编号",
                  "Name": "姓名"
              },
              {
                  "ID": 0,
                  "Name": "Demo",
                  "Age": 15
              },
              {
                  "ID": 1,
                  "Name": "Demo",
                  "Age": 15
              }
          ],
          "pagetag": {
              "page": 1,
              "limit": 2,
              "total": 500,
              "total_page": 250
          }
      },
      "request_id": "b4ef2d59-9887-4b9d-866b-2f8d787288b3"
  }
  ```
  
  `data.data[0]` 返回的第一行数据必须为EXCEL数据表头
  
  `data.pagetag` 必须为分页参数
  
  

### post数据源导出

在请求中直接携带数据导出

- 接口路径 `POST` `/v1/export/raw`

- 请求：

  Body:

    ```json
    {
        "timestamp": 1625106207,    // 当前时间戳
        "ext_type": "XLSX",         // 导出文件类型
        "title": "请求直接携带数据导出",  // 导出任务的标题
        "user_id": "1231",
        "callback": "",
        "source_raw": "[{\"name\":\"姓名\",\"age\":\"年龄\"},{\"name\":\"李雷\",\"age\":17}]"
    }
    ```

  `source_raw` 为json字符串，第一条为表头

- 返回参数：

    Body:
    
    ```json
    {
        "errcode": 0,  // 错误编码 0表示成功
        "msg": "",
        "data": {
            "hash_key": "5ad5f7b92a816b76e70a21e4154c2485" // 导出任务key
        },
        "request_id": "ba35f44b-d408-45d9-b384-69efd507deb0"
    }
    ```


### http数据源导出


导出某分页查询接口的数据，导出服务会并发去爬取该查询接口并将数据写入excel

- 接口路径: `POST` `/v1/export/http`

- 请求参数: 

  Header:

  | 参数名       | 值               |      |
  | ------------ | ---------------- | ---- |
  | Content-Type | application/json |      |

  Body:

  ```json
  {
      "timestamp": 1625106204,     // 当前时间戳
      "ext_type": "XLSX",          // 导出文件类型
      "title": "http接口数据源导出", // 导出任务的标题
      "user_id": "123",            // 导出者
      "callback": "http://www.baidu.com", // 任务结果回调通知地址
      "source_http": {  // 数据查询接口 
          "url": "http://127.0.0.1:8080/demo/page", // 数据查询地址
          "method": "GET", // 请求方式
          "header": {},    // 请求header
          "param": {       // 分页参数
              "limit": 20, // 每页数据条数 
          }
      }
  }
  ```

  `source_http` 数据查询接口必须遵循[分页数据demo](#分页数据demo)的规范

- 返回数据:

  Body:

    ```json
    {
        "errcode": 0,  // 错误编码 0表示成功
        "msg": "",
        "data": {
            "hash_key": "5ad5f7b92a816b76e70a21e4154c2485" // 导出任务key
        },
        "request_id": "ba35f44b-d408-45d9-b384-69efd507deb0"
    }
    ```

### 导出任务详情


- 接口路径: `GET` `/v1/export/detail`

- 请求参数: 

  Query:

  | 参数名 | 必须 | 示例                             | 描述          |
  | ------ | ---- | -------------------------------- | ------------- |
  | key    | 是   | 2c3a54f7120a0d31f847d43f2c9725b6 | 导出任务的key |

- 返回数据:

  Body:

  ```json
  {
      "errcode": 0,
      "msg": "",
      "data": {
          "created_at": "2021-09-06T17:24:34+08:00", // 任务创建时间
          "down_url": "http://127.0.0.1:8080/v1/down/20219\\6\\2c3a54f7120a0d31f847d43f2c9725b6.zip",  // 下载地址
          "ext_type": "XLSX", 
          "fail_reason": "", // 失败原因
          "hash_key": "2c3a54f7120a0d31f847d43f2c9725b6", // 任务ID
          "id": 4,
          "source_type": "http",  // 数据源
          "status": 2, // 状态：1处理中 2导出成功 3导出失败 4导出取消
          "title": "http接口数据源导出" // 任务标题
      },
      "request_id": "0452d509-de6e-4937-96e5-adbbdc41a82d"
  }
  ```


### 导出任务历史


历史记录查询只返回最近7天的数据，并且只返回了摘要信息，需要下载请再使用导出任务查询接口获取下载链接

- 接口路径: `GET` `/v1/export/history`

- 请求参数: 

  Query:

  | 参数名  | 必须 | 示例 | 描述     |
  | ------- | ---- | ---- | -------- |
  | user_id | 是   | 123  | 导出者ID |

- 返回数据:

  Body:

  ```json
  {
      "errcode": 0,
      "msg": "",
      "data": [
          {
              "created_at": "2021-09-06T17:24:34+08:00",
              "ext_type": "XLSX",
              "fail_reason": "",
              "hash_key": "2c3a54f7120a0d31f847d43f2c9725b6",
              "id": 4,
              "source_type": "http",
              "status": 2,
              "title": "http接口数据源导出"
          },
          {
              "created_at": "2021-09-06T17:01:29+08:00",
              "ext_type": "XLSX",
              "fail_reason": "",
              "hash_key": "68b5e20e3da35f7a4040844a63aaac89",
              "id": 1,
              "source_type": "http",
              "status": 2,
              "title": "http接口数据源导出"
          }
      ],
      "request_id": "c3aadcf8-d5d3-4a2b-888b-691781820a9f"
  }
  ```
