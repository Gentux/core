#!/bin/bash

# Configuration
MYPATH=/home/nanosoft/0.2

##################################
# Proxy Users                    #
##################################

# Run database script to list users 
echo -e "$(eval "echo -e \"`<$MYPATH/sql/users_list_mini.sql`\"")" | mysql -B -u root --password='pass' --database='guacamole' | $MYPATH/mysql_to_json.sh
