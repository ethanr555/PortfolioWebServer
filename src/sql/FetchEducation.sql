-- Gets all degrees earned, and sorts it by the year earned.
select Universities.name, Universities.link, Degrees.title, Degrees.major, Degrees.gpa, Degrees.startdate, Degrees.enddate
from Degrees
    left join Universities on Universities.id = Degrees.universityid
order by Degrees.enddate desc
;