
# bower or npm package name
pn=

# set to sudo on FreeBSD:
sudo=

# targetDir is without / at the end!
targetDir=/usr/local/vetkakb

# --info=name1 lists only updated files
rsync=rsync -a$(n) --info=name1

.PHONY: default bower-install npm-install go-build one-time-install \
	clean-data build go-build tsc-build web

default:
	@cat GNUmakefile

build: go-build tsc-build

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
	go run vetkakb.go -c -cf vetkakb.ini

web: tsc-build

tsc-build:
	cd frontend; tsc -p tsconfig.json

$(targetDir):
	sudo mkdir $(targetDir)
	sudo chown sergei $(targetDir)

# Update targetDir from dtree on build machine
#	rsync -av $(n) dtree/ $(targetDir)/
rsync: $(targetDir)
	$(rsync) vetkakb www db-sql t-html $(targetDir)/

# NOTE: -g option on FreeBSD requires sudo
# NOTE: -g option on Mac OS X does not require (?) sudo
# NOTE: I am not convinced beta is necessary to install.
onetime-build:
	$(sudo) npm install -g typescript@beta

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

	make bower-install pn=react-router

	cd frontend; typings install react-router

	cd www; bower install less

one-time-upgrade:
	cd frontend; typings un dt~react -S
	cd frontend; typings i dt~whatwg-fetch -GS
	cd frontend; typings i dt~whatwg-streams -GS
	typings un -S react react-dom react-router
	typings i -S react react-dom react-router

	cd frontend; typings un -S react react-dom react-router
	cd frontend; typings i -S react react-dom react-router

one-time-uninstall:
	cd www; bower uninstall system.js -S
	cd www; bower uninstall react -S
	cd www; bower uninstall react-router -S
	cd www; bower uninstall less -S
	cd www; bower uninstall fetch -S

onetime:

bundle:
	cd frontend; ../node_modules/.bin/webpack

clean-data:
	rm data/*
