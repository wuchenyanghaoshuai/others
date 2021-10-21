#1.备份mongodb内的yapi database
  (1)登录到mongo数据库内 
     # docker exec -it  ec784143af58 bash
  (2)备份数据库到目录
     # mongodump -h 127.0.0.1:27017 -d yapi -o /root/data/

2.启动一个新mongodb

  （1）docker volume create mongo-data
  （2）docker run -d \
  --name mongo-yapi \
  -v mongo-data:/data/db \
  -e MONGO_INITDB_ROOT_USERNAME=anoyi \
  -e MONGO_INITDB_ROOT_PASSWORD=anoyi.com \
  mongo

3.恢复刚刚备份的数据库到新的mongodb中
  (1)mongorestore -h 127.0.0.1:27017  -d yapi --drop  --dir /root/data/yapi
  注意: --drop 慎用; --drop:恢复的时候，先删除当前数据，然后恢复备份的数据。就是说，恢复后，备份后添加修改的数据都会被删除
  (2)mongorestore -h 127.0.0.1:27017  --authenticationDatabase admin    -u anoyi -p anoyi.com -d yapi --drop /yapi/

https://www.jianshu.com/p/de26ad0110b0


4.自定义yapi配置文件
	vim config.json
	```
{
  "port": "3000",
  "adminAccount": "admin@anoyi.com",
  "timeout":120000,
  "db": {
    "servername": "mongo",
    "DATABASE": "yapi",
    "port": 27017,
    "user": "anoyi",
    "pass": "anoyi.com",
    "authSource": "admin"
  },
  "ldapLogin": {
      "enable": true,
      "server": "ldap://ip:389",
      "baseDn": "cn=xxx,dc=xxx,dc=xxx",
      "bindPassword": "xxxx",
      "searchDn": "ou=xxx,dc=xxx,dc=xxx",
      "searchStandard": "&(objectClass=posixAccount)(cn=%s)"
   }
}
	```
#具体看这个链接https://hellosean1025.github.io/yapi/devops/index.html#%e6%9c%8d%e5%8a%a1%e5%99%a8%e7%ae%a1%e7%90%86
5.初始化yapi数据库索引以及管理员账号
```
docker run -it --rm \
  --link mongo-yapi:mongo \
  --entrypoint npm \
  --workdir /yapi/vendors \
  -v $PWD/config.json:/yapi/config.json \
  registry.cn-hangzhou.aliyuncs.com/anoyi/yapi \
  run install-server
```
6.docker启动yapi
docker run -d \
  --name yapi \
  --link mongo-yapi:mongo \
  --workdir /yapi/vendors \
  -p 3000:3000 \
  -v $PWD/config.json:/yapi/config.json \
  registry.cn-hangzhou.aliyuncs.com/anoyi/yapi \
  server/app.js
7.访问IP:Port看是否出现登录界面以及ldap集成
https://www.jianshu.com/p/a97d2efb23c5



