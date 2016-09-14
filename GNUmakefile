
# bower or npm package name
pn=

# set to sudo on FreeBSD:
sudo=

.PHONY: default bower-install npm-install go-build one-time-install \
	clean-data

default:
	@cat GNUmakefile

# install new bower pcakge with command
# 	make bower-install bpn=react
# Installs packages listed in bower.json run without arguments:
#	make bower-install
# To initialize new bower.json run:
#	bower init
# For additional info see:
#	http://frontendbabel.info/articles/bower-why-frontend-package-manager/
# To find new packages:
#	bower search
bower-install:
	cd www; bower install --save $(pn)

# To initialize npm package system use
# 	npm init
# Installed packages:
# 	babel-preset-es2015 babel-preset-react
npm-install:
	npm install --save-dev $(pn)

# Main go build:
go-build:
	go build

run:
	go run vetkakb.go

tsc-build:
	cd frontend; tsc -p tsconfig.json

# NOTE: -g option on FreeBSD requires sudo
# NOTE: -g option on Mac OS X does not require (?) sudo
one-time-install:
	$(sudo) npm install -g typescript@beta
	$(sudo) npm install typings --global

	# the following line will fail on fresh build box
	cd frontend; typings init

	cd www; bower install system.js --save

	# fetch project installed here
	# has a second name whatwg-fetch
	cd www; bower install fetch --save
	cd frontend; $(sudo) typings install dt~whatwg-fetch --global --save

	cd frontend; $(sudo) typings install dt~react --global --save
	cd frontend; $(sudo) typings install dt~react-dom --global --save

	cd frontend; typings install es6-promise
	cd frontend; typings uninstall es6-promise

onetime:

clean-data:
	rm data/*
