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
step5
  下面附上我自己的yaml，大家根据我的集群版本修改哈，我是1.23.1版本
  另外还有测试的ingress和yaml，有证书的话就保留tls字段，没有的话就删除
```
![image](https://github.com/wuchenyanghaoshuai/others/assets/39818267/1eece38c-8fb2-493d-89f9-3a4925789f07)
![image](https://github.com/wuchenyanghaoshuai/others/assets/39818267/200605a8-621f-4933-82ba-f55e21750776)
![image](https://github.com/wuchenyanghaoshuai/others/assets/39818267/7bdc26b5-3cd2-4501-b787-7142cf40dc5b)
![image](https://github.com/wuchenyanghaoshuai/others/assets/39818267/efb52877-5741-465d-bd69-2b501e224fef)


