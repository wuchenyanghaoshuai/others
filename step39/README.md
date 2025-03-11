# 本文介绍如何使用Opentelemetry+jaeger实现分布式链路追踪。
## 具体关于opentelemetry的介绍可以看下面的文章
``` shell
 https://blog.csdn.net/qq_62368250/article/details/143516314
```
## 1. 安装CertManager 首选 step01
```shell
# 地址在 https://cert-manager.io/docs/installation/kubectl
# 请注意我部署的时候k8s集群1.23.1 如果是低于我这个版本的话,会对go环境支持的不是这么友好
# 去官方github去下载跟k8s集群版本对应的yaml文件，我这里使用的是v1.13.6
# 可以先wget下来，然后查看yaml里的具体内容然后再apply
# 如果你的空k8s版本低于1.19.0 可以去github里寻找低版本的yaml，一般会有版本对照表
```

## 2. 安装Jaeger all-in-one模式 step02
```yaml 
# 我这里使用的ingress暴露的，你也可以使用nodeport或者ingress暴露
# 关于image,我这里使用的是官方的镜像转存到阿里云cr的

apiVersion: apps/v1
kind: Deployment
metadata:
  name: jaeger
spec:
  replicas: 1
  selector:
    matchLabels:
      app: jaeger
  template:
    metadata:
      labels:
        app: jaeger
    spec:
      containers:
      - name: jaeger
        image: jaegertracing/all-in-one:latest
        # image: registry.cn-hangzhou.aliyuncs.com/wuchenyang/jaegertracing-all-in-one:latest # 阿里云镜像
        env:
        - name: COLLECTOR_OTLP_ENABLED
          value: "true"
        ports:
        - containerPort: 16686
        - containerPort: 14268
---
apiVersion: v1
kind: Service
metadata:
  name: jaeger
spec:
  selector:
    app: jaeger
  type: ClusterIP
  ports:
    - name: ui
      port: 16686
      targetPort: 16686
    - name: collector
      port: 14268
      targetPort: 14268
    - name: http
      protocol: TCP
      port: 4318
      targetPort: 4318
    - name: grpc
      protocol: TCP
      port: 4317
      targetPort: 4317      
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-jager
  namespace: default
spec:
  ingressClassName: nginx
  rules:
    - host: jager.wuchenyangzhendehaoshuai.com
      http:
        paths:
          - backend:
              service:
                name: jaeger
                port:
                  number: 16686
            path: /
            pathType: ImplementationSpecific
```


## 3. 安装opentelemetry-operator step03

```shell
#opentelemetry-operator与k8s版本有关，需要根据k8s版本进行安装，具体看具体看下面的地址
# https://github.com/open-telemetry/opentelemetry-operator/blob/main/docs/compatibility.md
# 如果你的k8s集群高于1.23 直接看这个md里的对照表就行，如果低于1.23，需要看下面的地址(本人1.23.1)里的对照关系
# https://github.com/open-telemetry/opentelemetry-operator/tree/v0.85.0 || k8s 1.22.1
# 然后在tag里找到v0.100.0 的installation manifest for kubernetes 右键复制链接地址 ,然后在终端wget 即可
# 在终端执行
kubectl apply -f opentelemetry-operator-v0.100.0.yaml
# 查看是否安装成功
kubectl get pod -n opentelemetry-operator-system

# 如果下载不了image可以使用我的镜像(我pull下来tag的)
#registry.cn-hangzhou.aliyuncs.com/wuchenyang/kube-rbac-proxy:v0.13.1
#registry.cn-hangzhou.aliyuncs.com/wuchenyang/opentelemetry-operator:0.100.0
#实际使用的时候也可以将镜像pull下来，tag成官方的image
```
## 4. 部署otel-collector 跟Instrumentation step04

```yaml

# cat 04otlp/Instrumentation.yaml
#这个因为我是java的demo所以下面那块写了java，如果是python或者go去官方照文档即可都是一样的，这里刚刚好有一个go的demo跟一个java的demo
# 有了这个配置以后，你可以在你的deploy里添加一个annotation instrumentation.opentelemetry.io/inject-java: "true"
#下面有我的deployoy 配置
apiVersion: opentelemetry.io/v1alpha1
kind: Instrumentation
metadata:
  name: instrumentation-sample
spec:
  propagators:
    - tracecontext
    - baggage
    - b3
  sampler:
    type: parentbased_traceidratio
    argument: "1"
  env:
    - name: OTEL_EXPORTER_OTLP_ENDPOINT
      value: otel-collector.default:4318
  java:    
    env:
      - name: OTEL_EXPORTER_OTLP_ENDPOINT
        value: http://otel-collector.default:4317
  go:    
    env:
      - name: OTEL_EXPORTER_OTLP_ENDPOINT
        value: http://otel-collector.default:4318
```

