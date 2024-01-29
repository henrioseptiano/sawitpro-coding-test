/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */
create table users (
   id serial PRIMARY KEY,
   user_id text NOT null,
   full_name char(60) NOT NULL,
   phone_number char(13) NOT NULL,
   password text NOT null,
   successfull_login_attempts bigint not null,
   last_login timestamp null,
   created_at timestamp not null,
   updated_at timestamp not null
);

-- password : maulana
INSERT INTO public.users (id, user_id, full_name, phone_number, "password", successfull_login_attempts, last_login, created_at, updated_at) VALUES(2, 'd9982291-e467-4594-ab1c-18d1e2d7bbc1', 'maulana', '+6278231212', '$2a$10$mDMtvDh4opF/dzjO1W4v2ePoEbJafSYjlXqkNgGvCsokGd7qaO462', 3, '2024-01-29 01:27:44.996', '2024-01-29 01:00:00.851', '2024-01-29 01:00:00.851');
