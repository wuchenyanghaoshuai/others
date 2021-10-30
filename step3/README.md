# 日志收集方案 log-pilot
```
基于阿里云的开源日志搜集方案,并且使用daemonset的方式部署到node节点上用来自定义搜集控制台或者固定目录下的日志文件
本次的yaml是自定义输出到后端到kafka然后由logstash过滤以后发送到es中,kibana直接读es的索引
在k8s集群內以daemonset的方式来进行部署
注：在阿里云上开通阿里云sls日志收集与arms应用链路监控，
	会导致自定义安装的log-polit日志收集无法收集到日志，
	默认的阿里云自定义前缀是aliyun_logs，
	会导致冲突，需要自定义 - name: PILOT_LOG_PREFIX 日志收集格式的前缀可以解决。 
```
```
#所有的细节都在我下面这个链接里面
# https://thoughts.teambition.com/share/6020ff282917070042b9d3b7#title=k8s日志收集
```