```yaml
#其他配置的代码
java:  # Java 特定配置
    env:
      - name: OTEL_EXPORTER_OTLP_ENDPOINT
        value: "http://otel-collector.default:4317"  # gRPC 端点（可选覆盖全局）
  python:  # Python 特定配置
    env:
      - name: OTEL_EXPORTER_OTLP_ENDPOINT
        value: "http://otel-collector.default:4317"  # gRPC 端点（可选覆盖全局）
  go:  # Go 特定配置（实验性 eBPF）
    env:
      - name: OTEL_EXPORTER_OTLP_ENDPOINT
        value: "http://otel-collector.default:4318"  # gRPC 端点（可选覆盖全局）
```

```yaml
# cat 04otlp/otlp.yaml
#这个配置就是otel-collector的配置，这里我使用的是jaeger的exporter，你也可以使用其他的exporter
#因为jager的后端没有配置es等存储，我这里使用的本机，所以数据保存不了30天往上
apiVersion: opentelemetry.io/v1beta1
kind: OpenTelemetryCollector
metadata:
  name: otel
spec:
  config:
    receivers:
      otlp:
        protocols:
          grpc: {}
          http: {}
    processors:
      memory_limiter:
        check_interval: 1s
        limit_percentage: 75
        spike_limit_percentage: 15
      batch:
        send_batch_size: 10000
        timeout: 10s
    exporters:
      otlp/jaeger:
        endpoint: "jaeger.default:4317"
        tls:
          insecure: true
    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: [memory_limiter, batch]  # 添加 processors
          exporters: [otlp/jaeger]
```

## 5. 部署Java-Demo step05

```yaml
#可以直接使用我的yaml
#svc那块你可以直接用ingress或者nodeport都可以，我是用了ingress，本地做了域名映射
#可以清晰的看到下面的yaml使用了annotation instrumentation.opentelemetry.io/inject-java: "true"
#如果是go项目的话需要添加instrumentation.opentelemetry.io/inject-go: "true"，还有一个需要注意的地方，go项目需要单独提权限,而且是必须的
# 具体看官方文档的267行-287行,https://github.com/open-telemetry/opentelemetry-operator/blob/v0.85.0/README.md
#```yaml
#securityContext:
#    capabilities:
#     add:
#     - SYS_PTRACE
#    privileged: true
#    runAsUser: 0
#```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: java-demo
spec:
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  replicas: 1
  selector:
    matchLabels:
      app: java-demo
  template:
    metadata:
      labels:
        app: java-demo
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9180"
        prometheus.io/path: /metrics
        instrumentation.opentelemetry.io/inject-java: "true"
    spec:
      containers:
        - name: java-demo
          image: registry.cn-hangzhou.aliyuncs.com/wuchenyang/java-demo:test-b0cbb8
          ports:
            - containerPort: 8080
          readinessProbe:
            httpGet:
              path: /
              port: 8080
            initialDelaySeconds: 20
            periodSeconds: 5
            timeoutSeconds: 3
            failureThreshold: 2
          resources:
            requests:
              memory: "512Mi"
              cpu: "0.5"
            limits:
              memory: "1024Mi"
              cpu: "1"
          env:
          - name: "JAVA_OPTS"
            value: "-Xms1g -Xmx1g -javaagent:/data/agent/skywalking-agent.jar -javaagent:/data/agent/jmx_prometheus_javaagent-0.14.0.jar=9180:/data/agent/jmx_exporter.yml"
          - name : PROFILE
            value: test
---
apiVersion: v1
kind: Service
metadata:
  name: java-demo
spec:
  selector:
    app: java-demo
  ports:
    - protocol: TCP
      port: 8080
      targetPort: 8080
      nodePort: 30000
  type: NodePort
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-java-demo
spec:
  ingressClassName: nginx
  rules:
  - host: test.wuchenyangzhendehaoshuai.com
    http:
      paths:
      - backend:
          service:
            name: java-demo
            port:
              number: 8080
        path: /
        pathType: ImplementationSpecific
