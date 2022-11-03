## gitlab备份


### 检查自己的备份目录
```
root@test:~# cat /etc/gitlab/gitlab.rb | grep backup_path
# gitlab_rails['manage_backup_path'] = true
# gitlab_rails['backup_path'] = "/var/opt/gitlab/backups"
```


### 手动备份
```
gitlab-rake gitlab:backup:create
#crontab
crontab -e
 0 5 * * * /opt/gitlab/bin/gitlab-rake gitlab:backup:create
```

## gitlab恢复

```
# 注意备份的gitlab版本要和恢复的gitlab版本一致
# 停止数据写入服务
```

```
gitlab-ctl stop unicorn
gitlab-ctl stop sideki
gitlab-ctl status
```
# 启动恢复

```
cd /var/opt/gitlab/backups
gitlab-rake gitlab:backup:restore
```


