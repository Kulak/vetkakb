-- In general either host or path shall be set to a non-empty value, not both.
-- If both are set, then more than one record can match query that locates
-- site by either value.  That query result shall resolve to only one client.
create table if not exists site(
	siteID integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	host text default '' not null,
	-- If path is '', then it is not set.
	path text default '' not null,
	dbname text not null,
	theme text default 'basic' not null,
	created timestamp default (strftime('%s', 'now')) NOT NULL,
	updated timestamp default (strftime('%s', 'now')) NOT NULL
);

insert into site (dbname) values('default');
