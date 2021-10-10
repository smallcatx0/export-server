# common-export-server

[简体中文](readme.md)|English

High-performance, enterprise-level general data export service implemented by golang

- Asynchronous task queue
- Resource reuse coroutine pool
- Golang ordered concurrent crawler
- Pagination to excel


## Background introduction

Various management backgrounds will meet the needs of exporting query results into excel (and uploading to cloud storage). Under microservices. There are dozens of hundreds of projects in the enterprise back-end management system. Every project has to build wheels and dock with cloud storage. That is a waste of time, it is also easy to cause the key distribution of cloud storage and various projects to be difficult to maintain and standardize.

Therefore, the enterprise export service came into being. Each business only needs to provide a data query interface according to the [分页查询接口](#http数据源示例) specification. You only need to pay attention to the data to connect to the export service very simply. Convenient and efficient method to complete the data export requirements.

## Quick install

Require:

- go  >= 1.14
- redis >=4.0
- mysql >= 5.6

Compile:

```shell
# Download source code
git clone git@github.com:smallcatx0/export-server.git
# Compile
cd export-server
go build -o ./export-serv .
# Run
./export-serv -help
#-----output--------
#  -config string
#       配置文件地址 (default "conf/app.yaml")
#  -help
#       帮助
```

### Configuration

```yaml
# bases
app_name: export-server # Service unique name
debug: true
env: dev
http_port: 80
describe: common-export-server
# redis config
redis:
  open: true
  addr: redis-serv:6379
  db: 0
  pwd: ''
  pool_size: 30
  max_reties: 3
  idle_timeout: 1000
# mysql config
mysql:
  dsn: 'root:123123@tcp(mysql-serv:3306)/export?charset=utf8&parseTime=True&loc=Local'
  debug: true
  maxIdleConns: 10
  maxOpenConns: 100
  connMaxLifetime: 6
# log config
log:
  type: file # Optional file stdout
  level: DEBUG 
  path: "./logs/log.log"
  max_size: 32
  max_age: 30
  max_backups: 300

# Temporary storage file config
storage:
  outexcel_tmp: "/tmp/outexcel" # Generate a temporary directory for excel storage Data in progress of the task will be automatically deleted after the task is over
  source_raw: "/tmp/params"     # json source temporary parameter storage directory

http_req_conn: 5     # http source Crawl the default number of concurrent requests can override this value in the request parameter
excel_maxlines: 5000 # The maximum number of rows of excel file (more than one file will be generated if the amount of data exceeds)

# Number of consumers
taskPool:
  httpWorker: 1 # http source worker number
  rowWorker: 1  # raw source worker number

# Export result storage medium
exp_storage:
  # Select storage channel
  channel: local # Support local storage, Alibaba Cloud OSS
  # local storage
  local:
    path: "/tmp/down"
    down_url: http://127.0.0.1  # Need to start a nginx for file download service
  # Alibaba Cloud OSS
  alioss:
    bucket: heykui
    endpoint_up: oss-cn-beijing.aliyuncs.com
    key: '*****'
    secret: '*******'
    endpoint_down: https://heykui.oss-cn-beijing.aliyuncs.com
  # TODO: ftp
```

### Database structure
The data table definition is detailed in`database/export_log.sql` [Click to jump](https://gitee.com/smallcatx0/export-server/blob/master/database/export_log.sql)


## Quick access

The data volume is small (<5000 rows), and the data can be directly carried into the request body for use[Post data source export](#post数据源导出)

Interface summary(curl **bash**):

```bash
curl --location --request POST 'http://127.0.0.1:8080/v1/export/raw' \
--header 'Content-Type: application/json' \
--data-raw '{
    "timestamp": 1625106207,
    "ext_type": "XLSX",
    "title": "task title",
    "user_id": "export user id",
    "callback": "http://you-serv/export-notice",
	"source_raw": "[{\"name\":\"姓名\",\"age\":\"年龄\"},{\"name\":\"李雷\",\"age\":17}]"
}'
```

If the amount of data is large, a [paged query interface] (#httpdata source example) can be provided by the export service to crawl the data. Use [http data source export] (#http data source export)

Interface summary:

```bash
curl --location --request POST 'http://127.0.0.1:8080/v1/export/http' \
--header 'Content-Type: application/json' \
--data-raw '{
    "timestamp": 1625106204,
    "ext_type": "XLSX",
    "title": "task title",
    "user_id": "export uid",
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


## Design and architecture



![时序图](https://github.com/smallcatx0/export-server/blob/main/Doc/%E6%97%B6%E5%BA%8F%E5%9B%BE.png?raw=true)



## Interface documentation

See the interface documentation for details[interface documentation](https://github.com/smallcatx0/export-server/blob/main/Doc/接口文档.md)
