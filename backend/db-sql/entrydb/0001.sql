
create table if not exists entry(
	entryID integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	title text unique not null,
	rawType integer not null,
	raw blob,
	html text,
	created timestamp default (strftime('%s', 'now')) NOT NULL,
	updated timestamp default (strftime('%s', 'now')) NOT NULL
);

create virtual table if not exists entrySearch using fts4 (
	entryFK integer primary key,
	plain,
	tags,
	tokenize=porter
);
