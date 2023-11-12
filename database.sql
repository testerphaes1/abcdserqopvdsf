create table if not exists auth_info
(
    id          SERIAL primary key,
    private_key text
);

create table if not exists accounts
(
    id           SERIAL primary key,
    first_name   text,
    last_name    text,
    phone_number text,
    email        text,
    username     text,
    password     text,

    updated_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at   TIMESTAMP NOT NULL,
    deleted_at   TIMESTAMP,
    unique (username)
);

create table if not exists packages
(
    id             SERIAL primary key,
    title          text,
    price          int       NOT NULL,
    description    text,
    length_in_days int       NOT NULL,
    limits         jsonb,

    updated_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at     TIMESTAMP NOT NULL,
    deleted_at     TIMESTAMP
);

create table if not exists projects
(
    id            SERIAL primary key,
    title         text,
    is_active     bool,
    expire_at     TIMESTAMP,
    account_id    int       NOT NULL,
    package_id    int       NOT NULL,
    notifications jsonb,
    members       jsonb,

    updated_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at    TIMESTAMP NOT NULL,
    deleted_at    TIMESTAMP,

    foreign key (account_id) references accounts (id),
    foreign key (package_id) references packages (id)
);

create table if not exists endpoints
(
    id         SERIAL primary key,
    data       jsonb,
    project_id int       not null,
    disabled   text,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,

    foreign key (project_id) references projects (id)
);

create table if not exists trace_routes
(
    id         SERIAL primary key,
    data       jsonb,
    project_id int       not null,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,

    foreign key (project_id) references projects (id)
);

create table if not exists net_cats
(
    id         SERIAL primary key,
    data       jsonb,
    project_id int       not null,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,

    foreign key (project_id) references projects (id)
);

create table if not exists pings
(
    id         SERIAL primary key,
    data       jsonb,
    project_id int       not null,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,

    foreign key (project_id) references projects (id)
);

create table if not exists page_speeds
(
    id         SERIAL primary key,
    data       jsonb,
    project_id int       not null,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,

    foreign key (project_id) references projects (id)
);

create table if not exists drafts
(
    id         SERIAL primary key,
    data       text,
    project_id int       not null,
    type       text,

    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,

    foreign key (project_id) references projects (id)
);

create table if not exists datacenters
(
    id              SERIAL primary key,
    baseurl         text      not null,
    title           text      not null,
    connection_rate int,
    lat             float,
    lng             float,
    location_name   text,
    country_name    text,

    updated_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMP NOT NULL,
    deleted_at      TIMESTAMP
);

create table if not exists relation_datacenters
(
    id            SERIAL primary key,
    endpoint_id   int,
    datacenter_id int,

    updated_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at    TIMESTAMP NOT NULL,
    deleted_at    TIMESTAMP,

    foreign key (endpoint_id) references endpoints (id),
    foreign key (datacenter_id) references datacenters (id)
);

create table if not exists gateways
(
    id              SERIAL primary key,
    baseurl         text      not null,
    title           text      not null,
    connection_rate int,
    is_active       bool,
    is_default      bool,

    data            jsonb,

    updated_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMP NOT NULL,
    deleted_at      TIMESTAMP
);

create table if not exists orders
(
    id                      SERIAL primary key,
    account_id              int       not null,
    project_id              int       not null,
    package_id              int       not null,
    gateway_id              int       not null,
    status                  text      not null,
    amount                  int       not null,
    gateway_order_id        text,
    gateway_create_response jsonb,
    gateway_verify_response jsonb,

    foreign key (account_id) references accounts (id),
    foreign key (project_id) references projects (id),
    foreign key (package_id) references packages (id),
    foreign key (gateway_id) references gateways (id),

    created_at              TIMESTAMP NOT NULL,
    updated_at              TIMESTAMP NOT NULL,
    deleted_at              TIMESTAMP
);

create table if not exists tickets
(
    id            SERIAL primary key,
    account_id    int       not null,
    project_id    int,
    message       text,
    ticket_status int       not null,
    title         text,
    reply_to      int,

    foreign key (reply_to) references tickets (id),
    foreign key (project_id) references projects (id),
    foreign key (account_id) references accounts (id),

    created_at    TIMESTAMP NOT NULL,
    updated_at    TIMESTAMP NOT NULL,
    deleted_at    TIMESTAMP
);

