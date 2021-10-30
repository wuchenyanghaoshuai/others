#基础监控方案
```
1.阿里云资源(ECS,MySql,Elasticsearch,Redis,MongoDb)状态监控,内存告警CPU告警,重启等一系列告警
2.Prometheus node资源的监控: CPU整体使用百分比，磁盘使用率，TCP状态监控，内存使用率，负载监控，文件描述符，磁盘的读写速率，进程监控，带宽流量
3.nginx  ng状态,http状态码,活跃进程数,等待进程数
  rabbitmq 队列积压,集群状态
  nacos 集群状态,jvm监控
  kafka消费者组积压
4.告警是prometheus的altermanager来发送+ 夜莺平台 两种方式来做
```

