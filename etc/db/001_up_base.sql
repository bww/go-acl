
create table acl_authorization (
  id              uuid                      primary key,
  key             varchar(64)               unique not null,
  secret          varchar(64)               not null,
  description     varchar(1024),
  active          boolean                   not null default false,
  created_at      timestamp with time zone  not null default now()
);

create table acl_policy (
  id              uuid                      primary key,
  type            varchar(32)               not null,
  data            jsonb                     not null,
  created_at      timestamp with time zone  not null default now()
);

create table acl_authorization_r_policy (
  auth_id         uuid                      not null references acl_authorization (id),
  policy_id       uuid                      not null references acl_policy (id),
  primary key (auth_id, policy_id)
);
