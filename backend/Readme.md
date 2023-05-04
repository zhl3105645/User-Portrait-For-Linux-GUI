### 更新IDL
- hz update -idl idl/backend.thrift



### MQ（window系统启动）
1. 进入源码目录 /bin
2. 执行下述两条命令
   1. windows:
      1. start mqnamesrv.cmd
      2. start mqbroker.cmd -n 127.0.0.1:9876 autoCreateTopicEnable=true
   2. linux(虚拟机):
      1. nohup ./mqnamesrv &
      2. nohup ./mqbroker -n 192.168.81.131:9876 &

### Hadoop 
1. 启动Hadoop: ./sbin/start-dfs.sh
2. 启动HiveServer2: hive --service hiveserver2