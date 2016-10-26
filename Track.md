*** Vetka KB Project Track ***

<!-- TOC -->

- [Design](#design)
	- [Data Model](#data-model)
		- [Entry](#entry)
		- [Product Subsystem Design](#product-subsystem-design)
			- [Entry Properties mapped to Product Description:](#entry-properties-mapped-to-product-description)
			- [EntrySearch Properties mapped to Product Description:](#entrysearch-properties-mapped-to-product-description)
			- [Product Properties](#product-properties)
			- [Digital Download Properties](#digital-download-properties)
	- [Datastore Design](#datastore-design)
		- [FTS Search Conditions](#fts-search-conditions)
		- [SQL Column Affinity](#sql-column-affinity)
		- [SQLite3 Version](#sqlite3-version)
		- [SQLite3 'if not exists'](#sqlite3-if-not-exists)
		- [SQLite3 Length](#sqlite3-length)
		- [SQLite3 Copy Record into History Table](#sqlite3-copy-record-into-history-table)
	- [UI Design](#ui-design)
	- [UUID](#uuid)
	- [Functional Approaches](#functional-approaches)
	- [External vs Internal Blob](#external-vs-internal-blob)
	- [HTML Tokenization](#html-tokenization)
	- [Google OAuth](#google-oauth)
	- [Permalink, redirect, alternative URL, slug](#permalink-redirect-alternative-url-slug)
- [Alternatives: Custom Entry Type](#alternatives-custom-entry-type)
	- [Custom Data is only in Raw blob](#custom-data-is-only-in-raw-blob)
	- [Custom data is stored in custom table](#custom-data-is-stored-in-custom-table)
- [Log](#log)
	- [TypeScript: namespaces and modules](#typescript-namespaces-and-modules)
	- [Fetch API Support](#fetch-api-support)
	- [Seach Typings](#seach-typings)
	- [Correct way to install typings](#correct-way-to-install-typings)
	- [Update to TypeScript 2.0](#update-to-typescript-20)
	- [Promise API](#promise-api)
- [Version 1.0 Capability](#version-10-capability)
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

### Product Subsystem Design

To be able to reuse existing Entry UI controls, we choose to use
Entry with it current capabilities mapped to Product parameters.
The mapping from Product to Entry is natural.

Entry is the center of the star design.  Product and Digital downloads
are edges of the star.

Shall Product and Digital Download be its own microservices?
If edges are microservices, then additional server join or
REST client call are necessary.

In a star design it is impossible to know what extensions are in
use in which entry.  Thus a separate call to load each extension
is desirable.  This separate call can be made on the server side
or on the client side.  That justifies Service oriented approach to
edges of the star.

To simplify deployment service is going to be built-into the
main application process.  At the same time UI depends on product...

#### Entry Properties mapped to Product Description:

* `entryID` is product ID
* `raw` is contains primary product description.
* `rawType` indicates raw format (markdown, plain, HTML, etc.)
* `rawConentType` is a pass-through parameter for files.
  Content-Type is recorded as is in `rawConentType` and then
	served back when binary file is served.
* `rawFileName` is not in use
* `titleIcon` is path to title image
* `html` rendered HTML text for `raw` content.
* `intro` is a one paragraph of plain text introduction
* `ownerFK` is the owner of the product description entry
* `requiredClearance` reflects who can access (read or write) to the entry (?)
* `slug` is used as usual

#### EntrySearch Properties mapped to Product Description:

* `title` is a product title or name
* `plain` is an aggregated search field
* `tas` is used as usual

#### Product Properties

Product is a standalone data table that contains product specific parameters
unrelated to digital downloads and not covered in Entry or EntrySearch tables.

Product can be associated with any entry.

* `imagePaths` is an array of strings, each string is URL paths to product image
* `price`
* `SKU`
* `VendorID`
* `paypalID`
* `etc`

#### Digital Download Properties

Each product may have more than one digital download associated with it.
Digital downloads are controlled by per user security access policy.

* `ID` unique ID of the download file.
* `entryID` points to a product this download belongs to.
* `fileContent` is a download file content.
* `fileName` is a name of the download file.
* `revision` is a sequencial revision number which goes up with each updated
* `updated` is a date of when download file was updated
* `created`

Digital Download is not specific to the Product.  It can be associated with
any Entry.

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

### SQLite3 Length

If the declared type of the column contains any of the strings "CHAR", "CLOB", or "TEXT" then that column has TEXT affinity. Notice that the type VARCHAR contains the string "CHAR" and is thus assigned TEXT affinity.

Note that numeric arguments in parentheses that following the type name (ex: "VARCHAR(255)") are ignored by SQLite - SQLite does not impose any length restrictions (other than the large global SQLITE_MAX_LENGTH limit) on the length of strings, BLOBs or numeric values.

### SQLite3 Copy Record into History Table

From /Users/sergei/SWDev/py-projs/a-rebeccaslp.com/extsw/alitkaPy/alitka/model/product.py file:

	insert or fail into productInfoHistory (path, productInfoId, productId, name, iconUrl, price, downloadUrl, buyHtml,
	intro, description, type, fileRev, fileNameBase, fileNameExt, created, updated, published, productType)
	select path, productInfoId, productId, name, iconUrl, price, downloadUrl, buyHtml,
	intro, description, type, fileRev, fileNameBase, fileNameExt, created, updated, published, productType
	from productInfo where productInfoId = ?;


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

## Google OAuth

https://console.developers.google.com/apis/credentials?project=vetka-142905

## Permalink, redirect, alternative URL, slug

These are all terms for alternative URLs used to point to an entry.

# Alternatives: Custom Entry Type

The following were my initial thoughts.  However, designing
product and digital downloads exposed the need to create star based
design to extend softly the original Entry concept.

Star design keeps original UI controls unmodified and pushes
extension effort into new composite elements such as Product and Digital Download.

## Custom Data is only in Raw blob

Ideally raw blob stores all original custom data.

Cons: SQL queries cannot be used to implement custom search condition
on blob data.

Pros: standard search mechanism can be used as usual.  There is less
confusion about the model.

## Custom data is stored in custom table

This solution negaes cons of the RAW only extension.

However, custom query code has to be written to support lookup.

In a way this route can be taken later as search optimization solution and not as
a core data storage solution.

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

# Version 1.0 Capability

Core functionality is fully implemented in entry and entrySearch tables.

Features:

* plain text
* html
* markdown
* binary image file
* any other binary file
* can display images

Missing features:

* Authentication is not available
* Permission system is missing, because authentication is not there

# Plan

* Need to test redirection again.