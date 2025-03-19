-- Get project tags associated with project
select tags.id, tags.name
from projects
    inner join projecttags on projects.id = projecttags.projectid
    inner join tags on projecttags.tagid = tags.id
where projects.id = $1
;