-- Delete all users that have read 0 or 1 stories.
delete
from
    app_user
where
    id in (
        select
            usr.user_id
        from
            user_story_read usr
        group by
            usr.user_id
        having
            count(story_id) < 2
    );

-- Delete user_story_reads for non-existing users.
delete
from
    user_story_read
where
    user_id in (
        select
            usr.user_id
        from
            user_story_read usr
            left join app_user au on usr.user_id = au.id
        where
            au.id is null
    );
