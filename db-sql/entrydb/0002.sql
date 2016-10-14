alter table entry add column slug text not null default '';

update entry
set slug=(
	select replace(replace(replace(replace(replace(replace(title,
		',', ''), ' ', '-'), '?', ''), '!', ''), '(', ''), ')', '')
	from entrySearch where entry.entryId=entrySearch.entryFK
);

create unique index entrySlugIdx on entry(slug);