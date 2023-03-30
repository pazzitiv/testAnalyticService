CREATE TABLE user_actions
(
    id      serial PRIMARY KEY,
    user_id varchar(255) not null,
    data    jsonb not null ,
    time    timestamp not null
);