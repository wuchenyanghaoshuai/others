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
![image](https://user-images.githubusercontent.com/39818267/139527649-a6a26b2a-5af4-4f8e-ad4a-b990dc39d330.png)
![image](https://user-images.githubusercontent.com/39818267/139527655-36d3be5d-daf6-43d0-9301-0580b39ccd7f.png)

```
storage.selector: elasticsearch7
storage.elasticsearch7.namespace: xxx-logs
storage.elasticsearch7.clusterNodes: 192.168.xxx:9200,192.168.xxx:9200,192.168.xxx:9200
```
![image](https://user-images.githubusercontent.com/39818267/139527674-1d3dc465-bab4-4a1e-b679-552a892be15d.png)

```
3.4 skywalking-agent 
只需要打开配置文件修改一下即可
collector.backend_service
找到对应的地址

制作一个包含skywalking-agent的镜像即可或者直接使用我的
docker pull wuchenyanghaoshuai/jre-sky-agent:v1
在/data/agent下有配置文件，这个不需要动，只需要在k8s的yaml里指定evn即可
然后访问skywalking的ui即可访问到
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: kkb-corgi-tenant
spec:
  replicas: 1
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge:
      maxUnavailable: 0
  minReadySeconds: 5
  revisionHistoryLimit: 10
  template:
    metadata:
      labels:
        app: kkb-corgi-tenant
    spec:
      containers:
      - name: kkb-corgi-tenant
        image: registry.cn-beijing.aliyuncs.com/kaikeba/kkb-corgi-tenant-biz:{{build_image_tag}}
        ports:
        - containerPort: 8080
        readinessProbe:
          httpGet:
            path: /manage/health
            port: 8081
          initialDelaySeconds: 20
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 2
        resources:
          requests:
            memory: "256Mi"
          limits:
            memory: "2300Mi"
        env:
        - name: APP_ENV
          value: pre
        - name: JAVA_OPTS
          value: "-Xms2g -Xmx2g -javaagent:/data/agent/skywalking-agent.jar -javaagent:/data/agent/jmx_prometheus_javaagent-0.14.0.jar=9180:/data/agent/jmx_exporter.yml"
        - name: APP_ID
          value: kkb-corgi-tenant
        - name: SW_AGENT_NAME
          value: kkb-corgi-tenant
        - name: SW_AGENT_COLLECTOR_BACKEND_SERVICES
          value: xxxxxxxxx:11800
        - name: SERVICE_DISCOVERY_ADDR
          value: xxxxxx:8848            
        - name: SERVICE_DISCOVERY_NAMESPACE
          value: xxxx
```
