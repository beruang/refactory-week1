create table if not exists notes."user"
(
    id          serial not null
    constraint user_pk
    primary key,
    first_name  varchar,
    last_name   varchar,
    email       varchar,
    password    varchar,
    username    varchar,
    is_verified boolean default false
);

create unique index if not exists user_id_uindex
    on notes."user" (id);

create unique index if not exists user_username_uindex
    on notes."user" (username);