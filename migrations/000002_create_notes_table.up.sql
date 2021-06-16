create table notes.notes
(
    id serial not null,
    user_id int not null
        constraint notes_user_id_fk
            references notes."user",
    type varchar not null,
    title varchar,
    body varchar,
    secret varchar,
    is_active bool default true not null
);

create unique index notes_id_uindex
	on notes.notes (id);

alter table notes.notes
    add constraint notes_pk
        primary key (id);