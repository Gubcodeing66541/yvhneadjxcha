#### 安装前
开放6379  3306端口


#### 安装docker
sudo yum install -y yum-utils  device-mapper-persistent-data lvm2
  
yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
      
yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

systemctl start docker

#### 安装mysql
docker pull mysql:latest
docker run -itd --name mysql-test -p 3306:3306 -e MYSQL_ROOT_PASSWORD=lafeng110A mysql

#### 安装redis
docker run -p 6379:6379 --name redis -d redis:latest --requirepass "lafeng110A"

#### docker 进入mysql 创建数据库`
docker exec -it mysql-test /bin/bash
mysql -h 127.0.0.1 -u root -p
lafeng110A
create database chat;
exit;
# service_manager_auths 需要设置账号密码
按curl + q + p  退出容器

#### 修改config配置文件的IP
vi config.json
修改head ip

#### 启动服务
chmod -R 777 server & nohup ./server &

#### 日志命令行安装
go get github.com/go-swagger/go-swagger
cd Gopath/pkg/mod/github.com/go-swagger
go install github.com/swaggo/swag/cmd/swag
swagger version


#### 注释日志生成
swag init -g ./server.go

#### 打包提交代码
set goos=linux
go build server.go
git add .
git commit -m 'u'
git push

#### 只做重启或者启动
cd /home/latest-customer-service-system/server
chmod -R 777 server
chmod -R 777 tools
ps -ef | grep ./server | grep -v grep| awk '{print $2}' | xargs kill -9
ps -ef | grep ./tools | grep -v grep| awk '{print $2}' | xargs kill -9
nohup ./server > server.log 2>&1 &
nohup ./tools > tools.log 2>&1 &

#### 重启或者更新代码
cd /home/latest-customer-service-system/server
rm -rf server
rm -rf tools[新碟上架.url](..%2F..%2F..%2FAppData%2FLocal%2FTemp%2F%D0%C2%B5%FA%C9%CF%BC%DC.url)
git reset --hard HEAD && git pull
chmod -R 777 server
chmod -R 777 tools

#### 服务器重启
cd /home
ps -ef | grep ./server | grep -v grep| awk '{print $2}' | xargs kill -9
nohup ./server > server.log 2>&1 &

#### 定时任务重启 注意只能主服务器启动 其他服务区不能启动
ps -ef | grep ./tools | grep -v grep| awk '{print $2}' | xargs kill -9
nohup ./tools > tools.log 2>&1 &

#### 查看请求日志
tail -fn50 server.log


#### 整个服务器重启处理
systemctl start docker
docker start redis
ps -ef | grep ./tools | grep -v grep| awk '{print $2}' | xargs kill -9
nohup ./tools > tools.log 2>&1 &
## 以上四个命令只能在主服务器执行（前两句是重启docker和redis,后两句是重启定时任务）
cd /home
ps -ef | grep ./server | grep -v grep| awk '{print $2}' | xargs kill -9
nohup ./server > server.log 2>&1 &
## 以上三个命令是重启服务器（主服务器和副服务器都要执行）

#https://gatewaynew.azure-api.net/user/code/actions?code=cbbeb6ec2a5e0c041

#### 安装按照模板的创建【需要注意部署规范】
# 安装Certbot和Nginx（如果尚未安装）
sudo apt-get update
sudo apt-get install software-properties-common
sudo add-apt-repository universe
sudo add-apt-repository ppa:certbot/certbot
sudo apt-get update
sudo apt-get install certbot python3-certbot-nginx nginx

# 运行Certbot以自动获取证书
sudo certbot --nginx

# 如果你想手动指定域名，可以使用以下命令
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

sudo certbot --nginx -d fnukywkw.md.ci -d fnukywkw.md.ci

certbot certonly --email admin@www.xxx.cn  --webroot -w ./ -d fnukywkw.md.ci
certbot certonly --email admin@www.xxx.cn  --webroot -w ./ -d fnukywkw.md.ci

echo /etc/letsencrypt/live/fnukywkw.md.ci/fullchain.pem >> Test/build/
/etc/letsencrypt/live/fnukywkw.md.ci/privkey.pem
ps -ef | grep nginx | grep -v grep| awk '{print $2}' | xargs kill -9
