# 日志收集方案 filebeat
# 日志收集方案 log-pilot
```
如果runtime是docker 就使用log-pilot，如果是containerd，就使用filebeat，错误使用的话会导致出现问题
```