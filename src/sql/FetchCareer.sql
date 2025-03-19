-- Get a specific career entry
select careers.id, title, companies.name, description, startmonth, startyear, endmonth, endyear, current 
from careers
left join companies on companies.id = careers.companyid
where careers.id=$1
;