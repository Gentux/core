# Exit function with errors
# $1: Exit type (SUCCESS, ERROR, DEBUG, ...)
# $2: Exit message
function script_exit {
        printf -v result '{"code":"%s","message":"%s"}' "$2" "${MESSAGE[$2]}"
        log_it $1 "Output result --> $result"
        echo $result
        exit
}

# Log function
# $1: Log type (INFO, ERROR, DEBUG, ...)
# $2: Log message
function log_it {
        echo -e "$(date '+%F %T %Z') - $1: $SCRIPTNAME $PARAMETERS : $2" >> $LOGFILE
}

# Configuration
LOGFILE=$MYPATH/log/users.log
SCRIPTNAME="$(basename "$(test -L "$0" && readlink "$0" || echo "$0")")"
PARAMETERS="${@}"
