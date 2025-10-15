create table users
(
    id   uuid,
    name varchar not null check ( CHAR_LENGTH(name) > 2 ),
    constraint users_pk primary key (id)
);

create table user_activity
(
    id          uuid,
    user_id     uuid         not null,
    action_date timestamp(3) not null,
    action      varchar      not null,
    metadata    jsonb,
    constraint user_activity_pk primary key (id),
    constraint user_activity_users_fk foreign key (user_id) references users (id) on delete cascade
);

create table user_activity_history
(
    period_id     uuid not null,
    user_id       uuid not null,
    actions_count bigint default 0 check ( actions_count >= 0 ),
    constraint user_activity_history_pk primary key (period_id, user_id),
    constraint user_activity_history_users_fk foreign key (user_id) references users (id) on delete cascade
);

create table activity_periods
(
    id        uuid,
    from_date timestamp(3) not null,
    to_date   timestamp(3) not null check ( to_date > from_date ),
    constraint activity_periods_pk primary key (id)
);