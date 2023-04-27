### 更新IDL
- hz update -idl idl/backend.thrift



### MQ（window系统启动）
1. 进入源码目录 /bin
2. 执行下述两条命令
   1. start mqnamesrv.cmd
   2. start mqbroker.cmd -n 127.0.0.1:9876 autoCreateTopicEnable=true
