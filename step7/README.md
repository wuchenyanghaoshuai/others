```
nacos+elasticsearch+skywalking+skywalking-agent

172.16.92.2(nacos+elastic+skywalking)
172.16.92.3(nacos+elastic+skywalking)
172.16.92.4(nacos+elastic+skywalking)
172.16.92.5(es的存储节点,只需要配置node.data=true即可)
172.16.92.6(es的存储节点,只需要配置node.data=true即可)

1.elastic安装(集群,自行安排机器)
https://thoughts.teambition.com/share/6020ffe12917070042b9d3b9#title=centos7es集群安装
2.nacos集群安装(集群,自行安排机器)
https://blog.csdn.net/qq_32352777/article/details/86560333
3.skywalking安装(集群自行安排机器)
3.1因为我们安装的es的版本是7.7.1,所以对应的skywalking也是需要下载兼容es7的
wget https://archive.apache.org/dist/skywalking/8.3.0/apache-skywalking-apm-es7-8.3.0.tar.gz
tar -xf apache-skywalking-apm-es7-8.3.0.tar.gz -C /data/service/
cd apache-skywalking-apm-bin-es7/


3.2设置UI高可用(主要是listOfServers)
vim webapp.yml
server:
  port: 8080

collector:
  path: /graphql
  ribbon:
    ReadTimeout: 10000
    # Point to all backend's restHost:restPort, split by ,
    listOfServers: 192.168.200.96:12800,192.168.200.97:12800,192.168.200.98:12800
3.3修改配置文件(只需要修改我截图的部分即可)

cluster.standalon注销
cluster.selector: xxxx  改为nacos
cluster.nacos.serviceName: 起一个名字 或者用默认的(在nacos中可以看的到)
cluster.nacos.hostPort 写上面安装的nacos集群
cluster.nacos.user/password 可以不写
新增cluster.nacos.internalComHost和cluster.nacos.internalComPort
host是本机器的ip
port是11180端口，记得集群内的每个机器的ip不一样，scp的时候记得修改
```

