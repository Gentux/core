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
