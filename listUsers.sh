#!/bin/bash

# Configuration
MYPATH=/home/nanosoft/0.2

##################################
# Proxy Users                    #
##################################

#echo -e "$(eval "echo -e \"`<$MYPATH/sql/users_list.sql`\"")" | mysql -B -u root --password='pass' --database='guacamole' | $MYPATH/mysql_to_json.sh

# Run database script to list users 
# and format timestamp with ISO 8601
echo -e "$(eval "echo -e \"`<$MYPATH/sql/users_list.sql`\"")" | mysql -B -u root --password='pass' --database='guacamole' | awk  -F $'\t' '
																BEGIN { OFS = FS }
																{
																	if (NR==1) {
																		print
																	} else {
																		str = "\""$3" "$4"\""
																		cmd = "date -d " str " --iso-8601=second"
																		while (cmd | getline line) {
																			print $1 "\t" $2 "\t" line
																		}
																		close(cmd)
																	}
																}' | $MYPATH/mysql_to_json.sh

