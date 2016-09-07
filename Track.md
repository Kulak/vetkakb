*** Vetka KB Project Track ***

<!-- TOC -->

- [Design](#design)
	- [Data Model](#data-model)
		- [Entry](#entry)
	- [Datastore Design](#datastore-design)
		- [FTS Search Conditions](#fts-search-conditions)
		- [SQL Column Affinity](#sql-column-affinity)
		- [SQLite3 Version](#sqlite3-version)
		- [SQLite3 'if not exists'](#sqlite3-if-not-exists)
		- [SQLite Length](#sqlite-length)
	- [UI Design](#ui-design)
	- [UUID](#uuid)
	- [Functional Approaches](#functional-approaches)
	- [External vs Internal Blob](#external-vs-internal-blob)
	- [HTML Tokenization](#html-tokenization)
- [Log](#log)
	- [TypeScript: namespaces and modules](#typescript-namespaces-and-modules)
	- [Fetch API Support](#fetch-api-support)
	- [Seach Typings](#seach-typings)
	- [Correct way to install typings](#correct-way-to-install-typings)
	- [Update to TypeScript 2.0](#update-to-typescript-20)
	- [Promise API](#promise-api)
- [Plan](#plan)

<!-- /TOC -->

# Design

## Data Model

### Entry

* ID as primary key
* Title, always filled
* RawType, always filled
* Raw as blob, always filled
* Plain as text for indexing or NULL if N/A
* Html as text for display or NULL if N/A
* Tags as text for indexing by tag keyword

## Datastore Design


### FTS Search Conditions

Porter tokenizer has the following conditions that apply to tags format restrictions.

Tags format:

1. Each tag is space separated.
   DO NOT use comma (,) to separate multiple words.
2. To create a phrase from multiple words use underscore
   DO NOT USE dash (-) to separate multiple words as tokenizer will fail to match a phrase.

use rowid or docid to manipulate unique id

### SQL Column Affinity

	select articleId, typeof(articleId) from article;

### SQLite3 Version

* Mac OS X runs SQLite 3.8.8.1 2015-01-20 16:51:25 f73337e3e289915a76ca96e7a05a1a8d4e890d55
* FreeBSD run sqlite3-3.13.0

### SQLite3 'if not exists'

"if not exists" added in SQLite version 3.7.11,

### SQLite Length

If the declared type of the column contains any of the strings "CHAR", "CLOB", or "TEXT" then that column has TEXT affinity. Notice that the type VARCHAR contains the string "CHAR" and is thus assigned TEXT affinity.

Note that numeric arguments in parentheses that following the type name (ex: "VARCHAR(255)") are ignored by SQLite - SQLite does not impose any length restrictions (other than the large global SQLITE_MAX_LENGTH limit) on the length of strings, BLOBs or numeric values.

## UI Design

Search List
	- "Most Recent" entries button.
	- "Search Field" to just type to search.
	- "Entrie List" displays search results.

## UUID

The choice is https://godoc.org/github.com/satori/go.uuid#UUID.Variant because it is used the most.  The alternative is https://godoc.org/github.com/twinj/uuid#UUID

It is the best store UUID as 16 byte array (128 bit) in a blob.  It is not the optimal solution to use that blob as the primary key, because it does not increment monotonously.

Type 3 and Type 5 UUIDs are just a technique of stuffing a hash into a UUID.

* Type 1: stuffs MAC address, datetime into 128 bits
* Type 3: stuffs an MD5 hash into 128 bits
* Type 4: stuffs random data into 128 bits
* Type 5: stuffs an SHA1 hash into 128 bits

Because UUID can be reliably regenerated it is not necessary to use do it right now.

## Functional Approaches


## External vs Internal Blob

Embedded blob storage is more optimal for data that's less than 20K.

File system storage is more optimal for 50K file sizes and above.  At 100K and above it is guarnateed to be faster with file system.  Source: <https://www.sqlite.org/intern-v-extern-blob.html>

## HTML Tokenization

The simplest way to avoid writing custom HTML tokenizer is to convert HTML to text with lynx text browser:

	curl -s http://www.sqlite.org | lynx -nolist -stdin -dump

Source: <https://www.mail-archive.com/sqlite-users@mailinglists.sqlite.org/msg82126.html>

Other options are available.

# Log

## TypeScript: namespaces and modules

Namespace can be defined in any directory.  To create a complex name separate it with dot as in the following example:

	namespace app.common {
		export class CommonData {}
	}

To reference namespace from another file sitting in the parent directory:

	import './common/test.js';
	new app.common.CommonData();

The above example uses relative to current client file URL.  I am **NOT** sure if

	import './common/test.ts';

would work better.

Other forms of import are possible.


Code without namespace is referenceable as is.

About module imports:
https://www.typescriptlang.org/docs/handbook/module-resolution.html


## Fetch API Support

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

## Seach Typings

	typings search aName

Note that source column lists name of the source.  NPM is default (?).  If source is DT, then it has to be specified in install command.  See install `dt~promise` below.

## Correct way to install typings

These commands were rolled back, but they represent good example of how to work with typings:

	# installes into typings/modules/es6-promise
	cd frontend; typings install es6-promise --save
	cd frontend; typings uninstall es6-promise --save

or

	# installs into typings/globals/promise
	cd frontend; typings install dt~promise --global --save
	cd frontend; typings uninstall promise --global --save

NOTE: in package above `dt~promise` uninstall lists package name without source (dt~).

## Update to TypeScript 2.0

It turns out that "lig": ["es6"] is supported only by TypeScript 2.0.  I have TypeScript 1.8.10.

I uninstalled 1.8.10 with:

	sudo npm uninstall typescript -g

To install beta version:

	sudo npm install -g typescript@beta

VS Code needs to be told to use updated version of TypeScript through change in `.vscode\settings.json" file:

	"typescript.tsdk": "/usr/local/lib/node_modules/typescript/lib"

## Promise API

About chains: <http://stackoverflow.com/questions/32032588/how-to-return-from-a-promises-catch-then-block>


# Plan

Issues:  During image upload we need to preserve Content-Type as it is the easiest way to manage it.

Binary Image vs Custom is messed up.

Cleanup: entry rawType shall be removed and replaced with rawTypeName, because DB service can resolve number to a name now.

Add title to entrySearch table, because it will result in more uniform experience.  Since operator match works on entire table at once use of new column is the best.