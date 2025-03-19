-- Gets all tools associated with the specified project
select tools.id, tools.name
from projects
    inner join projecttools on projects.id = projecttools.projectid
    inner join tools on projecttools.toolid = tools.id
where projects.id = $1
;