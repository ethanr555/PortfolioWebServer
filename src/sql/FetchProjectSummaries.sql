-- Finds projects for project summaries, ordered from newest to oldest.
select      id, name, thumbnaillink, description
from        projects
order by    endyear desc
    limit $1
;