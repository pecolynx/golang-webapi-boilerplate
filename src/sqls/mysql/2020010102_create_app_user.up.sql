create table `app_user` (
 `id` int auto_increment
,`version` int not null
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp on update current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`login_id` varchar(200) character set ascii
,`hashed_password` varchar(200) character set ascii
,`username` varchar(40)
,`removed` tinyint(1) not null
,primary key(`id`)
,unique(`login_id`)
);
