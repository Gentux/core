#!/bin/bash

ZEROINSTAL_DIR="/home/gentux/projects/gocode/src/nanocloud.com/zeroinstall"
HAPTIC_SRC_DIR="${ZEROINSTAL_DIR}/agent/haptic"
PLUGINS_SRC_DIR="${ZEROINSTAL_DIR}/plugins"
BIN_DIR="${ZEROINSTAL_DIR}/bin/haptic"

echo "###### Set NANOCONF"
export NANOCONF="${HAPTIC_SRC_DIR}/config-local.json"

echo "###### Killing remaining plugins"
JOBS_PID=$(ps aux | grep "pingo" | awk '/:proto=/ { print $2; }')
if [ ! "" = "${JOBS_PID}" ]; then
    for job_pid in ${JOBS_PID}; do
        kill "${job_pid}"
    done
fi

echo "###### Removing previous binaries"
rm ${BIN_DIR}/plugins/iaas/iaas
rm ${BIN_DIR}/plugins/ldap/ldap
rm ${BIN_DIR}/plugins/owncloud/owncloud
rm ${BIN_DIR}/haptic

echo "###### Building IaaS Plugin"
( cd ${PLUGINS_SRC_DIR}/iaas ; go build -a -v -o ${BIN_DIR}/plugins/iaas/iaas)

echo "###### Building  Plugin"
( cd ${PLUGINS_SRC_DIR}/ldap ; go build -a -v -o ${BIN_DIR}/plugins/ldap/ldap)

echo "###### Building IaaS Plugin"
( cd ${PLUGINS_SRC_DIR}/owncloud ; go build -a -v -o ${BIN_DIR}/plugins/owncloud/owncloud)

echo "###### Rebuild Haptic"
( cd ${HAPTIC_SRC_DIR} ; go build -a -v -o ${BIN_DIR}/haptic)

echo "###### Launching Haptic"
( cd ${BIN_DIR} ; ./haptic serve )
