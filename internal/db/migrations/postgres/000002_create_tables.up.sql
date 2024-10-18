-- auth::orgs::create
create table orgs (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  name varchar(255) not null,
  slug varchar(255) not null
);

-- auth::orgs::trigger
create trigger update_orgs_updated_at
  after update on orgs
  for each row
  execute function update_updated_at();

-- auth::teams::create
create table teams (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  org_id uuid not null references orgs (id),
  name varchar(255) not null,
  slug varchar(255) not null
);

-- auth::teams::trigger
create trigger update_teams_updated_at
  after update on teams
  for each row
  execute function update_updated_at();

-- auth::users::create
create table users (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  org_id uuid not null references orgs (id),
  email varchar(255) not null,
  first_name varchar(255),
  last_name varchar(255),
  password text,
  is_active boolean not null default true,
  is_verified boolean not null default false
);

-- auth::users::trigger
create trigger update_users_updated_at
  after update on users
  for each row
  execute function update_updated_at();

-- auth::oauth_accounts::create
create table oauth_accounts (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  user_id uuid not null references users (id),
  provider varchar(255) not null,
  provider_account_id varchar(255) not null,
  expires_at timestamptz,
  type varchar(255)
);

-- auth::oauth_accounts::trigger
create trigger update_oauth_accounts_updated_at
  after update on oauth_accounts
  for each row
  execute function update_updated_at();

-- auth::team_users::create
create table team_users (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  team_id uuid not null references teams (id),
  user_id uuid not null references users (id),
  role varchar(255),
  is_active boolean not null default true,
  is_admin boolean not null default false
);

-- auth::team_users::trigger
create trigger update_team_users_updated_at
  after update on team_users
  for each row
  execute function update_updated_at();

-- core::repos::create
create table repos (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  org_id uuid not null references orgs (id),
  name varchar(255) not null,
  provider varchar(255) not null,
  provider_id varchar(255) not null,
  default_branch varchar(255),
  is_monorepo boolean,
  threshold integer,
  stale_duration interval
);

-- core::repos::trigger
create trigger update_repos_updated_at
  after update on repos
  for each row
  execute function update_updated_at();

-- integrations/github::github_installations::create
create table github_installations (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  org_id uuid not null references orgs (id),
  installation_id bigint not null,
  installation_login varchar(255) not null,
  installation_login_id bigint not null,
  installation_type varchar(255),
  sender_id bigint not null,
  sender_login varchar(255) not null,
  status varchar(255)
);

-- integrations/github::github_installations::index
create unique index github_installations_installation_id_idx on github_installations (installation_id);

-- integrations/github::github_installations::trigger
create trigger update_github_installations_updated_at
  after update on github_installations
  for each row
  execute function update_updated_at();

-- integrations/github::github_orgs::create
create table github_orgs (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  installation_id uuid not null references github_installations (id),
  github_org_id bigint not null,
  name varchar(255) not null
);

-- integrations/github::github_orgs::index
create index github_orgs_installation_id_idx on github_orgs (installation_id);

-- integrations/github::github_orgs::trigger
create trigger update_github_orgs_updated_at
  after update on github_orgs
  for each row
  execute function update_updated_at();

-- integrations/github::github_users::create
create table github_users (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  user_id uuid references users (id),
  github_id bigint not null,
  github_org_id uuid not null references github_orgs (id),
  login varchar(255) not null
);

-- integrations/github::github_users::trigger
create trigger update_github_users_updated_at
  after update on github_users
  for each row
  execute function update_updated_at();

-- integrations/github::github_repos::create
create table github_repos (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  repo_id uuid not null references repos (id),
  installation_id uuid not null references github_installations (id),
  github_id bigint not null,
  name varchar(255) not null,
  full_name varchar(255) not null,
  url varchar(255) not null,
  is_active boolean
);

-- integrations/github::github_repos::index
create index github_repos_installation_id_idx on github_repos (installation_id);

-- integrations/github::github_repos::trigger
create trigger update_github_repos_updated_at
  after update on github_repos
  for each row
  execute function update_updated_at();

-- messaging::messaging::create
create table messaging (
  id uuid primary key default uuid_generate_v7(),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  provider varchar(255) not null,
  kind varchar(255) not null,
  link_to uuid not null,
  data jsonb not null
);

-- messaging::messaging::trigger
create trigger update_messaging_updated_at
  after update on messaging
  for each row
  execute function update_updated_at();

-- events::flat_events::create
create type event_provider as enum ('github', 'slack');

create table flat_events (
  id uuid primary key default uuid_generate_v7(),
  version varchar(255) not null,
  parent_id uuid,
  provider event_provider not null,
  scope varchar(255) not null,
  action varchar(255) not null,
  source varchar(255) not null,
  subject_id uuid not null,
  subject_name varchar(255) not null,
  payload jsonb not null,
  team_id uuid not null references teams (id),
  user_id uuid not null references users (id),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

-- events::flat_events::trigger
create trigger update_flat_events_updated_at
  after update on flat_events
  for each row
  execute function update_updated_at();