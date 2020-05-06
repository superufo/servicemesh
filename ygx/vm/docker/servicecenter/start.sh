#!/bin/sh

umask 027

cd /opt/service-center

sed -i "s|^registry_plugin.*=.*$|registry_plugin = etcd|g" conf/app.conf
sed -i "s|^manager_cluster.*=.*$|manager_cluster = ${BACKEND_ADDRESS}|g"  conf/app.conf
sed -i "s/^httpaddr.*=.*$/httpaddr = ${IP_ADDRESS}/g" conf/app.conf

./service-center