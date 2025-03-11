注意:
    1. 因为我是macos环境，所以需要指定平台为linux/amd64
    2. 源代码里有sql文件,所以需要在数据库里执行一下sql文件
    3. 需要修改main函数里的mysql数据库的连接地址以及用户名密码
    4. 默认登录后台的用户名密码是admin/password
    5. 打包，上传阿里云镜像

第一步 打包可执行文件
GOOS=linux GOARCH=amd64 go build -o main
第二步
docker build  --platform linux/amd64 -f Dockerfile -t registry.cn-hangzhou.aliyuncs.com/wuchenyang/gindemo:v1 .
