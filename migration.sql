create table task
(
	id serial
		constraint url_pk
			primary key,
	url varchar(255) not null,
	count int not null,
	success int default 0 not null,
	error int default 0 not null,
);
