-- Delete user_story_reads for users that have read 1 story.
delete
from
	user_story_read
where
	user_id in (
		select
			user_id
		from
			user_story_read 
		group by
			user_id
		having
			count(story_id) < 2
	);

-- Delete all users that have no read stories.
delete
from
	app_user
where
	id in (
		select 
			au.id
		from
			app_user au
			left join user_story_read usr on au.id = usr.user_id
		where
			usr.user_id is null
	);
