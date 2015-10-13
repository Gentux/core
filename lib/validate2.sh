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

# Check if this is an email address
function check_email {
	EMAIL_TO_CHECK=$1
	valid_email_regex="^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,4}$"
	
	if [[ "$EMAIL_TO_CHECK" =~ $valid_email_regex ]]; then
		return 0
	else
		return 1
	fi
}

# Check if the password meets the following requirements:
# 	at least 7 and less than 65 characters long
# 	has at least one digit
# 	has at least one Upper case Alphabet
# 	has at least one Lower case Alphabet
#       characters that can be used:
#           any alphanumeric character 0 to 9 OR A to Z or a to z
#           punctuation symbols . , " ' ? ! ; : # $ % & ( ) * + - / < > = @ [ ] \ ^ _ { } | ~ 
function check_password {
	s=$1

	if [[ ${#s} -ge 7 && ${#s} -le 64 && "$s" == *[[:upper:]]* && "$s" == *[[:lower:]]* && "$s" == *[[:digit:]]* && "$s" =~ ^[[:alnum:][:punct:]]+$ ]]; then
		return 0
        else
		return 1
        fi
}

# Check if the name meets the following requirements:
#       at least 1 and less than 65 characters long
#       characters that can be used:
#           any alphanumeric character 0 to 9 OR A to Z or a to z
#           punctuation symbols . , " ' ? ! ; : # $ % & ( ) * + - / < > = @ [ ] \ ^ _ { } | ~ 
function check_name {
        s=$1

        if [[ ${#s} -ge 1 && ${#s} -le 64 && "$s" =~ ^[[:alnum:][:punct:]]+$ ]]; then
                return 0
        else
                return 1
        fi
}

function validate_url {
  	if [[ `wget -S --spider $1  2>&1 | grep 'HTTP/1.1 200 OK'` ]]; then 
                return 1
        else
                return 0
        fi
}
