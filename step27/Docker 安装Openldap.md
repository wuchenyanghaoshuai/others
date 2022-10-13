# Docker 安装Openldap

## ldap基本概念

> (一) 目录树概念

```shell
1. 目录树：在一个目录服务系统中，整个目录信息集可以表示为一个目录信息树，树中的每个节点是一个条目。

2. 条目：每个条目就是一条记录，每个条目有自己的唯一可区别的名称（DN）。

3. 对象类：与某个实体类型对应的一组属性，对象类是可以继承的，这样父类的必须属性也会被继承下来。

4. 属性：描述条目的某个方面的信息，一个属性由一个属性类型和一个或多个属性值组成，属性有必须属性和非必须属性。

```

> (二)  **DC、UID、OU、CN、SN、DN、RDN**

| 关键字 | 英文全称           | 含义                                                         |
| ------ | :----------------- | ------------------------------------------------------------ |
| dc     | Domain Component   | 域名的部分，其格式是将完整的域名分成几部分，如域名为example.com变成dc=example,dc=com（一条记录的所属位置） |
| uid    | User Id            | 用户ID songtao.xu（一条记录的ID）                            |
| Ou     | Organization Unit  | 组织单位，组织单位可以包含其他各种对象（包括其他组织单元），如“oa组”（一条记录的所属组织） |
| cn     | Common Name        | 公共名称，如“Thomas Johansson”（一条记录的名称）             |
| sn     | Surname            | 姓，如“许”                                                   |
| dn     | Distinguished Name | “uid=songtao.xu,ou=oa组,dc=example,dc=com”，一条记录的位置（唯一） |
| Ran    | Relative dn        | 相对辨别名，类似于文件系统中的相对路径，它是与目录树结构无关的部分，如“uid=tom”或“cn= Thomas Johansson” |

> (示例图)

<img width="470" alt="image" src="https://user-images.githubusercontent.com/39818267/195546352-bf07b60e-f93a-40d1-8d47-ce6938a1b0c7.png">

## 安装LDAP

```shell
docker pull osixia/openldap:1.4.0

mkdir /opt/data/ldap/database -p

mkdir /opt/data/ldap/config 

docker run -itd -p 389:389 --name openldap \
-v /opt/data/ldap/database:/var/lib/ldap \
-v /opt/data/ldap/config:/etc/openldap/slapd.d \
--env LDAP_ORGANISATION="xyz" \
--env LDAP_DOMAIN="xyz.com" \
--env LDAP_ADMIN_PASSWORD="123456" \
 osixia/openldap:1.4.0
```

- 配置LDAP组织者：`LDAP_ORGANISATION`
- 配置LDAP域：`LDAP_DOMAIN`
- 配置LDAP密码：`LDAP_ADMIN_PASSWORD`
- 默认登录用户名：`admin`

## LDAP管理工具安装

```shell
docker run -d --privileged -p 80:80 --name myphpldapadmin \
--env PHPLDAPADMIN_HTTPS=false \
--env PHPLDAPADMIN_LDAP_HOSTS=192.168.5.13 \
--detach osixia/phpldapadmin
```

<img width="872" alt="image" src="https://user-images.githubusercontent.com/39818267/195546445-b87dd66f-d869-41b2-af76-02bbba1afc72.png">

- `login DN` ：cn=admin,dc=xyz,dc=com
- `Password`: 123456
