![image](https://github.com/user-attachments/assets/b57eb1d8-5a45-4ce5-98d3-e95a84df7ec2)![image](https://github.com/user-attachments/assets/7232a569-a652-4171-8f3a-3430e2df4295)# 本文介绍如何使用Opentelemetry+jaeger实现分布式链路追踪。
## 1. 安装CertManager 首选
```shell
# 地址在 https://cert-manager.io/docs/installation/kubectl
# 请注意我部署的时候kubectl version >=1.19.0
# kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.17.0/cert-manager.yaml
# 可以先wget下来，然后查看yaml里的具体内容然后再apply
# 如果你的空k8s版本低于1.19.0 可以去github里寻找低版本的yaml，一般会有版本对照表
```

## 2. 安装Jaeger all-in-one模式 step01
```yaml 
# 该yaml与step01提供的yaml区别在于使用了外部的存储es，具体看env部分

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
```

```yaml
# 我这里使用的ingress暴露的，你也可以使用nodeport或者ingress暴露
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


## 3. 安装opentelemetry-operator

```shell
#opentelemetry-operator与k8s版本有关，需要根据k8s版本进行安装，具体看具体看下面的地址
# https://github.com/open-telemetry/opentelemetry-operator/blob/main/docs/compatibility.md
# 如果你的k8s集群高于1.23 直接看这个md里的对照表就行，如果低于1.23，需要看下面的地址(本人1.22.1)里的对照关系
# https://github.com/open-telemetry/opentelemetry-operator/tree/v0.85.0
# 然后在tag里找到v0.85.0 的installation manifest for kubernetes 右键复制链接地址 ,然后在终端wget 即可
# 在终端执行
kubectl apply -f opentelemetry-operator.yaml
# 查看是否安装成功
kubectl get pod -n opentelemetry-operator-system

# 如果下载不了image可以使用我的镜像(我pull下来tag的)
#registry.cn-hangzhou.aliyuncs.com/wuchenyang/kube-rbac-proxy:v0.13.1
#registry.cn-hangzhou.aliyuncs.com/wuchenyang/opentelemetry-operator:0.85.0
#实际使用的时候也可以将镜像pull下来，tag成官方的image
```
## 4. 部署otel-collector 跟Instrumentation

```yaml

# cat 04otlp/Instrumentation.yaml
#这个因为我是java的demo所以下面那块写了java，如果是python或者go去官方照文档即可都是一样的，
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
        value: "http://otel-collector.default:4317"  # gRPC 端点（可选覆盖全局）
```

```yaml
# cat 04otlp/otlp.yaml
#这个配置就是otel-collector的配置，这里我使用的是jaeger的exporter，你也可以使用其他的exporter
#因为jager的后端没有配置es等存储，我这里使用的本机，所以数据保存不了30天往上
apiVersion: opentelemetry.io/v1alpha1
kind: OpenTelemetryCollector
metadata:
  name: otel
spec:
  config: |
    receivers:
      otlp:
        protocols:
          grpc:
          http:
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
          processors: []
          exporters: [otlp/jaeger]
```

## 5. 部署Java-Demo

```yaml
#可以直接使用我的yaml
#svc那块你可以直接用ingress或者nodeport都可以，我是用了ingress，本地做了域名映射
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