```

## 6. 访问java-demo以及jaeger-ui

### 浏览器输入 http://test.wuchenyangzhendehaoshuai.com/
![image](https://github.com/user-attachments/assets/069664b8-3444-4f80-bdf0-e152580bb442)
### 点击页面选择数据选择提交
![image](https://github.com/user-attachments/assets/22063593-ffab-4547-a536-ec194f5741c1)
### 浏览器输入 http://jager.wuchenyangzhendehaoshuai.com/
![image](https://github.com/user-attachments/assets/8df26c48-7d98-4a72-846d-54e9b9fd8b7b)
### 查询刚刚提提交动作m的trace
![image](https://github.com/user-attachments/assets/c39c0aa3-9ee4-4ecd-88c3-87dbb19a9852)


## 7. Go项目demo step06

### 注意
```shell
go跟java服务不一样的地方在于java是以initcontainer的方式，go是sidcar的方式，go项目启动以后是2/2的方式,如果你kubectl get pods 去看的话
如果你要使用我的go程序，需要提前在你的环境里准备好mysql环境，并且创建好数据库，表结构可以参考我的sql文件，具体看docker目录中详细解释
由于go的ebpf插桩在operator的0.88版本开始，我的k8s版本为1.23.1,所以我本次使用的是v0.100.0，做测试的时候，一定要看好自己的k8s版本与operator版本的对应关系
按照官方文档，需要在注解里添加主程序的位置
由于我的系统使用的是centos7.9，内核版本默认的是3.10.x，所以需要升级内核版本，要不然的话在go自动插桩的时候那个会直接报错，但是不会影响服务正常
升级步骤如下(所有节点都需要执行)：
  参考文档: https://www.jb51.net/article/265146.htm
  下载地址: http://193.49.22.109/elrepo/kernel/el7/x86_64/RPMS/
  wget http://193.49.22.109/elrepo/kernel/el7/x86_64/RPMS/kernel-ml-devel-5.6.14-1.el7.elrepo.x86_64.rpm
  wget http://193.49.22.109/elrepo/kernel/el7/x86_64/RPMS/kernel-ml-5.6.14-1.el7.elrepo.x86_64.rpm
  rpm -ivh *.rpm
  grub2-set-default 0
  grub2-mkconfig -o /boot/grub2/grub.cfg
  reboot
  [root@km ~]# uname -r
  5.6.14-1.el7.elrepo.x86_64
```
```yaml
# 我是使用ingress暴露，可以根据自身情况选择使用nodeport，
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gin-demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gin-demo
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
        app: gin-demo
      annotations:
        instrumentation.opentelemetry.io/inject-go: "true"
        instrumentation.opentelemetry.io/otel-go-auto-target-exe: "/app/main"
        instrumentation.opentelemetry.io/container-names: "gin-demo"
    spec:
      containers:
      - name: gin-demo
        image: registry.cn-hangzhou.aliyuncs.com/wuchenyang/gin-demo:v4
        securityContext:
          privileged: true
          runAsUser: 0
        imagePullPolicy: Always
        ports:
        - containerPort: 8081
        readinessProbe:
          httpGet:
            path: /ping
            port: 8081
          initialDelaySeconds: 20
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 2
        resources:
          requests:
            memory: "256Mi"
            cpu: 0.5
          limits:
            memory: "512Mi"
            cpu: 1
---
apiVersion: v1
kind: Service
metadata:
  name: gin-demo-service
  namespace: default
spec:
  ports:
  - name: http
    port: 8081
    targetPort: 8081
    protocol: TCP
  type: NodePort
  selector:
    app: gin-demo
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-gin-demo
spec:
  ingressClassName: nginx
  rules:
  - host: gin.wuchenyangzhendehaoshuai.com
    http:
      paths:
      - backend:
          service:
            name: gin-demo-service
            port:
              number: 8081
        path: /
        pathType: ImplementationSpecific
```

## 8. 访问Go项目以及jaeger-ui

### 浏览器输入 gin.wuchenyangzhendehaoshuai.com
### 登录后台系统(用户名密码为admin/password),默认是/login页面,也可以直接/home跳过登录
![image](https://github.com/user-attachments/assets/312e875d-bbb2-4b79-9db9-7a0e69598c3a)
![image](https://github.com/user-attachments/assets/eb9e7e3a-55a8-4b3b-94e8-c2b6e9b55dc5)
### 浏览器输入 jaeger.wuchenyangzhendehaoshuai.com,选择相对应的服务
![image](https://github.com/user-attachments/assets/e0469779-7a7f-44e9-8a07-7f964a702e23)
![image](https://github.com/user-attachments/assets/f8e2b5be-d7f8-463a-a07f-b3d3fa5dd3fa)

