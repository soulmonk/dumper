CREATE TABLE ideas
(
  id BIGSERIAL
    constraint user_pk
    primary key,
  title text,
  body  text
);
