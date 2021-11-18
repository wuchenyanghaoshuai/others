# velero 概述
```
Velero（以前称为 Heptio Ark）为您提供了用于备份和恢复 Kubernetes 集群资源和持久卷的工具。您可以通过云提供商或本地运行 Velero。Velero 让您：
1. 备份集群并在丢失时恢复。
2. 将集群资源迁移到其他集群。
3. 将您的生产集群复制到开发和测试集群
```
# velero工作原理
```
每个 Velero 操作——按需备份、计划备份、恢复——都是一个自定义资源，使用 Kubernetes 自定义资源定义 (CRD) 定义并存储在 etcd 中。Velero 还包括处理自定义资源以执行备份、恢复和所有相关操作的控制器。
您可以备份或恢复集群中的所有对象，也可以按类型、命名空间和/或标签过滤对象。
Velero 非常适合灾难恢复用例，以及在集群上执行系统操作（如升级）之前对应用程序状态进行快照。
```
# 备份工作流程
```
1. Velero 客户端调用 Kubernetes API 服务器来创建一个Backup对象。
2. 该BackupController注意到新的Backup对象，并进行验证。
3. 在BackupController开始备份过程。它通过向 API 服务器查询资源来收集要备份的数据。
4. 将BackupController使得对象存储服务的调用-例如，AWS S3 -上传备份文件
```
![image](https://user-images.githubusercontent.com/39818267/142493899-5fdc5c99-ec7e-4259-bff3-8a4fb5bfb796.png)
# 安装minio
```
注意此块的--volume 要跟自己环境一致,然后打开ip:9000,输入账号密码即可登录,登录上去自己创建一个bucket
docker run -d  --name minio     --publish 9000:9000     --publish 9001:9001     --env MINIO_ACCESS_KEY="admin"     --env MINIO_SECRET_KEY="admin123"     --volume /chenyang:/data     bitnami/minio:latest
```
![image](https://user-images.githubusercontent.com/39818267/142497549-1900982e-27fd-4840-ae58-084135161375.png)
# 安装velero
```
1. 下载安装velero
注意请确保你的.kube下有可以正常链接到k8s集群的config
wget https://github.com/vmware-tanzu/velero/releases/download/v1.7.0/velero-v1.7.0-linux-amd64.tar.gz
tar -xf velero-v1.7.0-linux-amd64.tar.gz
mv velero-v1.7.0-linux-amd64/velero /usr/local/bin
velero version (查看到版本就可以了,看不到的话应该是没指定k8s集群)
2. 创建minio凭证
vim  credentials-velero
[default]
aws_access_key_id = admin
aws_secret_access_key = admin123
3.安装velero
   velero install    \
     --provider aws   \
     --bucket k8s-backup   \
     --image velero/velero:v1.6.3  \
     --plugins velero/velero-plugin-for-aws:v1.2.1  \
     --namespace velero  \
     --secret-file ./credentials-velero  \
     --use-volume-snapshots=false \
     --use-restic \
     --backup-location-config region=minio,s3ForcePathStyle="true",s3Url=http://192.168.1.213:9000
4. 查看k8s集群内是否多一个namespace,然后查看对应的pod日志,看有没有报错,没有报错继续进行
5. 对某一个namespace进行备份,完成以后立即到minio的页面查看是否生成对应的目录
   velero backup create test-backup --include-namespaces test  --wait
6. 到k8s集群删除namespace，并且恢复,restore完成以后去集群查看,发现刚刚删除的kkb-test命名空间又存在了,并且资源都是正确的
   kubectl delete ns kkb-test
   velero restore create --from-backup test-backup
```
![image](https://user-images.githubusercontent.com/39818267/142497611-f4a7ed76-ad2d-4639-9a29-2bdeff6f65e1.png)
