create table `user_ticket` (
 `created_at` datetime not null default current_timestamp
,`created_by` int not null
,`app_user_id` int not null
,`ticket_id` int not null
,primary key(`app_user_id`, `ticket_id`)
,foreign key(`created_by`) references `app_user`(`id`) on delete cascade
,foreign key(`app_user_id`) references `app_user`(`id`) on delete cascade
,foreign key(`ticket_id`) references `ticket`(`id`) on delete cascade
);
