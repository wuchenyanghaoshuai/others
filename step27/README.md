# docker pull image

```
一个是ldap一个是ldap的web页面
docker pull osixia/openldap:latest or  docker pull osixia/openldap:1.4.0
docker pull osixia/phpldapadmin:latest
```


# docker 运行ldap

```
docker run -itd -p 389:389 --name openldap \
--env LDAP_ORGANISATION="ailieyun" \
--env LDAP_DOMAIN="ailieyun.com" \
--env LDAP_ADMIN_PASSWORD="123456" \
--detach osixia/openldap:1.4.0


# -v /opt/data/ldap/database:/var/lib/ldap \
 -v /opt/data/ldap/config:/etc/openldap/slapd.d \
#注意 -v 是将数据持久化以及共享配置文件 此处是做测试我没有使用-v参数

配置LDAP组织者：LDAP_ORGANISATION
配置LDAP域：LDAP_DOMAIN
配置LDAP密码：LDAP_ADMIN_PASSWORD
默认登录用户名：admin
```


# docker运行web页面

```
docker run -d --privileged -p 80:80 --name myphpldapadmin \
--env PHPLDAPADMIN_HTTPS=false \
--env PHPLDAPADMIN_LDAP_HOSTS=172.17.0.2 \
--detach osixia/phpldapadmin
# 注意PHPLDAPADMIN_LDAP_HOSTS这个字段修改为自己docker运行的ldap的container的ip地址

```

# 登陆ldap
```
IP地址: 127.0.0.1
login DN ：cn=admin,dc=ailieyun,dc=com
Password: 123456
```

#Windows 管理界面
```
http://directory.apache.org/studio/
```
