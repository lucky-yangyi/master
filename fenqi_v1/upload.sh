#! /bin/bash
echo '------------------pack------------------'

tar -zcvf fenqi_v1_release.tar.gz cache controllers routers \
models services utils views static main.go
echo '-----------------upload-----------------'

scp fenqi_v1_release.tar.gz yg@10.139.101.57:/home/yg/go/src/fenqi_v1
scp fenqi_v1_release.tar.gz yg@10.139.100.117:/home/yg/go/src/fenqi_v1
echo '-------------------server-126-restart-------------------'

ssh yg@10.139.101.57  << remotessh
cd /home/yg/go/src/fenqi_v1
tar -zxvf fenqi_v1_release.tar.gz
go build -o fenqi_v1 main.go
exit
remotessh
echo '-------------------server-170-restart-------------------'
ssh yg@10.139.100.117  << remotessh
cd /home/yg/go/src/fenqi_v1
tar -zxvf fenqi_v1_release.tar.gz
go build -o fenqi_v1 main.go
exit
remotessh
echo '------------------delete------------------'
rm -f fenqi_v1_release.tar.gz
echo '------------------done------------------'
