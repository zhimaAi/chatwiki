-- +goose Up

/* robot*/
update  "public"."chat_ai_robot" set robot_avatar = '/upload/default/robot_avatar.svg'  where robot_avatar = '/upload/default/robot_avatar.png' and application_type = 0;

/* workflow*/
update  "public"."chat_ai_robot" set robot_avatar = '/upload/default/workflow_avatar.svg'  where robot_avatar = '/upload/default/robot_avatar.png' and application_type = 1;

/*library*/

update "public"."chat_ai_library" set avatar = 'upload/default/library_normal_avatar.svg' where avatar = 'upload/default/library_avatar.png' and  "type" = 0;


/*open*/

update "public"."chat_ai_library" set avatar = 'upload/default/library_open_avatar.svg' where avatar = 'upload/default/library_avatar.png' and  "type" = 1;
/*qa*/

update "public"."chat_ai_library" set avatar = 'upload/default/library_qa_avatar.svg' where avatar = 'upload/default/library_avatar.png' and  "type" = 2;