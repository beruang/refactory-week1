create table if not exists notes.media
(
    id serial not null,
    mime_type varchar not null,
    file bytea not null
);

create unique index media_id_uindex
	on notes.media (id);

alter table notes.media
    add constraint media_pk
        primary key (id);

