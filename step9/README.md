```
yum install -y gcc gcc-c++ zlib zlib-devel openssl openssl-devel pcre pcre-devel jemalloc jemalloc-devel git


wget http://tengine.taobao.org/download/tengine-2.3.1.tar.gz
wget https://kkb-op.oss-cn-beijing.aliyuncs.com/kkb-test/openssl-OpenSSL_1_1_1b.tar.gz
git clone https://github.com/vozlt/nginx-module-vts.git

yum install -y gcc gcc-c++ zlib zlib-devel openssl openssl-devel pcre pcre-devel jemalloc jemalloc-devel


cd /root/ && mkdir -p /data/nginx
tar zxf openssl-OpenSSL_1_1_1b.tar.gz
mv openssl-OpenSSL_1_1_1b /usr/local/openssl2

tar zxf tengine-2.3.1.tar.gz && cd tengine-2.3.1

./configure --prefix=/data/nginx \
--with-http_gzip_static_module \
--with-pcre \
--with-http_ssl_module \
--with-http_realip_module \
--with-http_mp4_module \
--with-http_v2_module \
--with-http_sub_module \
--with-openssl=/usr/local/openssl2 \
--add-module=modules/ngx_http_upstream_check_module \
--add-module=/root/nginx-module-vts

make -j8 && make install
```
