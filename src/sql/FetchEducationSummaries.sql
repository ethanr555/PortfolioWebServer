-- Gets the degree title and major and orders it by year received. 
select degrees.title, degrees.major, degrees.startdate, degrees.enddate
from degrees
order by enddate desc
limit $1
;