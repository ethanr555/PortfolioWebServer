
// images: {'ImageLink': string, 'ImageThumbnail': string}[]

// The spagehtti increases expontentially
function CreateCarousel(id,images, videoids) {
    const carouselbase = document.createElement('div')
    carouselbase.className = "block pb-20"
    const tophalf = document.createElement('div')
    tophalf.className = "flex justify-center items-center"
    const gallery = document.createElement('div')
    gallery.className = "min-h-[66vw] min-w-lvw lg:min-h-100 lg:min-w-150 overflow-hidden"
    const thumbnailselection = document.createElement('div')
    thumbnailselection.className = "flex h-50 min-w-50 overflow-x-scroll snap-mandatory snap-x"
    const thumbnailclass = "cursor-pointer m-2 snap-center"
    let imagesobj = JSON.parse(images)
    for (let i = 0; i < videoids.length; i++)
    {
        let thumbnail = document.createElement('img')
        thumbnail.src = 'https://img.youtube.com/vi/' + videoids[i] + '/maxresdefault.jpg'
        thumbnail.className = thumbnailclass
        thumbnail.onclick = function () {
            console.log(id)
            for (let i = 0; i < thumbnailselection.childNodes.length; i++)
            {
                thumbnailselection.childNodes[i].className = thumbnailclass
            }
            thumbnail.className = thumbnailclass + " outline-4 outline-blue-500"
            gallery.innerHTML = '';
            const framewrap = document.createElement('div')
            framewrap.className = "relative pb-[56.25%] h-0 overflow-hidden"
            const youtubeframe = document.createElement('iframe');
            youtubeframe.className = "absolute top-0 left-0 w-full h-full"
            youtubeframe.src = 'https://www.youtube-nocookie.com/embed/' + videoids[i]
            youtubeframe.allowFullscreen = ''
            framewrap.appendChild(youtubeframe)
            gallery.appendChild(framewrap)
        }
        thumbnailselection.appendChild(thumbnail)
    }

    for (let i = 0; i < imagesobj.length; i++)
    {
        let thumbnail = document.createElement('img')
        thumbnail.src = imagesobj[i].ImageThumbnail
        thumbnail.className = thumbnailclass
        thumbnail.onclick = function () {
            gallery.innerHTML = '';
            for (let i = 0; i < thumbnailselection.childNodes.length; i++)
            {
                thumbnailselection.childNodes[i].className = thumbnailclass
            }
            thumbnail.className = thumbnailclass + " outline-4 outline-blue-500"
            const imageelement = document.createElement('img');
            imageelement.src = imagesobj[i].ImageLink;
            imageelement.className = "h-100 min-w-100 object-contain"
            gallery.appendChild(imageelement);
        }
        thumbnailselection.appendChild(thumbnail)
    }

    tophalf.appendChild(gallery)
    carouselbase.appendChild(tophalf)
    carouselbase.appendChild(thumbnailselection)
    const parent = document.getElementById(id)
    parent.appendChild(carouselbase)
    thumbnailselection.firstElementChild.click()
    console.log("End!")
}


