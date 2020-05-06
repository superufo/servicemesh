#!/usr/bin/env bash
set -e

umask 027

cd /opt/service-center

set +e

#sed -i "s/^httpaddr.*=.*$/httpaddr = $(hostname)/g" conf/app.conf
if [ ! -z "${BACKEND_ADDRESS}" ]; then
    sed -i "s|^registry_plugin.*=.*$|registry_plugin = etcd|g" conf/app.conf
    sed -i "s|^manager_cluster.*=.*$|manager_cluster = ${BACKEND_ADDRESS}|g" conf/app.conf
fi

sed -i "s/^httpaddr.*=.*$/httpaddr = 172.20.0.103/g" conf/app.conf

set -e

./service-center


######## 1.3.0
#umask 027
#
#cd /opt/service-center
#
#sed -i "s|^registry_plugin.*=.*$|registry_plugin = etcd|g" conf/app.conf
#sed -i "s|^manager_cluster.*=.*$|manager_cluster = 172.20.0.100:2379,172.20.0.101:2379,172.20.0.102:2379|g"  conf/app.conf
#sed -i "s/^httpaddr.*=.*$/httpaddr = 172.20.0.103/g" conf/app.conf
#
#./service-center

######## 1.3.0.1
#umask 027
#
#cd /opt/service-center
#
#sed -i "s|^registry_plugin.*=.*$|registry_plugin = etcd|g" conf/app.conf
#sed -i "s|^manager_cluster.*=.*$|manager_cluster = ${BACKEND_ADDRESS}|g"  conf/app.conf
#sed -i "s/^httpaddr.*=.*$/httpaddr = ${IP_ADDRESS}/g" conf/app.conf
#
#./service-center









