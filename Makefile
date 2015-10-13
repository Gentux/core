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

# SHELL = /bin/bash

.DEFAULT_GOAL := haptic

ldap: ../bin/haptic/plugins/ldap/ldap

#.PHONY: ../bin/haptic/plugins/ldap/ldap
../bin/haptic/plugins/ldap/ldap: ../plugins/ldap/main.go
	go build -o ../bin/haptic/plugins/ldap/ldap nanocloud.com/plugins/ldap

iaas: ../bin/haptic/plugins/iaas/iaas

#.PHONY: ../bin/haptic/plugins/iaas
../bin/haptic/plugins/iaas/iaas: ../plugins/iaas/main.go
	go build -o ../bin/haptic/plugins/iaas/iaas nanocloud.com/plugins/iaas

owncloud: ../bin/haptic/plugins/owncloud/owncloud

#.PHONY: ../bin/haptic/plugins/owncloud/owncloud
../bin/haptic/plugins/owncloud/owncloud: ../plugins/owncloud/main.go
	go build -o ../bin/haptic/plugins/owncloud/owncloud nanocloud.com/plugins/owncloud

haptic: iaas ldap owncloud ../bin/haptic/haptic

.PHONY: ../bin/haptic/haptic
../bin/haptic/haptic:
	go build -o ../bin/haptic/haptic nanocloud.com/core/haptic

setup:
	mkdir -p ../bin/haptic/plugins
	mkdir -p ../bin/haptic/plugins/iaas
	mkdir -p ../bin/haptic/plugins/ldap
	mkdir -p ../bin/haptic/plugins/owncloud

	mkdir -p ../bin/haptic/external/bin

	echo "Installing go packages dependencies"
	go get github.com/dullgiulio/pingo
	go get golang.org/x/net/icmp
	go get golang.org/x/net/internal/iana
	go get golang.org/x/net/ipv4
	go get github.com/boltdb/bolt
	go get github.com/gorilla/rpc
	go get github.com/gorilla/rpc/json
	go get github.com/gorilla/securecookie
	go get github.com/hypersleep/easyssh
	go get github.com/go-sql-driver/mysql

        # Copy configuration files

	@ if [ ! -f ../bin/haptic/plugins/iaas/config.json ]; then \
		echo "One time creation of config file: ../bin/haptic/plugins/iaas/config.json" ; \
		cp ../plugins/iaas/config.json.sample ../bin/haptic/plugins/iaas/config.json; \
	fi

	@ if [ ! -f ../bin/haptic/plugins/ldap/config.json ]; then \
		echo "One time creation of config file: ../bin/haptic/plugins/ldap/config.json" ; \
		cp ../plugins/ldap/config.json.sample ../bin/haptic/plugins/ldap/config.json; \
	fi

	cp ../plugins/ldap/*.php ../bin/haptic/plugins/ldap/;

	@ if [ ! -f ../bin/haptic/plugins/owncloud/config.json ]; then \
		echo "One time creation of config file: ../bin/haptic/plugins/owncloud/config.json" ; \
		cp ../plugins/owncloud/config.json.sample ../bin/haptic/plugins/owncloud/config.json; \
	fi

	@ if [ ! -f ../bin/haptic/config.json ]; then \
		echo "One time creation of config file: ../bin/haptic/config.json" ; \
		cp ./haptic/config.json.sample ../bin/haptic/config.json; \
	fi

	# Copy binaries
	cp ./external/bin/winexe-static ../bin/haptic/external/bin/

serve:
	../bin/haptic/haptic serve

clean:
	rm ../bin/haptic/plugins/iaas/iaas
	rm ../bin/haptic/plugins/ldap/ldap
	rm ../bin/haptic/plugins/owncloud/owncloud
	rm ../bin/haptic/haptic
	go clean nanocloud.com/plugins/iaas
	go clean nanocloud.com/plugins/owncloud
	go clean nanocloud.com/plugins/ldap
	go clean nanocloud.com/core/haptic
