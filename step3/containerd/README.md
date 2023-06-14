# 日志收集方案 filebeat
```
	这套方案为基于namespace来收集日志，比如我A namespace里跑的业务，或者有相关的需求为只收集A namespace的日志，
	架构为 filebeat收集namespace日志，然后发送到kafka，经过logstash以后，到es中，主要收集和相关的数据清洗工作都在filebeat侧
```

# 配置文件cm解释

```
    filebeat.inputs:
    - type: container
      enabled: true
      paths:
        - /var/log/containers/*_ailieyun-prod_*.log
      multiline.pattern: '^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}|^[1-9]\d*\.[1-9]\d*\.[1-9]\d*\.[1-9]\d*'
      multiline.negate: true
      multiline.match: after
	  multiline.timeout: 10s
      fields: #主动添加字段，方便展示
        namespace: ailieyun-prod
        env: prod
        k8s: cluster-prod
      filebeat.config.modules:
        path: ${path.config}/modules.d/*.yml
        reload.enabled: false
    processors: #不添加这个的话，不会显示具体的pod的信息，例如pod的名称啊之类的
      - add_kubernetes_metadata:
          host: ${NODE_NAME}
          matchers:
          - logs_path:
              logs_path: "/var/log/containers/"
      - drop_fields:

          fields:  #下面的是属于丢弃字段，方便在es中根据字段查询，不丢弃的话在es查询中显得不中用的字段特别多，所以就丢弃了达到极简的风格
            - host
            - ecs
            - input
            - log
            - stream
            - container
            - agent
            - kubernetes.pod.uid
            - kubernetes.replicaset.name
            - kubernetes.labels.app
            - kubernetes.labels.pod-template-hash
          ignore_missing: true
    output.kafka:       #输出到kafka
      hosts: ["192.168.10.136:9092"]
      topic: '%{[fields.k8s]}-%{[fields.namespace]}' #自动创建topic，跟我我上面得配置，我创建的topic为cluster-prod-ailieyun-prod
```


```
#所有的细节都参考下面的链接,
# https://www.modb.pro/db/168772
```
