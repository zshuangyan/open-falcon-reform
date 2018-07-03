export HOME=/home/work
export WORKSPACE=$HOME/open-falcon
mkdir -p $WORKSPACE
DOWNLOAD="http://10.202.42.2:8089/open-falcon-v0.2.1.tar.gz"

#下载
wget $DOWNLOAD -O open-falcon-latest.tar.gz

#解压
tar -zxf open-falcon-latest.tar.gz -C $WORKSPACE

#进入工作目录
cd $WORKSPACE

#修改配置文件
sed -i "13,20s/0.0.0.0/10.202.42.2/" agent/config/cfg.json

#启动agent
./open-falcon start agent
