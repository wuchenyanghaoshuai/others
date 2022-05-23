1.系统规划(/etc/hosts)

        192.168.100.125    etcd-1

        192.168.100.126    etcd-2

        192.168.100.127    etcd-3

2.关闭selinux

    setenforce 0
    sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
3.配置系统的文件描述符
    echo "* soft nofile 65536" >> /etc/security/limits.conf
    echo "* hard nofile 65536" >> /etc/security/limits.conf
    echo "* soft nproc 65536"  >> /etc/security/limits.conf
    echo "* hard nproc 65536"  >> /etc/security/limits.conf
    echo "* soft  memlock  unlimited"  >> /etc/security/limits.conf
    echo "* hard memlock  unlimited"  >> /etc/security/limits.conf
4.下载etcd二进制包
    wget https://github.com/etcd-io/etcd/releases/download/v3.4.16/etcd-v3.4.16-linux-amd64.tar.gz
5.解压安装
    tar xf etcd-v3.4.1-linux-amd64.tar.gz
    cd etcd-v3.4.1-linux-amd64
    cp etcd etcdctl /usr/bin/
6.添加用户和组并且授权
    groupadd --system etcd
    useradd -s /sbin/nologin --system -g etcd etcd
    mkdir -p /var/lib/etcd
    chown -R etcd:etcd /var/lib/etcd
7.添加几个环境变量

    INT_NAME="ens192"   #注意这里，每个机器的网卡号不同
    ETCD_HOST_IP=$(ip addr show $INT_NAME | grep "inet\b" | awk '{print $2}' | cut -d/ -f1)
    ETCD_NAME=$(hostname -s)
8.配置system启动(注意下面的http部分要修改为对应的主机名)
cat <<EOF | sudo tee /etc/systemd/system/etcd.service
[Unit]
Description=etcd service
Documentation=https://github.com/etcd-io/etcd


[Service]
Type=notify
User=etcd
ExecStart=/usr/bin/etcd \\
--name ${ETCD_NAME} \\
--data-dir=/var/lib/etcd \\
--initial-advertise-peer-urls http://${ETCD_HOST_IP}:2380 \\
--listen-peer-urls http://${ETCD_HOST_IP}:2380 \\
--listen-client-urls http://${ETCD_HOST_IP}:2379,http://127.0.0.1:2379 \\
--advertise-client-urls http://${ETCD_HOST_IP}:2379 \\
--initial-cluster-token etcd-cluster-0 \\
--initial-cluster etcd-1=http://etcd-1:2380,etcd-2=http://etcd-2:2380,etcd-3=http://etcd-3:2380 \\
--initial-cluster-state new \


[Install]
WantedBy=multi-user.target
EOF

9.启动etcd

    systemctl enable etcd

    systemctl start etcd （只有两台或者以上的服务启动才算正常,如果只启动一台,会一直卡住）

10. 常用命令

 查看集群中成员列表

   etcdctl member list
查看节点信息
   etcdctl endpoint --cluster=true status

   etcdctl endpoint --cluster=true

   health etcdctl endpoint --cluster=true hashkv

列表方式输出
    etcdctl --write-out=table --cluster=true endpoint status

参考文档
 https://www.58jb.com/html/install-etcd-on-centos7.html




