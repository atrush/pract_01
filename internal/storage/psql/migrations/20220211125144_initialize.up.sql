CREATE TABLE users (
    id uuid not null,
    primary key (id)
);

CREATE TABLE urls (
    id uuid not null,
    user_id uuid not null,
    srcurl varchar(2050) not null,
    shorturl varchar (16) not null,
    unique (shorturl),
    primary key (id),
    foreign key (user_id) references users (id)
);