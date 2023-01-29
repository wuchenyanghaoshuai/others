# gitlab接入ldap
```
大概在478行开始配置

gitlab_rails['ldap_enabled'] = true
###! **remember to close this block with 'EOS' below**
gitlab_rails['ldap_servers'] = YAML.load <<-'EOS'
   main: # 'main' is the GitLab 'provider ID' of this LDAP server
     label: 'LDAP'
     host: '192.168.3.102'
     port: 389
     uid: 'cn'
     bind_dn: 'cn=admin,dc=ailieyun,dc=com'
     password: '123456'
     encryption: 'plain' # "start_tls" or "simple_tls" or "plain"
     verify_certificates: flase
     timeout: 10
     active_directory: true
     allow_username_or_email_login: false
     lowercase_usernames: false
     block_auto_created_users: false
     base: 'dc=ailieyun,dc=com'
     user_filter: ''
EOS

gitlab-ctl restart
重启以后再登陆gitlab就会发现使用原有的ldap账号就可以登陆了
```