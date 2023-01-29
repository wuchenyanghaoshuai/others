# 搭建个人网盘存储 docker-compose
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
# 启动以后需要修改php.config 下的两个内容(如果需要使用https访问就改如果使用ip加端口就不用改)
```
 'overwriteprotocol' => 'https',
 # array 这块是自己添加的
 array (
    0 => 'file.ailieyun.com',
  ),

 ```