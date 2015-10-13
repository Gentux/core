#!/bin/bash

# ./build.sh

DESC=Uploading
DEST=esi-proxy:/home/nanosoft/0.2
_CP=scp

if [[ $1 == "local" ]]
then
  DESC=Copying
  DEST=/home/fred/Nanocloud/Dev/Go/src/nanocloud/esi/bin
  _CP=cp
fi

pwd

echo $DESC haptic
$_CP ./haptic $DEST
$_CP ./plugins/ldap/ldap $DEST/plugins/ldap/ldap
# $_CP ./plugins/owncloud/owncloud $DEST/plugins/owncloud/owncloud

$_CP -r ./public $DEST


