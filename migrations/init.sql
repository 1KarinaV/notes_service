create extension if not exists "uuid-ossp";

create table if not exists "users" (
    id uuid primary key,
    login varchar(40) unique not null,
    password varchar(255) not null
);

create table if not exists "notes" (
    id uuid primary key,
    user_id uuid not null,
    header varchar(200) not null,
    content text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    unique(user_id, header),

    constraint fk_users
        foreign key(user_id)
        references users(id)
);