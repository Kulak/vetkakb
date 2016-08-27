
create table if not exists entry(
	entryID integer primary key,
	title text unique not null,
	rawType text not null,
	raw blob,
	html text,
	created datetime not null,
	updated datetime null
);

create virtual table if not exists entrySearch using fts4 (
	entryFK integer primary key,
	plain,
	tags,
	tokenize=porter
);
