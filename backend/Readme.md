### 目录说明

- cmd: gorm 生成 mysql 表模型和操作接口
- biz
  - entity: 一些元操作
  - hadoop: hadoop交互，hive相关操作
  - handler: 接口统一入口
  - microtype: 错误处理
  - mq: 消息队列相关
  - mw: 中间件
  - router: 路由
  - usecase: 接口处理
  - util: 一些公用处理函数
- consumer: mq消息组
- gohive: go的hive驱动，由[beltran/gohive: Go driver for Apache Hive (github.com)](https://github.com/beltran/gohive) 改造而来
- idl: idl文件
- impl: 前期数据处理代码
- optimize_prefixspan：简化版prefixspan算法
- prefixspan：prefixspan算法实现



### 相关命令

#### 更新IDL

- hz update -idl idl/backend.thrift



#### MQ
1. 进入源码目录 /bin
2. 执行下述两条命令
   1. windows:
      1. start mqnamesrv.cmd
      2. start mqbroker.cmd -n 127.0.0.1:9876 autoCreateTopicEnable=true
   2. linux(虚拟机):
      1. nohup ./mqnamesrv &
      2. nohup ./mqbroker -n 192.168.81.131:9876 &

#### Hadoop 
1. 启动Hadoop: start-dfs.sh start-yarn.sh
2. 启动HiveServer2: hive --service hiveserver2