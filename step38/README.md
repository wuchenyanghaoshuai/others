```shell
本文介绍如何在Gitlab-CI中使用kaniko构建镜像，并推送到ACR。
```

## 1 安装gitlab
```shell
1. (安装gitlab)以centos7为例
wget https://mirror.tuna.tsinghua.edu.cn/gitlab-ce/yum/el7/gitlab-ce-16.11.5-ce.0.el7.x86_64.rpm
rpm -ivh gitlab-ce-16.11.5-ce.0.el7.x86_64.rpm(如果在安装过程中出现缺少依赖的包，需要yum安装依赖包)

2. 修改gitlab的配置文件
vim /etc/gitlab/gitlab.rb
external_url 'http://192.168.140.30' (改为自己的ip地址)

3. 重新配置gitlab
gitlab-ctl reconfigure

4. 启动gitlab
gitlab-ctl start 

5. 查看gitlab的状态
gitlab-ctl status

6. 查看gitlab登录密码
cat /etc/gitlab/initial_root_password
```
## 2 安装gitlab-runner

```shell
1. 安装gitlab-runner(本次采用docker方式安装)
docker run -d --name gitlab-runner --restart always \
  -v /srv/gitlab-runner/config:/etc/gitlab-runner \
  -v /var/run/docker.sock:/var/run/docker.sock \
  gitlab/gitlab-runner:latest

2. 注册gitlab-runner(这个只是一个命令，需要在gitlab中执行，需要给runner打一个kubernetes的标签)
docker exec -it gitlab-runner gitlab-runner register
3. 查看gitlab-runner的状态
docker exec -it gitlab-runner gitlab-runner status

4. gitlab-runner安装kubectl
curl -LO "https://dl.k8s.io/release/v1.22.1/bin/linux/amd64/kubectl"
chmod +x kubectl &&  mv kubectl /usr/local/bin/kubectl
kubectl version --client

5. 配置config
su - gitlab-runner
mkdir .kube && cd .kube
vim config (将k8s的config文件复制到config文件中)

6. 测试
kubectl get nodes

7. 后话
如果有条件的话建议把这个runner封装一下把k8s的配置信息封装进去，后续如果再创建多个runner的话，就可以直接注册到gitlab中就可以了

8. runner注册可以参考 4.3-4.5 部分
https://www.cnblogs.com/wjx001/p/18199274
```

## 3 准备kaniko与阿里云CR的 secret
``` shell
1. 准备kaniko的secret(使用此secret将build好的镜像推送到阿里云CR中)
该config.json 就是dockerlogin以后所产生的,所以不管是华为云cr或阿里云cr或者是本地的harbor,都可以使用此secret
kubectl create secret generic kaniko-secret --from-file=./config.json


cat /root/.docker/config.json
{
	"auths": {
	  "registry.cn-hangzhou.aliyuncs.com": {
		"auth": "dGI1NjkwODU5MDpxxxxxxxxxxxxxx"
	  }
	}
}

```
## 4 准备kaniko的测试文件

```shell
1. 这块主要是为了测试kaniko是否能build image并且push 到阿里云CR中
2. 准备一个python的文件(因为我们的python文件是在本地的，而且volume又是hostpath)，需要将代码复制到两个node下
3. 下面的python源代码可以直接用我的,我的image也是国内的，yaml这块只需要修改一下image的地址改成自己的
4. yaml中的volume中的path需要改成自己的或者用我的都行，secret 就是刚刚使用docker 的config.json 生成的secret
5. 这个文件就是测试看你的kaniko是否能build image并且push 到阿里云CR中，如果不能则需要手动修改问题，这一步是整个流程中重要的一步
[root@kn1 python]# pwd
/root/python
[root@kn1 python]# ls
app.py  Dockerfile  requirements.txt
[root@kn1 python]# cat app.py
#!/usr/bin/python3
from flask import Flask
app = Flask(__name__)

@app.route('/')
def hello_world():
    return 'Hello, Kaniko!'

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
[root@kn1 python]# cat Dockerfile
FROM registry.cn-hangzhou.aliyuncs.com/wuchenyang/python:3.9-slim-buster

WORKDIR /app


COPY . /app
RUN pwd && ls -l
RUN pip3 install -r requirements.txt


CMD ["python3", "app.py" ]
[root@kn1 python]# cat requirements.txt
flask==3.1.0

```


