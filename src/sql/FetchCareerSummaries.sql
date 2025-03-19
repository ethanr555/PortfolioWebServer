-- Fetch just the summary information for the careers, sort for jobs that ended latest
select careers.id, title, companies.name, description
from careers
inner join companies on companies.id = careers.companyid
order by endyear desc, endmonth desc
limit $1
;