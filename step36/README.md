# 阿里云ECS自建k8s并且使用阿里云的slb

```
step1
  首先要有一个公网slb，并且后端监听 80 443，因为ingress-nginx 使用hostnetrok方式 所以自动会在节点监听80 443
step2
  首先根据自己的k8s集群版本选择自己想对应的ingress-nginx的版本,具体参考下面的链接
  https://github.com/kubernetes/ingress-nginx
step3
  修改一下几个位置并说明原因
  1 修改svc 字段下的 externalIPs ,这块修改成阿里云购买的slb 因为本身svc就是loadbalancer的模式所以不用修改
  2 deploy那块修改成hostnetwork: true 和 nodename ，防止deployment 的pod漂移
step4
  下面放修改的截图
```


