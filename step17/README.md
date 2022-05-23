➜ ~ kkkb-dev get pods -n kkb-dev -o wide|grep gateway-trip
kkb-message-gateway-tripartite-56c64d458d-p9nr6 0/1 ContainerCreating 0 24m <none> 192.168.100.208 <none> <none>
kkb-message-gateway-tripartite-697595465c-sthqg 1/1 Running 1 5d18h 10.3.10.183 192.168.70.62 <none> <none>

Warning FailedCreatePodContainer 0s (x3 over 26s) kubelet, 192.168.100.208 unable to ensure pod container exists: failed to create container for [kubepods burstable pod9a36f621-62ee-11ec-9fa4-0050568f7de2] : mkdir /sys/fs/cgroup/memory/kubepods/burstable/pod9a36f621-62ee-11ec-9fa4-0050568f7de2: cannot allocate memory

现象: 再创建pod的过程中一直处于ContainerCreating状态，然后describe pod 发现上面的报错,说没有内存的可用，但是通过kuboard界面看到还是有很多可用内存,初步怀疑内存泄漏
手动到该containercreating的node节点去执行上面的mkdir 命令发现还是报错,所以确定是内存泄漏
cat  /etc/default/grub

BOOT_IMAGE=/vmlinuz-3.10.0-1160.6.1.el7.x86_64 root=UUID=b499b035-85b1-451d-977b-0a402abfba1e ro crashkernel=auto spectre_v2=retpoline ipv6.disable=1 rhgb quiet nouveau.modeset=0 cgroup.memory=nokmem
[root@kkb-hc-208 ~]# cat /etc/default/grub
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR="$(sed 's, release .*$,,g' /etc/system-release)"
GRUB_DEFAULT=saved
GRUB_DISABLE_SUBMENU=true
GRUB_TERMINAL_OUTPUT="console"
GRUB_CMDLINE_LINUX="crashkernel=auto spectre_v2=retpoline ipv6.disable=1 rhgb quiet"
GRUB_DISABLE_RECOVERY="true"



修改这个文件，修改以后的是

BOOT_IMAGE=/vmlinuz-3.10.0-1160.6.1.el7.x86_64 root=UUID=b499b035-85b1-451d-977b-0a402abfba1e ro crashkernel=auto spectre_v2=retpoline ipv6.disable=1 rhgb quiet nouveau.modeset=0 cgroup.memory=nokmem
[root@kkb-hc-208 ~]# cat /etc/default/grub
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR="$(sed 's, release .*$,,g' /etc/system-release)"
GRUB_DEFAULT=saved
GRUB_DISABLE_SUBMENU=true
GRUB_TERMINAL_OUTPUT="console"
GRUB_CMDLINE_LINUX="crashkernel=auto spectre_v2=retpoline ipv6.disable=1 rhgb quiet nouveau.modeset=0 cgroup.memory=nokmem"
GRUB_DISABLE_RECOVERY="true"



然后执行grub2-mkconfig -o /boot/grub2/grub.cfg

然后reboot该node节点

然后kubectl get node 就发现之前处于containerCreating的pod没有问题了

➜ ~ kkkb-dev get pods -n kkb-dev -o wide|grep gateway-trip
kkb-message-gateway-tripartite-56c64d458d-p9nr6 1/1 Running 2 47m 10.3.6.214 192.168.100.208 <none> <none>



文档https://1week.tistory.com/24
