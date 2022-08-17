location /MP_verify_c36BZv71rxv4iB1O.txt {
add_header Content-Type 'text/plain;charset=utf-8';
return 200 'c36BZv71rxv4iB1O';
}


location /apple-app-site-association {
        charset UTF-8;
        default_type application/octet-stream;
        return 200  '{\"applinks\": {\"apps\": [],\"details\": [{\"appID\": \"XKPY665JK8.cn.porsche.app\",\"paths\": [ \"*\"]},{\"appID\": \"5DUT5478S6.cn.porsche.product.development\",\"paths\": [ \"*\"]}]}}' ;
    }
	
location /.well-known/apple-app-site-association {
        charset UTF-8;
        default_type application/octet-stream;
        return 200  '{\"applinks\": {\"apps\": [],\"details\": [{\"appID\": \"VSK98DGSS7.com.miduiedu.mdeapp\",\"paths\": [ \"/midui/*\"]}]}}' ;
    }


location /abc {
        alias /data/www/default/abc;
        index index.html index.htm;
        try_files $uri $uri /abc/index.html;
	}

location /def {
        alias /data/www/default/def;
        index index.html index.htm;
        try_files $uri $uri /def/index.html;
        }
