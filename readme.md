# 通用导出服务

golang实现的高性能、企业级的通用数据导出服务

- 异步任务队列
-  资源复用协程池
-  golang有序并发爬虫
- 分页写入excel

## 目录

- [背景介绍](##背景介绍)
- [快速安装](##快速安装)
  - [配置文件](###配置文件)
  - [数据库结构](###数据库结构)
- [快速接入](##快速接入)
- [设计与架构](##设计与架构)
- [接口文档](##接口文档)



## 背景介绍



## 快速安装

要求:

- go  >= 1.14
- redis >=4.0
- mysql >= 5.6

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

> **!!!一个小细节!!!** json 并不支持注释，所以请在启动前检擦配置文件中注释是否删干净

```json
{
  "base": { // 基础配置
    "debug": true,
    "app_name": "export-server",
    "describe" : "通用导出服务",
    "env": "dev",
    "http_port": "8080"
  },
  "redis": { // redis 配置
    "open": true,
    "addr": "redis-serv:6379",
    "db": 0,
    "pwd": "",
    "pool_size": 30, 
    "max_reties": 3,
    "idle_timeout": 1000,
    "key_prefix": "export-server"
  },
  "mysql":{ // mysql 配置
    "dsn": "root:root@tcp(mysql-serv:3306)/export?charset=utf8&parseTime=True&loc=Local",
    "debug": true,
    "maxIdleConns": 10,
    "maxOpenConns": 100,
    "connMaxLifetime" : 6
  },
  "log": {  // 日志文件配置
    "type": "file",
    "path": "./logs/log.log",
    "level": "DEBUG",
    "max_size": 32,
    "max_age": 30,
    "max_backups": 300
  },
  "http_req_conn": 10, // 使用http数据源时 默认并发数
  "taskPool": {        // 导出服务 worker 数量
    "httpWorker" : 1,  // http数据源 worker 数量
    "rowWorker" : 1    // 直接数据源 worker 数量
  },
  "tmp_storage": { // 临时目录
    "outexcel_tmp": "/tmp/outexcel",  // 生成excel存放的临时目录 任务进行中的数据 任务结束后会自动删除
    "source_raw" : "/tmp/params" // 直接数据源 的参数存储临时目录
  },
  "exp_storage": { // 导出结果存储介质
    "channel": "local",  // 选择存储介质:支持本地存储、阿里云oss
    "local" : {          // 本地存储配置
      "path": "/tmp/down",
      "down_url": "http://127.0.0.1:8080" // 下载域名
    },
    "alioss": {  // 阿里云oss 配置
      "bucket": "*",
      "endpoint_up": "com",
      "key": "",
      "secret": "",
      "down_url": "https://heykui.oss-cn-beijing.aliyuncs.com"
    },
    "ftp": { // TODO: FTP 存储配置

    }
  },
  "excel_maxlines": 5000  // excel 文件最大行数（数据量超过会生成多个文件）
}
```

### 数据库结构

```mysql
DROP TABLE IF EXISTS `export_log`;
CREATE TABLE `export_log`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `hash_key` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '参数哈希',
  `title` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '导出标题',
  `ext_type` char(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '导出类型(文件后缀)',
  `source_type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '数据源类型',
  `param` varchar(800) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '请求参数（json）',
  `user_id` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '用户id',
  `callback` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '回调地址',
  `status` tinyint(4) UNSIGNED NULL DEFAULT NULL COMMENT '状态：1处理中 2导出成功 3导出失败 4导出取消',
  `fail_reason` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '失败理由',
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

DROP TABLE IF EXISTS `export_file`;
CREATE TABLE `export_file`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `hash_key` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `type` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;
```



## 快速接入

数据量较小(<5000行) 可直接将数据携带至请求体中使用[post数据源导出](###post数据源导出)

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

如果数据量较大则可 提供一个[分页查询接口](###http数据源示例) 由导出服务来爬取数据 使用[http数据源导出](###http数据源导出)

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





## 接口文档



### http数据源示例



### post数据源导出



### http数据源导出



### 导出任务详情



### 导出任务历史