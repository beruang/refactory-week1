create table if not exists notes.roles
(
    id        serial               not null
    constraint roles_pk
    primary key,
    name      varchar              not null,
    is_active boolean default true not null
);

create unique index if not exists roles_id_uindex
    on notes.roles (id);

create unique index if not exists roles_name_uindex
    on notes.roles (name);