create DATABASE mesos;

create table slave_info (
	hostname varchar(20) not null primary key, 
	attachment varchar(255) not null
);


create table task_info (
	task_cpu float not null, 
	task_mem float not null,
	id varchar(255) primary key, 
	name varchar(20) not null, 
	cmd varchar(255)  not null, 
	env varchar(255)  not null, 
	image varchar(50) not null, 
	slave_id varchar(255)  not null, 
	hostname varchar(20) not null, 
	framework_id varchar(255) not null, 
	status int not null, 
	count int not null
);
