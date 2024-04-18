```
➜  ~ kubectl get node
Unable to connect to the server: tls: failed to verify certificate: x509: certificate is valid for 10.1.0.1, 172.31.54.213, not 47.76.237.146
```


```
通过查阅资料发现了一个kubectl的参数--insecure-skip-tls-verify，加上这个参数之后确实好使了，但是，总是感觉治标不治本，所以经过一番查阅是apiserver的证书中没有添加10.8.5.5这个ip导致的，需要重新生成一下证书，具体操作如下：
```
# 查看apiserver证书信息
```
cd /etc/kubernetes/pki
openssl x509 -noout -text -in apiserver.crt |grep IP
                DNS:k8s-master, DNS:kubernetes, DNS:kubernetes.default, DNS:kubernetes.default.svc, DNS:kubernetes.default.svc.cluster.local, IP Address:10.96.0.1, IP Address:172.25.1.100
```


# 删除旧证书

```
mkdir -pv /opt/cert
mv apiserver.* /opt/cert

```


# 生成新的证书

```

kubeadm init phase certs apiserver --apiserver-advertise-address 172.31.54.213 --apiserver-cert-extra-sans 10.1.0.1  --apiserver-cert-extra-sans 47.76.237.146

```
# 重启所有机器的kubelet

```
systemctl restart kubelet.service
```

##本文参考
```
https://www.jianshu.com/p/2550dba811e7
```
