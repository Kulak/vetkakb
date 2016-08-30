
# bower or npm package name
pn=

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

one-time-install:
	npm install typescript -g
	npm install typings --global

	cd frontend; typings init

	# cd frontend; typings install react
	# cd frontend; typings uninstall react --save
	# cd frontend; typings install react-dom
	# cd frontend; typings uninstall react-dom --save

	cd www; bower install system.js --save

	# fetch project installed here
	# has a second name whatwg-fetch
	cd www; bower install fetch --save
	cd frontend; typings install dt~whatwg-fetch --global --save

	cd frontend; typings install dt~react --global --save
	cd frontend; typings install dt~react-dom --global --save

onetime:

clean-data:
	rm data/*
