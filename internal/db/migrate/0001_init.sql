-- SQL migration script for initializing the database
create table users (
  id integer primary key autoincrement,
  email text unique not null,
  password_hash text,
  role text not null check (role in ('admin','manager','user')),
  created_at datetime not null default current_timestamp
);

create table tasks (
  id integer primary key autoincrement,
  owner_id integer not null references users(id) on delete cascade,
  title text not null,
  description text,
  status text not null default 'todo' check (status in ('todo','doing','done')),
  due_date date,
  created_at datetime not null default current_timestamp,
  updated_at datetime not null default current_timestamp
);
