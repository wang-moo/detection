# 使用
> 此程序转为探测 tcp 业务，和udp 存在，同一个host 只要有一个返回信息错误，同host都会下线切换到其余host
 - 程序启动时默认读取配置路径 ./config
 - 日志 保存 /tmp/log/detection.log 
    > 最大50M  保存7天，最多一个备份文件
## 配置说明
 - interval: 10
    >探测间隔（秒）

 - fail_threshold: 3
    >阈值设置 探测失败的次数
 - success_threshold: 2
    >阈值设置 探测成功的次数

# UDP 服务健康检查配置
> 原理构建udp 的ping包 接收返回的pong包，只有收到pong 才算成功
> 请求超时为1秒
```
udp:
    - host: 1.1.1.1
      port: 12345
    - host: 1.1.1.2
      port: 12345
```

# TCP 服务健康检查配置
> 检擦服务内部的提供的健康 
> GET 请求目前只支持 返回 {"status": true}
> 请求超时为2秒
```
tcp:
    - host: 1.1.1.1
      port: 80
      url: /health/
    - host: 1.1.1.2
      port: 80
      url: /health/
```

# 模板文件
> 此配置是指定 config 下面的模板文件名，模板内部只做server 的动态修改其余都必须写死
```
conf_file:
    - 12345
    - 80
```

# Nginx upstream 配置文件路径
```
nginx:
    command: nginx -s reload # 此配置是服务上线下线需要执行的命令
    conf_path: /etc/nginx/conf.d/xxx/   # 此路径是根据模板生成的文件保存目录，文件夹下是  {conf_file}.conf  文件
```

# 注
- nginx 根配置内部需要 手动调整每个配置的引用位置 /etc/nginx/conf.d/xxx/xxxx.conf