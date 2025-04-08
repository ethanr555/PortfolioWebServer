-- Finds projects for project summaries, ordered from newest to oldest.
(
select      id, name, thumbnaillink, description, endyear, false as iscareer
from        projects
where companyid is null
)
union
(
select      id, name, thumbnaillink, description, endyear, true as iscareer
from        projects
where companyid is not null
)
order by    iscareer asc, endyear desc
    limit $1
;