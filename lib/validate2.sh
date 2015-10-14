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
