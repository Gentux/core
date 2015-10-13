# SHELL = /bin/bash

.DEFAULT_GOAL := haptic

owncloud: ./bin/haptic/plugins/owncloud/owncloud

./bin/haptic/plugins/owncloud/owncloud: ./plugins/owncloud/main.go
	go build -o ./bin/haptic/plugins/owncloud/owncloud nanocloud.com/zeroinstall/plugins/owncloud

ldap: ./bin/haptic/plugins/ldap/ldap

./bin/haptic/plugins/ldap/ldap: ./plugins/ldap/main.go
	go build -o ./bin/haptic/plugins/ldap/ldap nanocloud.com/zeroinstall/plugins/ldap

haptic: ldap owncloud ./bin/haptic/haptic

./bin/haptic/haptic:
	go build -o ./bin/haptic/haptic nanocloud.com/zeroinstall/agent/haptic

dev: haptic
	mkdir -p ./bin/haptic/plugins
	mkdir -p ./bin/haptic/plugins/ldap
	mkdir -p ./bin/haptic/plugins/owncloud

	#cp ./plugins/ldap/ldap ./agent/haptic/plugins/ldap/
	@ if [ ! -f ./bin/haptic/plugins/ldap/config.json ]; then \
		echo "One time creation of config file: .bin/haptic/plugins/ldap/config.json" ; \
		cp ./plugins/ldap/config.json.sample ./bin/haptic/plugins/ldap/config.json; \
	fi

	# cp ./plugins/owncloud/owncloud ./agent/haptic/plugins/owncloud/
	@ if [ ! -f ./bin/haptic/plugins/owncloud/config.json ]; then \
		echo "One time creation of config file: ./bin/haptic/plugins/owncloud/config.json" ; \
		cp ./plugins/owncloud/config.json.sample ./bin/haptic/plugins/owncloud/config.json; \
	fi

#echo "Copying config.json to: ./agent/haptic/plugins/owncloud/config.json"; \

clean:
	rm ./bin/haptic/haptic
	rm ./bin/haptic/plugins/ldap/ldap
	rm ./bin/haptic/plugins/owncloud/owncloud
	go clean nanocloud.com/zeroinstall/plugins/owncloud
	go clean nanocloud.com/zeroinstall/plugins/ldap
	go clean nanocloud.com/zeroinstall/agent/haptic
