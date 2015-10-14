#!/bin/bash

# Configuration
MYPATH=/home/nanosoft/0.2
WINEXE=$MYPATH'/winexe-static -U 'intra.nanocloud.com/Administrator%password' --runas='intra.nanocloud.com/Administrator%password' //10.20.12.20'

# Add logs functions
. $MYPATH/lib/log.sh

# Get command line parameters
EMAIL=$1
PASSWORD=$2

# Workflow variables
CODE=0

# Workflow messages
MESSAGE[0]="Something went wrong"
MESSAGE[1]="The password has been changed"
MESSAGE[2]="User does not exist"
MESSAGE[3]="The password does not respect the security policy"
MESSAGE[4]="Problem with email format"
MESSAGE[5]="Problem with TAC server"

# Add validation functions
. $MYPATH/lib/validate2.sh

# Check parameters
if check_email $EMAIL ; then true ; else script_exit "ERROR" 4 ; fi
if check_password "$PASSWORD" ; then true ; else script_exit "ERROR" 3 ; fi
if $MYPATH/listUsers_mini.sh | grep \"$EMAIL\" >/dev/null 2>&1 ; then true ; else script_exit "ERROR" 2 ; fi

################################
# Windows Server Configuration #
################################

if [ $CODE -le 1 ]
then
	mkdir -p $MYPATH/studio/$EMAIL
        chmod o+w $MYPATH/studio/$EMAIL

        # Retrieve SAM user with the email adress
        SAM=$(echo -e "$(eval "echo -e \"`<$MYPATH/sql/sam.sql`\"")" | mysql -u root --password='pass' --database='guacamole' --skip-column-names)

        log_it "INFO" "Logoff user"
        cp $MYPATH/studio/logoffUser.cmd $MYPATH/studio/$EMAIL/$SAM.logoffUser.bat
        sed -i s/SAMUSER/"$SAM"/g $MYPATH/studio/$EMAIL/$SAM.logoffUser.bat

        unix2dos $MYPATH/studio/$EMAIL/$SAM.logoffUser.bat
        scp $MYPATH/studio/$EMAIL/$SAM.logoffUser.bat Administrator@10.20.12.20:/cygdrive/c/Windows/SYSVOL/domain/scripts/ > /dev/null 2>&1

        # Run the logoff script
        eval "$WINEXE 'cmd.exe /C \\\winad.intra.nxbay.com\NETLOGON\'\$SAM'.logoffUser.bat'" > /dev/null 2>&1

        # Delete the logoffUser File on Windows Server
        eval "$WINEXE 'cmd.exe /C DEL \\\winad.intra.nxbay.com\NETLOGON\'\$SAM'.logoffUser.bat'" > /dev/null 2>&1

        # Setup the Workspace file
        cp $MYPATH/studio/workspace.cmd $MYPATH/studio/$EMAIL/$SAM.config.bat
        sed -i s/USEREMAIL/"$EMAIL"/g $MYPATH/studio/$EMAIL/$SAM.config.bat
        sed -i s/USERPASSWORD/"$PASSWORD"/g $MYPATH/studio/$EMAIL/$SAM.config.bat
        sed -i "s|USERFQDN|$TAC_URL|g" $MYPATH/studio/$EMAIL/$SAM.config.bat

        unix2dos $MYPATH/studio/$EMAIL/$SAM.config.bat

        # Put the Workspace file on the server
        scp $MYPATH/studio/$EMAIL/$SAM.config.bat Administrator@10.20.12.20:/cygdrive/c/Windows/SYSVOL/domain/scripts/ > /dev/null 2>&1

        # Setup the Windows Server Security file
        cp $MYPATH/studio/setSecurity.cmd $MYPATH/studio/$EMAIL/$SAM.setSecurity.bat
        sed -i s/SAMUSER/"$SAM"/g $MYPATH/studio/$EMAIL/$SAM.setSecurity.bat

        unix2dos $MYPATH/studio/$EMAIL/$SAM.setSecurity.bat

        scp $MYPATH/studio/$EMAIL/$SAM.setSecurity.bat Administrator@10.20.12.20:/cygdrive/c/Windows/SYSVOL/domain/scripts/ > /dev/null 2>&1

        # Execute Security File on Windows Server
        eval "$WINEXE 'cmd.exe /C \\\winad.intra.nxbay.com\NETLOGON\'\$SAM'.setSecurity.bat'" > /dev/null 2>&1

        # Delete the Security File on Windows Server
        eval "$WINEXE 'cmd.exe /C DEL \\\winad.intra.nxbay.com\NETLOGON\'\$SAM'.setSecurity.bat'" > /dev/null 2>&1

fi


##################################
# Proxy User                     #
##################################

if [ $CODE -le 1 ]

then

        # Run database script to change user password
        echo -e "$(eval "echo -e \"`<$MYPATH/sql/change_user_password.sql`\"")" | mysql -u root --password='pass' --database='guacamole'

fi

# UUID=$(echo -e "$(eval "echo -e \"`<$MYPATH/sql/tac_get.sql`\"")" | mysql -N -u root --password='pass' --database='guacamole')

CODE=1

# Final code
script_exit "INFO" $CODE
