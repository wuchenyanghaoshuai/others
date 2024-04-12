```
具体看本文提供的图片或者链接
本文参考: https://blog.csdn.net/liuhongqwe/article/details/136037589
```


```
如果需要域名访问则需要修改一下两个方面
nginx and confluence的conf
```


```
nginx方面,如果你的端口不修改的话那么基本上就是我这个配置了
# cat wiki.ailieyun.com.conf
upstream lieyun-wiki {
    server 192.168.3.127:8090  weight=5 max_fails=5 fail_timeout=30s;
}

server {
    listen      443 ssl;
    server_name wiki.ailieyun.com;

    location / {
	client_max_body_size 100m;
        proxy_pass              http://lieyun-wiki/;
        proxy_set_header X-forwarded-Proto $scheme;
        proxy_set_header        Host            $host;
        proxy_set_header        X-Real-IP       $remote_addr;
        proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
server {
    listen 80;
    server_name  wiki.ailieyun.com;
    rewrite ^(.*)$ https://$host$1 permanent;
}
```


```
confluence conf配置方面
找到confluence的安装目录
以文档来说 就是/opt/atlassian/confluence/conf/server.xml 修改之前要先备份
然后只增加红色部分就可以了根据自己公司域名需要
scheme="https" secure="true" proxyName="wiki.ailieyun.com" proxyPort="443"
然后useHttpOnly="true" 去掉
```
![image](https://github.com/wuchenyanghaoshuai/others/assets/39818267/2820316f-8137-45bc-9727-454875c21dcc)
