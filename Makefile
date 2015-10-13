# SHELL = /bin/bash

owncloud: ./plugins/owncloud/owncloud

./plugins/owncloud/owncloud: ./plugins/owncloud/main.go
	go build -o ./plugins/owncloud/owncloud nanocloud.com/zeroinstall/plugins/owncloud

ldap: ./plugins/ldap/ldap

./plugins/ldap/ldap: ./plugins/ldap/main.go
	go build -o ./plugins/ldap/ldap nanocloud.com/zeroinstall/plugins/ldap

haptic: ldap owncloud ./agent/haptic/haptic

./agent/haptic/haptic:
	go build -o ./agent/haptic/haptic nanocloud.com/zeroinstall/agent/haptic

dev: haptic
	mkdir -p ./agent/haptic/plugins
	mkdir -p ./agent/haptic/plugins/ldap
	mkdir -p ./agent/haptic/plugins/owncloud

	cp ./plugins/ldap/ldap ./agent/haptic/plugins/ldap/
	@ if [ ! -f ./agent/haptic/plugins/ldap/config.json ]; then \
		echo "One time creation of config file: agent/haptic/plugins/ldap/config.json" ; \
		cp ./plugins/ldap/config.json ./agent/haptic/plugins/ldap/config.json; \
	fi

	cp ./plugins/owncloud/owncloud ./agent/haptic/plugins/owncloud/
	@ if [ ! -f ./agent/haptic/plugins/owncloud/config.json ]; then \
		echo "One time creation of config file: agent/haptic/plugins/owncloud/config.json" ; \
		cp ./plugins/owncloud/config.json ./agent/haptic/plugins/owncloud/config.json; \
	fi

#echo "Copying config.json to: ./agent/haptic/plugins/owncloud/config.json"; \

clean:
	go clean nanocloud.com/zeroinstall/plugins/owncloud
	go clean nanocloud.com/zeroinstall/plugins/ldap
	go clean nanocloud.com/zeroinstall/agent/haptic
