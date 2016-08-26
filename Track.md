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

	"lib": ["es6"],

This magically enabled ES6 Promise support and compilation error was gone.

If typing is installed with option --global, then it does not need to be imported (it is imported automatically).

The systemjs still needs to map polyfil name 'whatwg-fetch' to the actual library.  So, a change in HTML file was necessary.

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

## Plan

Continue on visual design through web UI development.

