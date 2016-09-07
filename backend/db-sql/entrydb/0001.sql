
create table if not exists entry(
	entryID integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	raw blob,
	rawType integer not null,
	rawContentType text,
	rawFileName text,
	html text,
	created timestamp default (strftime('%s', 'now')) NOT NULL,
	updated timestamp default (strftime('%s', 'now')) NOT NULL
);

create virtual table if not exists entrySearch using fts4 (
	entryFK integer primary key,
	title,
	plain,
	tags,
	tokenize=porter
);
