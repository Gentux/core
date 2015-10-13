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


PATH=$PATH:/usr/libexec
DATABASE_NAME="guacamole"
SQLDATAFILE=mysqlDump.sql


wait_for_line () {
    while read line
    do
        echo "$line" | grep -q "$1" && break
    done < "$2"

    # Read the fifo for ever otherwise process would block
    cat "$2" >/dev/null &
}


echo "Creating temporary directory to store database"
MYSQL_DATA=$(mktemp -d /tmp/haptic-mysql-XXXXX)
mkfifo ${MYSQL_DATA}/out

echo "Start MySQL process for tests"
mysqld \
    --no-defaults \
    --datadir=${MYSQL_DATA} \
    --pid-file=${MYSQL_DATA}/mysql.pid \
    --socket=${MYSQL_DATA}/mysql.socket \
    --log=${MYSQL_DATA}/log \
    --skip-networking \
    --skip-grant-tables &> ${MYSQL_DATA}/out &

echo "Waiting for MySQL to start listening to connections"
wait_for_line "mysqld: ready for connections." ${MYSQL_DATA}/out

echo "Connection String : root@unix(${MYSQL_DATA}/mysql.socket)/guacamole?loc=UTC&charset=utf8"

mysql --no-defaults -S ${MYSQL_DATA}/mysql.socket -e "CREATE DATABASE ${DATABASE_NAME};"
[ -f ${SQLDATAFILE} ] && mysql --no-defaults -S ${MYSQL_DATA}/mysql.socket ${DATABASE_NAME} < ${SQLDATAFILE}

cat << EOF > killDB.sh
#!/bin/bash

for job_pid in "$(jobs -p)"; do
    kill \${job_pid};
done
rm -rf ${MYSQL_DATA}
EOF

chmod +x killDB.sh

# Reset log file
echo "" > ${MYSQL_DATA}/log