```shell
cat kaniko-pod.yaml
```

```yaml

apiVersion: v1
kind: Pod
metadata:
  name: kaniko
spec:
  containers:
  - name: kaniko
    image: gcr.io/kaniko-project/executor:latest
    imagePullPolicy: IfNotPresent
    args:
    - "--dockerfile=/app/Dockerfile"
    - "--context=/app"
    - "--destination=registry.cn-hangzhou.aliyuncs.com/wuchenyang/kaniko-python-docker:v2.0"
    volumeMounts:
    - name: kaniko-secret
      mountPath: /kaniko/.docker
    - name: workspace
      mountPath: /app
  restartPolicy: Never
  volumes:
  - name: kaniko-secret
    secret:
      secretName: kaniko-secret
  - name: workspace
    hostPath:
      path: /root/python
```
## 5 准备CI的测试源代码(前后端)

```shell
1. java-demo 是后端的源代码 git@github.com:wuchenyanghaoshuai/java-demo.git(用release/v0.0.2的分支，java1.8去build)
2. tupian-demo  是前端的源代码 git@github.com:wuchenyanghaoshuai/tupian.git(用test分支)
3. 将上面的代码clone到你本地并且上传到你私有的gitlab中,.gitlab-ci.yml文件我贴到下面了，根据自己的仓库地址做修改，其余不用动，源代码只有push到test或者prod分支才会触发build
4. 这两个项目都需要在各自的项目中开启使用runner
5，GITLAB_PAT，这个是一个全局变量,创建步骤如下(如果不创建全局变量，则需要给每个项目都创建同样的变量)
    5.1 创建个人访问令牌(PAT)
        5.1.1 点击右上角的头像，选择 Edit Profile（编辑个人资料）
        5.1.2 点击左侧的 Access Tokens（访问令牌）
        5.1.3 点击 Create personal access token（创建个人访问令牌）
        5.1.4 填写 Token 名称和过期时间，然后点击 Create personal access token（创建个人访问令牌）
        5.1.5 权限给上所有权限 
        5.1.6 就会给出一个token
    5.2 创建全局CICD变量
        5.2.1 点击右上角的头像 编辑个人资料
        5.2.2 点击下侧的 CI/CD -> Settings -> CI/CD Variables
        5.2.3 点击右上角的 New Variable
        5.2.4 填写 Key (GITLAB_PAT)和 Value (<YOUR_GITLAB_ACCESS_TOKEN>) 
        5.2.5 保护变量隐藏变量不勾选
        5.2.6 保存即可
6. 我在下面的context中使用了gitee的仓库，因为我本身的代码仓不支持https，所以只能用外部的，如果你们有https的话，则不需要进行修改，做测试的话可以直接用我的，没任何权限


```

