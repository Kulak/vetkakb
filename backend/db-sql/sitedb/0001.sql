create table if not exists site(
	siteID integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	host text default '' not null,
	path text default '' not null,
	dbname text not null,
	theme text default 'basic' not null,
	created timestamp default (strftime('%s', 'now')) NOT NULL,
	updated timestamp default (strftime('%s', 'now')) NOT NULL
);

insert into site (dbname) values('default');
