create extension if not exists "uuid-ossp";

create table if not exists "users" (
    "id" uuid primary key default uuid_generate_v4(),
    "username" varchar(64) not null,
    "pseudonym" varchar(64) not null,
    "first_name" varchar(128) not null,
    "last_name" varchar(128) not null,
    "is_deleted" boolean not null default false,
    "created_at" timestamp not null default now(),
    "updated_at" timestamp not null default now()
);

alter table "users" add constraint "users_username_key" unique (username);
create index on "users" ("pseudonym");
create index on "users" ("username");
create index on "users" ("first_name");
create index on "users" ("last_name");

create type "admin_scale" as enum ('minor', 'major');

create table "users_admins" (
  "user_id" uuid primary key,
  "scale" admin_scale not null,
  "created_at" timestamp not null default now()
);

alter table "users_admins" add foreign key ("user_id") references "users" ("id");