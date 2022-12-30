vim compose.yml
```
version: '3'
services:
  nextcloud:
    image: nextcloud:latest
    ports:
      - 8080:80
    volumes:
      - /root/nextcloud/app:/var/www/html
    environment:
      - NEXTCLOUD_ADMIN_USER=admin
      - NEXTCLOUD_ADMIN_PASSWORD=password
      - NEXTCLOUD_REDIS_HOST=redis
      - NEXTCLOUD_DB_USER=nextcloud_user
      - NEXTCLOUD_DB_PASSWORD=password
    depends_on:
      - db
      - redis
  db:
    image: mariadb:latest
    restart: always
    command: --transaction-isolation=READ-COMMITTED --binlog-format=ROW
    environment:
      - MYSQL_ROOT_PASSWORD=123456
      - MYSQL_USER=nextcloud
      - MYSQL_PASSWORD=nextcloud
      - MYSQL_DATABASE=nextcloud
    volumes:
      - /root/nextcloud/db:/var/lib/mysql
  redis:
    image: redis:latest
    volumes:
      - /root/nextcloud/redis_data:/data
```

```
首先确认目录是否存在,没有的话就创建一下
然后docker-compose up -d
如果你不使用nginx作为代理的话，直接访问ip:port就不用看以下的文件,如果使用nginx作为代理，并且使用https的话请继续看下去

启动以后，去挂载的目录看下,修改添加config.php
 'overwriteprotocol' => 'https',

 并且在  array (
    0 => '127.0.0.1:8080',
    1 => 'example.com',
  ),
  example.com 为你自己的域名`
```