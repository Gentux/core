#!/bin/bash

# Load configuration file
source configuration.sh

if [ ${#} -lt 1 ]; then
  echo "Not enough arguments"
  exit 1
fi

FILENAME=${1}
SCP=$(which scp)

sshpass -p "${PASSWORD}" "${SCP}" -P "${PORT}" "${FILENAME}" "${USER}@${SERVER}:~/Documents/"
