# Vetka KB Project Track

## Design

### Datastore Design

### UI Design

Search List
	- "Most Recent" entries button.
	- "Search Field" to just type to search.
	- "Entrie List" displays search results.

### Functional Approaches


### External vs Internal Blob

Embedded blob storage is more optimal for data that's less than 20K.

File system storage is more optimal for 50K file sizes and above.  At 100K and above it is guarnateed to be faster with file system.  Source: <https://www.sqlite.org/intern-v-extern-blob.html>

### HTML Tokenization

The simplest way to avoid writing custom HTML tokenizer is to convert HTML to text with lynx text browser:

	curl -s http://www.sqlite.org | lynx -nolist -stdin -dump

Source: <https://www.mail-archive.com/sqlite-users@mailinglists.sqlite.org/msg82126.html>

Other options are available.

## Log

### Fetch API Support

Add fetch API at browser level with:

	cd frontend; typings install dt~whatwg-fetch --global --save

command.  However, this command depends on Promise and ES6 functionality.

Installing promise typings did not help.  The solution was to alter compiler configuration in `tsconfig.json` file:

	"lib": ["es6", "dom"],

Option "lib" requires TypeScript 2.0 to work.

This magically enabled ES6 Promise support and compilation error was gone.  Once "lib" is used other options have to be specified explicitly.  "dom" option is required for "react" dependencies to work properly.

If typing is installed with option --global, then it does not need to be imported (it is imported automatically).

The systemjs still needs to map polyfil name 'whatwg-fetch' to the actual library.  So, a change in HTML file was necessary.

Because "fetch" is global in typescript I dont' use "import".  Because there is no import in source code, "fetch" polyfill does not get imported and we have to import it directly in `index.html` file.

### Seach Typings

	typings search aName

Note that source column lists name of the source.  NPM is default (?).  If source is DT, then it has to be specified in install command.  See install `dt~promise` below.

### Correct way to install typings

These commands were rolled back, but they represent good example of how to work with typings:

	# installes into typings/modules/es6-promise
	cd frontend; typings install es6-promise --save
	cd frontend; typings uninstall es6-promise --save

or

	# installs into typings/globals/promise
	cd frontend; typings install dt~promise --global --save
	cd frontend; typings uninstall promise --global --save

NOTE: in package above `dt~promise` uninstall lists package name without source (dt~).

### Update to TypeScript 2.0

It turns out that "lig": ["es6"] is supported only by TypeScript 2.0.  I have TypeScript 1.8.10.

I uninstalled 1.8.10 with:

	sudo npm uninstall typescript -g

To install beta version:

	sudo npm install -g typescript@beta

VS Code needs to be told to use updated version of TypeScript through change in `.vscode\settings.json" file:

	"typescript.tsdk": "/usr/local/lib/node_modules/typescript/lib"



## Plan

Continue on visual design through web UI development.

