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
# 	at least 6 characters long
# 	has at least one digit
# 	has at least one Upper case Alphabet
# 	has at least one Lower case Alphabet
function check_password {
	PASSWORD_TO_CHECK=$1
	PASSWORD_REGEX="^(?=^.{8,255}$)((?=.*\d)(?=.*[A-Z])(?=.*[a-z])|(?=.*\d)(?=.*[^A-Za-z0-9])(?=.*[a-z])|(?=.*[^A-Za-z0-9])(?=.*[A-Z])(?=.*[a-z])|(?=.*\d)(?=.*[A-Z])(?=.*[^A-Za-z0-9]))^.*$"
	TEST=`echo "$PASSWORD_TO_CHECK" | grep -oP "$PASSWORD_REGEX"` 
	
	if [[ $TEST ]]; then
		return 0
        else
		return 1
        fi
}