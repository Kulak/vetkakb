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
	expiresAt timestamp,
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
	-- URL to titleIcon
	titleIcon text not null default '',
	html text,
	-- intro is a teaser text or short description displayed in the list of entries.
	-- It does not have to be set.  If it is set it will be used.
	-- If it is not set, then Intro can be generated from plain text.
	-- Intro must be plain text only and not markdown or any other format.
	intro text not null default '',
	-- owner has implicit access
	ownerFK integer NOT NULL,
	-- if BITWISE AND requiredClearance with user clearance is zero, then user has no access.
	-- if not zero, then user has access
	-- default access clearance access mask allows access to Administrator
	requiredClearance integer NOT NULL default 0x8,
	created timestamp NOT NULL default (strftime('%s', 'now')),
	updated timestamp NOT NULL default (strftime('%s', 'now')),
	-- Tracks when clearance of entry is set to 0x1 (Guest accessible)
	published timestamp default NULL
);

create virtual table if not exists entrySearch using fts4 (
	entryFK integer primary key,
	title,
	plain,
	tags,
	tokenize=porter
);

-- redirect table represents alternative path that maps into entry
-- It is especially useful for old path redirection
create table if not exists redirect (
	redirectID integer primary key,
	path text not null,
	entryFK integer not null,
	statusCode integer not null
);