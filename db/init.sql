CREATE TABLE users
(
    uid serial NOT NULL,
    username character varying(100) NOT NULL UNIQUE,
    email_address character varying(100) NOT NULL UNIQUE,
    password character varying(500) NOT NULL,
    insert_date timestamp default current_timestamp,
    PRIMARY KEY (uid)
); 