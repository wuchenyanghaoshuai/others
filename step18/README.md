kubernetes 集群在查看日志

kubectl logs -f <pod>

报错

failed to watch file "/var/log/pods/xxxx/xxxxxx.log": no space left on deviceroot

可能是节点 out of inotify watches  **节点**   **node**  **不是master**  **是pod运行的物理节点**

可以使用命令查看情况

sysctl fs.inotify.max_user_watches

默认是

fs.inotify.max_user_watches = 8192

可以尝试将数值改大

vi /etc/sysctl.conf

在文件中添加或者修改成自己想要的值，一般直接翻倍，一直到没有问题不报错了

fs.inotify.max_user_watches = 2097152

修改完后执行命令，读取修改后的值

/sbin/sysctl -p

再次查看就没有错了

kubectl logs -f <pod>
