select videos.id, videos.videoid
from projectvideos
    inner join videos on videos.id = projectvideos.videoid
where projectvideos.projectid = $1
;