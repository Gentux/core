#!/bin/bash


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
