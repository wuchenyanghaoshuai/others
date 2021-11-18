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
