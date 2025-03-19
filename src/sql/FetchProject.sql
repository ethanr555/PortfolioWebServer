-- Fetch an individual project
-- int, string, string, string, string, string, string, string, int, int
select projects.id, projects.name, companies.name, description, codereponame, coderepolink, projectsitename, projectsitelink, startyear, endyear
from projects
left join companies on companies.id = projects.companyid
where projects.id=$1
;