#小命令
```
#设置ssh连接超时时间
sed -i 's/^#ClientAliveInterval.*/ClientAliveInterval 30/g' /etc/ssh/sshd_config
sed -i 's/^#ClientAliveCountMax.*/ClientAliveCountMax 3600/g' /etc/ssh/sshd_config
systemctl restart sshd
```
