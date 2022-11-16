create table account
(
    id            bigint unsigned auto_increment comment 'id',
    account_owner varchar(64) not null default '',
    balance       bigint      not null default 0,
    currency      varchar(8)  not null default '',
    created_at    datetime    not null default now(),
    updated_at    datetime    not null default now() on update now(),
    primary key (id) using btree
) engine = innodb comment = 'account';

create index account_owner_index on account(account_owner);
create unique index unique_account_index on account(account_owner, currency);

create table entry
(
    id         bigint unsigned auto_increment comment 'id',
    account_id bigint unsigned not null default 0,
    amount     bigint          not null default 0,
    created_at datetime        not null default now(),
    updated_at datetime        not null default now() on update now(),
    primary key (id) using btree
) engine innodb comment = 'entries';

create index entry_account_index on entry (account_id);

create table transfer
(
    id         bigint unsigned auto_increment comment 'id',
    from_id    bigint unsigned not null default 0,
    to_id      bigint unsigned not null default 0,
    amount     bigint          not null default 0,
    created_at datetime        not null default now(),
    updated_at datetime        not null default now() on update now(),
    primary key (id) using btree
) engine = innodb comment = 'transfer';

create index transfer_from_to_index on transfer (from_id, to_id);
create index transfer_to_index on transfer (to_id);

create table users
(
    id bigint unsigned auto_increment comment 'id',
    username varchar(32) not null default '',
    email varchar(32) not null default '',
    hashed_passwd varchar(64) not null default '',
    password_changed_at datetime not null default '0000-01-01 00:00:00',
    created_at datetime        not null default now(),
    updated_at datetime        not null default now() on update now(),
    primary key (id) using btree
) engine = innodb comment 'users';

create unique index unique_email_index on users(email);