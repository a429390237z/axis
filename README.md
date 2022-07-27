# axis

## axis是一个通用项目
目前有定时任务服务端、定时任务客户端，之后会加入API接口及RPC接口。

## 项目构建及运行
```shell script
# 查看构建相关命令
make help
```
本地构建完镜像后可通过 `docker-compose -d` 命令，启动运行。也提供kubernetes deployment文件用于部署。
然后访问: <http://127.0.0.1:8980>查看定时任务WEBUI。