```yaml
1. 这个是java-demo的gitlab-ci.yml文件
cat .gitlab-ci.yml

default:
  tags:
    - kubernetes # 使用 k8s 的 runner

stages:
  - build
  - deploy

build:
  stage: build
  script:
    - echo "我是 builddddddd"
    - echo "${CI_JOB_TOKEN}"
    - echo "我是 builddddddd"
    - sleep 10
  rules:
    - if: '$CI_COMMIT_BRANCH == "test" || $CI_COMMIT_BRANCH == "prod"'

deploy:
  stage: deploy
  script:
    # 调试信息
    - echo "当前用户:$(whoami)  HOME目录:$HOME"
    - echo "重新给gitlab-runner用户设置k8s权限"
    - ls -l $HOME/.kube/config || echo "Kubeconfig 未找到"
    - export KUBECONFIG=$HOME/.kube/config
    - kubectl config view
    - kubectl get node

    # 计算短 commit ID
    - echo "CI_COMMIT_SHA:$CI_COMMIT_SHA"  # 调试，查看完整 SHA
    - COMMIT_SHORT=$(echo $CI_COMMIT_SHA | cut -c1-6)  # 截取前 6 位
    - echo "Short commit ID:$COMMIT_SHORT"

    # 触发 Kaniko 构建镜像的 Job
    - echo "🚀 在 Kubernetes 上触发 Kaniko 构建镜像 🚀"
    - |
      kubectl apply -f - <<EOF
      apiVersion: batch/v1
      kind: Job
      metadata:
        name: kaniko-build-${CI_COMMIT_REF_NAME}-${COMMIT_SHORT}
        namespace: default
      spec:
        backoffLimit: 1
        ttlSecondsAfterFinished: 300
        template:
          spec:
            containers:
            - name: kaniko
              image: gcr.io/kaniko-project/executor:latest
              imagePullPolicy: IfNotPresent
              args:
                - "--dockerfile=Dockerfile"
                - "--context=git://gitee.com/chenyang_wu/java-demo.git#refs/heads/${CI_COMMIT_REF_NAME}"
                - "--destination=registry.cn-hangzhou.aliyuncs.com/wuchenyang/java-demo:${CI_COMMIT_REF_NAME}-${COMMIT_SHORT}"
                - "--insecure"
                - "--verbosity=debug"
              volumeMounts:
              - name: kaniko-secret
                mountPath: /kaniko/.docker
            restartPolicy: Never
            volumes:
            - name: kaniko-secret
              secret:
                secretName: kaniko-secret
      EOF

    # 等待 Kaniko Job 完成
    - echo "等待 Kaniko Job 完成..."
    - kubectl wait --for=condition=complete job/kaniko-build-${CI_COMMIT_REF_NAME}-${COMMIT_SHORT} --timeout=600s || kubectl wait --for=condition=failed job/kaniko-build-${CI_COMMIT_REF_NAME}-${COMMIT_SHORT} --timeout=600s

    # 使用 sed 修改 deploy.yaml 中的镜像 Tag
    - echo "🚀 修改 deploy.yaml 并应用到 Kubernetes 🚀"
    - sed -i "s|registry.cn-hangzhou.aliyuncs.com/wuchenyang/java-demo:.*|registry.cn-hangzhou.aliyuncs.com/wuchenyang/java-demo:${CI_COMMIT_REF_NAME}-${COMMIT_SHORT}|" deploy.yaml
    - cat deploy.yaml  # 调试：查看修改后的文件内容
    - kubectl apply -f deploy.yaml

    # 配置 Git 并推送修改后的 deploy.yaml
    - echo "🚀 将修改后的 deploy.yaml 推回 Git 仓库 🚀"
    - git config --global user.email "gitlab-runner@example.com"
    - git config --global user.name "GitLab Runner"
    - git add .  # 先暂存修改
    - git commit -m "Update deploy.yaml with new image tag:${CI_COMMIT_REF_NAME}-${COMMIT_SHORT}" || echo "无变更需要提交"  # 提交修改
    - git checkout $CI_COMMIT_BRANCH  # 再切换分支
    - git status  # 调试：查看当前状态
    - git pull origin $CI_COMMIT_BRANCH --rebase  # 拉取最新代码并 rebase
    - git push http://gitlab:$GITLAB_PAT@192.168.140.30/backend/java-demo.git $CI_COMMIT_BRANCH

    # 验证 Deployment 滚动更新
    - echo "验证 Deployment 更新状态..."
    - kubectl rollout status deployment/java-demo --namespace=default --timeout=300s

  rules:
    - if: '$CI_COMMIT_BRANCH == "test" || $CI_COMMIT_BRANCH == "prod"'


```