create table if not exists faq
(
    id         SERIAL primary key,
    question   text,
    answer     text,

    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

-----------------------------------------------------------------------------------------

create table if not exists endpoint_stats
(
    time              TIMESTAMPTZ      NOT NULL,
    session_id        text             NOT NULL,
    project_id        int              NOT NULL,
    endpoint_name     text,
    endpoint_id       int              not null,
    url               text,
    datacenter_id     int              not null,
    is_heart_beat     bool             not null,
    success           int              not null,
    response_time     double precision not null,
    response_times    text,
    response_bodies   BYTEA,
    response_headers  text,
    response_statuses text,

    foreign key (project_id) references projects (id),
    foreign key (endpoint_id) references endpoints (id),
    foreign key (datacenter_id) references datacenters (id),

    PRIMARY KEY (time, endpoint_id)
);

create table if not exists net_cats_stats
(
    time          TIMESTAMPTZ NOT NULL,
    session_id    text        NOT NULL,
    project_id    int         NOT NULL,
    netcat_id     int         not null,
    url           text,
    datacenter_id int         not null,
    is_heart_beat bool        not null,
    success       int         not null,

    foreign key (project_id) references projects (id),
    foreign key (netcat_id) references net_cats (id),
    foreign key (datacenter_id) references datacenters (id),

    PRIMARY KEY (time, netcat_id)
);

create table if not exists page_speeds_stats
(
    time          TIMESTAMPTZ NOT NULL,
    session_id    text        NOT NULL,
    project_id    int         NOT NULL,
    pagespeed_id  int         not null,
    url           text,
    datacenter_id int         not null,
    is_heart_beat bool        not null,
    success       int         not null,

    foreign key (project_id) references projects (id),
    foreign key (pagespeed_id) references page_speeds (id),
    foreign key (datacenter_id) references datacenters (id),

    PRIMARY KEY (time, pagespeed_id)
);

create table if not exists pings_stats
(
    time          TIMESTAMPTZ NOT NULL,
    session_id    text        NOT NULL,
    project_id    int         NOT NULL,
    ping_id       int         not null,
    url           text,
    datacenter_id int         not null,
    is_heart_beat bool        not null,
    success       int         not null,

    foreign key (project_id) references projects (id),
    foreign key (ping_id) references pings (id),
    foreign key (datacenter_id) references datacenters (id),

    PRIMARY KEY (time, ping_id)
);

create table if not exists trace_routes_stats
(
    time          TIMESTAMPTZ NOT NULL,
    session_id    text        NOT NULL,
    project_id    int         NOT NULL,
    traceroute_id int         not null,
    url           text,
    datacenter_id int         not null,
    is_heart_beat bool        not null,
    success       int         not null,

    foreign key (project_id) references projects (id),
    foreign key (traceroute_id) references trace_routes (id),
    foreign key (datacenter_id) references datacenters (id),

    PRIMARY KEY (time, traceroute_id)
);

SELECT create_hypertable('endpoint_stats', 'time');
SELECT create_hypertable('net_cats_stats', 'time');
SELECT create_hypertable('page_speeds_stats', 'time');
SELECT create_hypertable('pings_stats', 'time');
SELECT create_hypertable('trace_routes_stats', 'time');

INSERT INTO accounts (first_name, last_name, phone_number, email, username, password, updated_at, created_at,
                      deleted_at)
VALUES ('mohammad'::text, 'safakhou'::text, '09337942924'::text, 'mohammad.sf220@gmail.com'::text, 'mnim'::text,
        'password'::text, '2022-10-29 11:32:50.000000'::timestamp, '2022-10-29 11:32:49.000000'::timestamp,
        null::timestamp);
INSERT INTO datacenters (baseurl, title, connection_rate, lat, lng, location_name, updated_at, created_at,
                         deleted_at)
VALUES ('http://93.113.233.131:10002'::text, 'default DS'::text, 1::integer, 35.7448459::double precision,
        51.3731325::double precision, 'asia tech'::text, '2022-10-29 11:51:02.000000'::timestamp,
        '2022-10-29 11:51:03.000000'::timestamp, null::timestamp);
INSERT INTO gateways (baseurl, title, connection_rate, is_active, is_default, data, updated_at, created_at,
                      deleted_at)
VALUES ('http://api.idpay.ir'::text, 'idpay'::text, 2::integer, true::boolean, true::boolean, '{
  "call_back_url": "http://93.113.233.131:222/idpay-callback-url"
}'::jsonb, '2022-10-29 11:52:15.000000'::timestamp, '2022-10-29 11:52:14.000000'::timestamp, null::timestamp);
INSERT INTO packages (title, price, description, length_in_days, limits, updated_at, created_at, deleted_at)
VALUES ('Free'::text, 0::integer, 'To Test Out Your Plans'::text, 0::integer, '{
  "ping_limits": {
    "duration_limit": 10,
    "number_of_monitoring": 1
  },
  "net_cat_limits": {
    "duration_limit": 0,
    "number_of_monitoring": 0
  },
  "endpoint_limits": {
    "duration_limit": 60,
    "number_of_monitoring": 1
  },
  "page_speed_limits": {
    "duration_limit": 0,
    "number_of_monitoring": 0
  },
  "trace_route_limits": {
    "duration_limit": 60,
    "number_of_monitoring": 1
  }
}'::jsonb, '2022-10-29 11:54:54.000000'::timestamp, '2022-10-29 11:54:51.000000'::timestamp, null::timestamp);
INSERT INTO packages (title, price, description, length_in_days, limits, updated_at, created_at, deleted_at)
VALUES ('Silver'::text, 5000000::integer, 'Small Businesses'::text, 30::integer, '{
  "ping_limits": {
    "duration_limit": 10,
    "number_of_monitoring": 10
  },
  "net_cat_limits": {
    "duration_limit": 10,
    "number_of_monitoring": 10
  },
  "endpoint_limits": {
    "duration_limit": 20,
    "number_of_monitoring": 10
  },
  "page_speed_limits": {
    "duration_limit": 0,
    "number_of_monitoring": 0
  },
  "trace_route_limits": {
    "duration_limit": 10,
    "number_of_monitoring": 10
  }
}'::jsonb, '2022-10-29 11:54:53.000000'::timestamp, '2022-10-29 11:54:50.000000'::timestamp, null::timestamp);
INSERT INTO packages (title, price, description, length_in_days, limits, updated_at, created_at, deleted_at)
VALUES ('Gold'::text, 15000000::integer, 'StartUp Engines'::text, 30::integer, '{
  "ping_limits": {
    "duration_limit": 5,
    "number_of_monitoring": 50
  },
  "net_cat_limits": {
    "duration_limit": 5,
    "number_of_monitoring": 50
  },
  "endpoint_limits": {
    "duration_limit": 5,
    "number_of_monitoring": 30
  },
  "page_speed_limits": {
    "duration_limit": 0,
    "number_of_monitoring": 0
  },
  "trace_route_limits": {
    "duration_limit": 5,
    "number_of_monitoring": 10
  }
}'::jsonb, '2022-10-29 11:54:52.000000'::timestamp, '2022-10-29 11:54:48.000000'::timestamp, null::timestamp);
INSERT INTO packages (title, price, description, length_in_days, limits, updated_at, created_at, deleted_at)
VALUES ('Platinum'::text, 30000000::integer, 'You Care About Future'::text, 30::integer, '{
  "ping_limits": {
    "duration_limit": 1,
    "number_of_monitoring": 100
  },
  "net_cat_limits": {
    "duration_limit": 1,
    "number_of_monitoring": 100
  },
  "endpoint_limits": {
    "duration_limit": 1,
    "number_of_monitoring": 50
  },
  "page_speed_limits": {
    "duration_limit": 0,
    "number_of_monitoring": 0
  },
  "trace_route_limits": {
    "duration_limit": 1,
    "number_of_monitoring": 100
  }
}'::jsonb, '2022-10-29 11:54:54.000000'::timestamp, '2022-10-29 11:54:50.000000'::timestamp, null::timestamp);
INSERT INTO projects (title, is_active, expire_at, account_id, package_id, notifications, updated_at, created_at,
                      deleted_at)
VALUES ('default'::text, true::boolean, '2100-01-01 00:00:00.000000'::timestamp, 1::integer, 1::integer, '{
  "email": [
    "mohammad.sf220@gmail.com"
  ],
  "slack": [
    ""
  ],
  "telegram": [
    ""
  ]
}'::jsonb, '2022-10-29 11:57:30.000000'::timestamp, '2022-10-29 11:57:31.000000'::timestamp, null::timestamp);
INSERT INTO public.endpoints (data, project_id, updated_at, created_at, deleted_at)
VALUES ('{
  "endpoints": [
    {
      "url": "https://www.digikala.com",
      "body": "",
      "header": {},
      "method": "GET",
      "endpoint_name": "https://www.digikala.com",
      "acceptance_model": {
        "statuses": [
          "200"
        ],
        "response_bodies": []
      }
    }
  ],
  "scheduling": {
    "end_at": "2022-12-18 16:55:26.433925 +0000 +0000",
    "duration": 1,
    "is_active": true,
    "project_id": 1,
    "pipeline_id": 0,
    "data_centers": [
      1
    ],
    "is_heart_beat": true
  }
}'::jsonb, 1::integer, '2022-10-29 12:37:14.000000'::timestamp, '2022-10-29 12:37:14.000000'::timestamp,
        null::timestamp);
INSERT INTO public.endpoints (data, project_id, updated_at, created_at, deleted_at)
VALUES ('{
  "endpoints": [
    {
      "url": "https://www.digikala.com",
      "body": "",
      "header": {},
      "method": "GET",
      "endpoint_name": "https://www.digikala.com",
      "acceptance_model": {
        "statuses": [
          "200"
        ],
        "response_bodies": []
      }
    }
  ],
  "scheduling": {
    "end_at": "2022-12-18 16:55:26.433925 +0000 +0000",
    "duration": 1,
    "is_active": true,
    "project_id": 1,
    "pipeline_id": 0,
    "data_centers": [
      1
    ],
    "is_heart_beat": false
  }
}'::jsonb, 1::integer, '2022-10-29 12:37:15.000000'::timestamp, '2022-10-29 12:37:16.000000'::timestamp,
        null::timestamp);
