select id, imagelink, imagethumbnaillink
from images
where images.projectid = $1
;