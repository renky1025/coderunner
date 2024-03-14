# coderunner

###  在线编译器原理
原理是每次用户提交代码就创建一个容器执行编译和运行，获取运行结果后删除容器，非常简单。

### 资源限制
程序运行时间的限制
内存的限制
cpu限制
提交代码和输入内容大小的限制

### 需要安装的docker镜像
docker pull gcc
docker pull rust
docker pull golang
docker pull python:3.10.13
docker pull khipu/openjdk17-alpine:latest
docker pull ruby
docker pull node

```shell

docker run -i ubuntu /bin/bash
cat>test.py<<\anythinghere
print("hello")
anythinghere
python test.py
```
