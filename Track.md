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

## Plan

Continue on visual design through web UI development.