```yaml
2. 这个是tupian的gitlab-ci.yml文件
default:
  tags:
    - kubernetes # 使用 k8s 的 runner

stages:
  - build

build:
  stage: build
  script:
    # 调试信息
    - echo "当前用户:$(whoami)  HOME目录:$HOME"
    - echo "重新给gitlab-runner用户设置k8s权限"
    - ls -l $HOME/.kube/config || echo "Kubeconfig 未找到"
    - export KUBECONFIG=$HOME/.kube/config
    - kubectl config view
    - kubectl get node

    # 计算短 commit ID
    - echo "CI_COMMIT_SHA:$CI_COMMIT_SHA"  # 调试，查看完整 SHA
    - COMMIT_SHORT=$(echo $CI_COMMIT_SHA | cut -c1-6)  # 截取前 6 位
    - echo "Short commit ID:$COMMIT_SHORT"

    # 触发 Kaniko 构建镜像的 Job
    - echo "🚀 在 Kubernetes 上触发 Kaniko 构建镜像 🚀"
    - |
      kubectl apply -f - <<EOF
      apiVersion: batch/v1
      kind: Job
      metadata:
        name: kaniko-build-${CI_COMMIT_REF_NAME}-${COMMIT_SHORT}
        namespace: default
      spec:
        backoffLimit: 1
        ttlSecondsAfterFinished: 300
        template:
          spec:
            containers:
            - name: kaniko
              image: gcr.io/kaniko-project/executor:latest
              imagePullPolicy: IfNotPresent
              args:
                - "--dockerfile=Dockerfile"
                - "--context=git://gitee.com/chenyang_wu/tupian-demo.git#refs/heads/${CI_COMMIT_REF_NAME}"
                - "--destination=registry.cn-hangzhou.aliyuncs.com/wuchenyang/tupian-demo:${CI_COMMIT_REF_NAME}-${COMMIT_SHORT}"
                - "--insecure"
                - "--verbosity=debug"
              volumeMounts:
              - name: kaniko-secret
                mountPath: /kaniko/.docker
            restartPolicy: Never
            volumes:
            - name: kaniko-secret
              secret:
                secretName: kaniko-secret
      EOF

    # 等待 Kaniko Job 完成
    - echo "等待 Kaniko Job 完成..."
    - kubectl wait --for=condition=complete job/kaniko-build-${CI_COMMIT_REF_NAME}-${COMMIT_SHORT} --timeout=600s || kubectl wait --for=condition=failed job/kaniko-build-${CI_COMMIT_REF_NAME}-${COMMIT_SHORT} --timeout=600s

    # 使用 sed 修改 deploy.yaml 中的镜像 Tag
    - echo "🚀 修改 deploy.yaml 并应用到 Kubernetes 🚀"
    - sed -i "s|registry.cn-hangzhou.aliyuncs.com/wuchenyang/tupian-demo:.*|registry.cn-hangzhou.aliyuncs.com/wuchenyang/tupian-demo:${CI_COMMIT_REF_NAME}-${COMMIT_SHORT}|" deploy.yaml
    - cat deploy.yaml  # 调试：查看修改后的文件内容
    - kubectl apply -f deploy.yaml

    # 配置 Git 并推送修改后的 deploy.yaml
    - echo "🚀 将修改后的 deploy.yaml 推回 Git 仓库 🚀"
    - git config --global user.email "gitlab-runner@example.com"
    - git config --global user.name "GitLab Runner"
    - git add deploy.yaml  # 先暂存修改
    - git commit -m "Update deploy.yaml by CI with new image tag:${CI_COMMIT_REF_NAME}-${COMMIT_SHORT} [skip ci]" || true  # 提交修改
    - git checkout $CI_COMMIT_BRANCH  # 切换到目标分支
    - sed -i "s|registry.cn-hangzhou.aliyuncs.com/wuchenyang/tupian-demo:.*|registry.cn-hangzhou.aliyuncs.com/wuchenyang/tupian-demo:${CI_COMMIT_REF_NAME}-${COMMIT_SHORT}|" deploy.yaml
    - git add deploy.yaml  # 先暂存修改
    - git commit -m "Update deploy.yaml by CI with new image tag:${CI_COMMIT_REF_NAME}-${COMMIT_SHORT} [skip ci]" || true  # 提交修改
    - git pull origin $CI_COMMIT_BRANCH
    - git push http://gitlab:$GITLAB_PAT@192.168.140.30/front/tupian-demo.git $CI_COMMIT_BRANCH
    - echo "🚀 强制推送完成 🚀"
    # 验证 Deployment 滚动更新
    - echo "验证 Deployment 更新状态..."
    - kubectl rollout status deployment/tupian-demo --namespace=default --timeout=300s
    - echo "🚀 验证完成 🚀"

  rules:
    - if: '$CI_COMMIT_BRANCH == "test" || $CI_COMMIT_BRANCH == "prod"'
      changes:
        - "Dockerfile"
        - "**/src/**/*"
        - "!deploy.yaml"


```


## 6 下面我贴上你把源代码push到test分支以后gitlab-ci做的一些事情
