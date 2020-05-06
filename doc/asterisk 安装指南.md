asterisk 安装指南

```shell
install https://wiki.asterisk.org/wiki/display/AST/Installing+Asterisk+From+Source



wget   http://downloads.asterisk.org/pub/telephony/asterisk/releases/asterisk-17.1.0.tar.gz
wget  http://downloads.asterisk.org/pub/telephony/asterisk/releases/asterisk-16.6.0.tar.gz


wget  https://downloads.asterisk.org/pub/telephony/libpri/libpri-current.tar.gz

wget  https://downloads.asterisk.org/pub/telephony/dahdi-linux/dahdi-linux-current.tar.gz &&

wget https://downloads.asterisk.org/pub/telephony/dahdi-tools/dahdi-tools-current.tar.gz &&

wget https://downloads.asterisk.org/pub/telephony/dahdi-linux-complete/dahdi-linux-complete-current.tar.gz



 wget  http://repo.iotti.biz/CentOS/7/x86_64/kernel-devel-3.10.0-957.21.3.el7.centos.plus.x86_64.rpm

 rpm -Uvh  kernel-devel-3.10.0-957.21.3.el7.centos.plus.x86_64.rpm

 yum install kernel-devel-$(uname -r)

重启机器 reboot 

 tar zxvf  dahdi-linux-complete-current.tar.gz  &&  cd  dahdi-linux-complete-3.1.0+3.1.0  &&  make && make install &&  make install-config && cd ../



 tar zxvf libpri-current.tar.gz && cd  libpri-1.6.0 && make && make install  && cd ../
```



```shell
yum install -y bzip2
wget http://www.pjsip.org/release/2.6/pjproject-2.6.tar.bz2
wget https://www.pjsip.org/release/2.9/pjproject-2.9.tar.bz2
tar -xjf  pjproject-2.6.tar.bz2  
tar -xjvf pjproject-2.6.tar.bz2 && cd pjproject-2.6 && ./configure --prefix=/usr --enable-shared --disable-sound --disable-resample --disable-video --disable-opencore-amr CFLAGS='-O2 -DNDEBUG'
&& make dep &&  make &&  make install && ldconfig  && cd ../

ldconfig -p | grep pj
asterisk 中 make menuselect 的时候存在 res_pjsip 
```



```
http://downloads.digium.com/pub/telephony/codec_opus/asterisk-17.0/x86-64/codec_opus-17.0_1.3.0-x86_64.tar.gz

tar zxvf codec_opus-17.0_1.3.0-x86_64.tar.gz 
cp /download/codec_opus-17.0_1.3.0-x86_64/*  /usr/lib/
cp /download/codec_opus-17.0_1.3.0-x86_64/*  /usr/lib64/
cp /download/codec_opus-17.0_1.3.0-x86_64/*  /usr/libexec/
```





```
tar zxvf asterisk-17.1.0.tar.gz
cd  asterisk-17.1.0
yum install patch  libedit  libedit-devel.i686  libedit-devel.x86_64  libedit.i686 libedit.x86_64  mpg123

#https://wiki.asterisk.org/wiki/display/AST/Installing+the+Asterisk+Test+Suite#InstallingtheAsteriskTestSuite-pjsua_installationPJSUAInstallation
## 会自动下载pjproject-2.9.tar.bz2编译

#清理缓存
make distclean

#检测需要安装哪些依赖项
contrib/scripts/install_prereq test

yum install mpg123  (播放mp3文件)
yum install --skip-broken --assumeyes speexdsp-devel libogg-devel libvorbis-devel alsa-lib-devel portaudio-devel xmlstarlet bison flex postgresql-devel unixODBC-devel gmime-devel lua-devel uriparser-devel libxslt-devel mysql-devel bluez-libs-devel radcli-devel freetds-devel jack-audio-connection-kit-devel iksemel-devel corosynclib-devel libical-devel spandsp-devel libresample-devel uw-imap-devel libsrtp-devel graphviz openldap-devel hoard codec2-devel fftw-devel libsndfile-devel unbound-devel python-devel

#自动安装依赖项
contrib/scripts/install_prereq install

./configure   --with-pjproject-bundled   --with-jansson-bundled   

make menuselect  
(webrtc 必须要的模块 res_crypto res_http_websocket res_pjsip_transport_websocket codec_opus )

make install
#安装配置文件例子
make samples
#安装服务脚本，这样程序可以作为系统服务自动启动
make config
#安装日志管理工具
make install-logrotate

make && make install && make samples && make config && make install-logrotate

service asterisk start
asterisk -rvvvvvvvvvv
# 加载 sip 通道
module load chan_sip.so
# 
module reload chan_sip.so
```



```
#设置防火墙，开放5060端口
firewall-cmd --permanent --add-port 5060/udp
firewall-cmd --permanent --add-port 5060/tcp
firewall-cmd --reload
firewall-cmd --list-port

iptables -A INPUT -p tcp -m tcp --dport 3306 -j ACCEP
iptables -A INPUT -p tcp -m tcp --dport 3306 -j ACCEPT
```

参考:http://www.jiazi.cn/blog/?id=63