#!/bin/bash

scriptPath=`readlink -f $0̀`
dirPath=`dirname $scriptPath`

DESC=Uploading
DEST=esi-proxy:/home/nanosoft/0.2
_CP="scp -q"

if [[ $1 == "local" ]]
then
  DESC=Copying
  DEST=/home/fred/Nanocloud/Dev/Go/src/nanocloud/esi/bin
  _CP=cp
fi

#echo $DESC $dirPath/haptic
 
$_CP ${dirPath}/haptic $DEST
$_CP ${dirPath}/plugins/ldap/ldap $DEST/plugins/ldap/ldap
$_CP ${dirPath}/plugins/owncloud/owncloud $DEST/plugins/owncloud/owncloud

# $_CP ./plugins/owncloud/owncloud $DEST/plugins/owncloud/owncloud

# $_CP -r $dirPath/public $DEST
