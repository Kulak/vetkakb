create table if not exists user (
	userID integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	-- one means a guest access; zero means account is disabled as no access is granted.
	clearances integer default 0x1 NOT NULL,
	created timestamp default (strftime('%s', 'now')) NOT NULL,
	updated timestamp default (strftime('%s', 'now')) NOT NULL
);

create table if not exists oauthUser (
	userFK integer NOT NULL PRIMARY KEY,
	provider text,
	email text,
	name text,
	firstName text,
	lastName text,
	nickName text,
	description text,
	provUserID text,
	avatarUrl text,
	location text,
	accessToken text,
	accessTokenSecret text,
	refreshToken text,
	expiresAt text,
	created timestamp default (strftime('%s', 'now')) NOT NULL,
	updated timestamp default (strftime('%s', 'now')) NOT NULL
);

-- this table is for documentation purpose
create table if not exists clearance (
	mask integer not null primary key,
	name text not null
);
insert into clearance (mask, name) values(0x8, 'Administrator');
-- Items with Guest access are effectively "published" to public
insert into clearance (mask, name) values(0x1, 'Guest');

create table if not exists entry(
	entryID integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	raw blob,
	rawType integer not null,
	rawContentType text,
	rawFileName text,
	html text,
	-- owner has implicit access
	ownerFK integer NOT NULL,
	-- if BITWISE AND requiredClearance with user clearance is zero, then user has no access.
	-- if not zero, then user has access
	-- default access clearance access mask allows access to Administrator
	requiredClearance integer default 0x8 NOT NULL,
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
