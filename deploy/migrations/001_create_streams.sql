CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table video_streams(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  title varchar not null,
  created timestamp not null,
  updated timestamp not null
);

create table questions(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  stream uuid REFERENCES video_streams(id),
  text varchar not null
);

create table answers(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  question uuid REFERENCES questions(id),
  text varchar not null,
  correct boolean not null
);

---- create above / drop below ----

drop table answers;
drop table questions;
drop table video_streams;
