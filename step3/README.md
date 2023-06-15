# 日志收集方案 filebeat
# 日志收集方案 log-pilot
```
如果runtime是docker 就使用log-pilot;如果是containerd就使用filebeat，错误使用的话会导致出现runtime的问题问题
具体体现在docker.sock 和containerd.sock 那块
```