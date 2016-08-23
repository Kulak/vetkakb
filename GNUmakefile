
# bower or npm package name
pn=

.PHONY: default bower-install npm-install go-build one-time-install

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

one-time-install:
	npm install typescript -g
	npm install typings --global

	cd frontend; typings init
	cd frontend; typings install react
	cd frontend; typings install react-dom

	cd www; bower install system.js --save
