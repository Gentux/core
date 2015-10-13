#!/bin/bash

# Nanocloud community -- transform any application into SaaS solution
#
# Copyright (C) 2015 Nanocloud Software
#
# This file is part of Nanocloud community.
#
# Nanocloud community is free software; you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# Nanocloud community is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.


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
