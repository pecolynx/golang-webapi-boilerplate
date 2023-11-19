create table `user_document` (
 `created_at` datetime not null default current_timestamp
,`created_by` int not null
,`app_user_id` int not null
,`document_id` varchar(26) not null
,primary key(`app_user_id`, `document_id`)
,foreign key(`created_by`) references `app_user`(`id`) on delete cascade
,foreign key(`app_user_id`) references `app_user`(`id`) on delete cascade
,foreign key(`document_id`) references `document`(`id`) on delete cascade
);
