# 高版本的k8s在部署nfs-client-provisioner 的时候会报错,本人1.20版本 虽然pod起来了但是在实际使用storageclass的时候还是有问题的
# 如果你创建一个使用storageclass的deployment 那么这个pod会pending，接着describe deploy 会报错 waiting for a volume to be created, either by external provisioner “fuseim.pri/ifs” or manually created by system administrator
# 1.20版本（我的是1.20.0）默认禁止使用selfLink
# 开启方法 - --feature-gates=RemoveSelfLink=false (kube-apiserver.yaml)
