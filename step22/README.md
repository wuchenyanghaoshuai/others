```
1. 如果集群有storageclass忽略deploy下的nfs,因为这个是部署一个由nfs提供的storageclass
  1.1 如果没有storageclass的话,如果由mysql也可以，因为storageclass是为了作为mysql的存储来部署的
  1.2 如果你既没有mysql也没有storageclass，那就修改deploy下的nfs的地址,和nfs的路径就可以了，然后直接apply

2. 安装mysql
  2.1 deploy/mysql/persistent.yaml 创建一个pvc
  2.2 检查这个目录下的文件配置deploy/mysql/mysql-ceph.yaml,如果没有问题就apply就可以看到mysql运行

3. 安装 nacos
  3.1 /deploy/nacos/nacos-quick-start.yaml 这个位置需要注意是我修改了暴露方式,本来是clusterip
4. 注意
  4.1 目录下的其他文件，任选其一都可以
```
