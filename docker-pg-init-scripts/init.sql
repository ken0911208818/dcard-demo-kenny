/*
    https://docs.postgresql.tw/reference/sql-commands/create-role
    CREATE ROLE 定義一個資料庫的角色
    NOSUPERUSER : 一般使用者 若是 SUPERUSER 則為超級使用者
    INHERIT : 自動使用已授予其直接或間接成員的所有角色的任何資料庫特權 NOINHERIT 則為預設
    NOCREATEDB : 不可建立資料庫的角色
    NOCREATEROLE : 不可建立新User
    NOREPLICATION : 不可以進行複製工作的角色
 */
CREATE ROLE dcard_admin LOGIN PASSWORD 'admin_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;
CREATE ROLE dcard_user LOGIN PASSWORD 'user_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;
CREATE ROLE dcard_readonly LOGIN PASSWORD 'readonly_password' NOSUPERUSER INHERIT NOCREATEDB NOCREATEROLE NOREPLICATION;

/*
    setup db
    ENCODING : 資料庫字元編碼
    LC_COLLATE : 排序順序
    LC_CTYPE : 字元分類方式 影響 大小寫數字 等等
    CONNECTION LIMIT : 當兩個新連線幾乎同時開始 連線數只剩一個時 兩者則都有可能失敗 不會對超級使用者以及後台工作程序限制
    temple : 複製標準系統資料庫
    OWNER TO : 將擁有者給 to ??
    SET timezone : 設定時區
    REVOKE : 撤銷某個功能權限
    GRANT — 賦予存取權限
 */
CREATE DATABASE dcard_db with ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8' CONNECTION LIMIT = -1 template=template0;
ALTER DATABASE dcard_db OWNER TO dcard_admin;
ALTER DATABASE dcard_db SET timezone TO 'UTC';
REVOKE USAGE ON SCHEMA public FROM PUBLIC;
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public to dcard_admin;
GRANT CREATE ON SCHEMA public to dcard_admin;
/* 將所有 USAGE CREATE 從所有人轉移至 dcard_admin能使用*/
GRANT USAGE ON SCHEMA public to dcard_user;
GRANT USAGE ON SCHEMA public to dcard_readonly;
/*
    create table
 */
 --the script to remove all tables in the database
\connect dcard_db;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS pairs CASCADE;

create table users
(
    id uuid,
    email character varying(200) not null,
    password_digest character varying(1000) not null,
    name character varying(255) not null,
    create_time timestamp without time zone not null default current_timestamp,
    update_time timestamp without time zone not null default current_timestamp,

    CONSTRAINT "users_pk" PRIMARY KEY (id)
);
ALTER TABLE users ADD CONSTRAINT users_u1 UNIQUE (email);

create table pairs
(
    user_id_one uuid,
    user_id_two uuid,

    CONSTRAINT "pairs_pk" PRIMARY KEY (user_id_one, user_id_two)
);
ALTER TABLE pairs ADD CONSTRAINT pairs_u1 UNIQUE (user_id_one);

\connect dcard_db;
ALTER TABLE pairs ADD CONSTRAINT pairs_fk1 FOREIGN KEY (user_id_one) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE pairs ADD CONSTRAINT pairs_fk2 FOREIGN KEY (user_id_two) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;

\connect dcard_db;

/*
    grant_table_privilege
*/
/*for normal tables */
ALTER TABLE users OWNER TO dcard_admin;
ALTER TABLE pairs OWNER TO dcard_admin;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE users to dcard_user;
GRANT SELECT, INSERT, UPDATE, DELETE, REFERENCES ON TABLE pairs to dcard_user;
GRANT SELECT ON TABLE users to dcard_readonly;
GRANT SELECT ON TABLE pairs to dcard_readonly;

/*
    insert testing data
*/
\connect dcard_db;
-- each test user's plain password is 0000 --
insert into users(id, email, password_digest, name)
values('97327413-6b65-486f-b299-91be0871f898', 'ken@example.com', '$2a$10$gVtjNk4YL.O4I//ZBtvfN.YEebwR1Ci3.5OBHan4PWFzniSFqpzce', 'kenny');

insert into users(id, email, password_digest, name)
values('eb3c75df-b0df-4e06-a02f-e2ba77eba68a', 'nicole@example.com', '$2a$10$6tsb.2dRzV5gSTEJmtwkgeKpPIMO0VbMv2E6hP9xuAytwFlf0trVm', 'nicole');

insert into users(id, email, password_digest, name)
values('80695811-0bf2-44fd-980d-1635de7734a8', 'jack@example.com', '$2a$10$WkWwIpCbMyB1A2OuMC9LI.4LtQZtxNb1djcYqzeP0IayazJQgVkHG', 'jack');

insert into pairs(user_id_one, user_id_two)
values('97327413-6b65-486f-b299-91be0871f898', 'eb3c75df-b0df-4e06-a02f-e2ba77eba68a')