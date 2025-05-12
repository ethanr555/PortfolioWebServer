
// images: {'ImageLink': string, 'ImageThumbnail': string}[]

// Creates an image/video carousel in parent element
class Carousel {
    //id: Parent Element ID
    //images: array of structs for images, as specified above
    //videoids: array of strings containing a youtube video id
    constructor(id, images, videoids) {
        this.id = id
        this.images = images
        this.videoids = videoids
        this.currentselectedthumbnail = null
        this.thumbnailclass = "cursor-pointer m-2 snap-center" 
        this.selectedthumbnailclass = this.thumbnailclass + " " + " outline-4 outline-blue-500"
        this.gallery = null
        this.Create()
    }

    //Function to switch currently viewed media to specified video
    //thumbnail: Thumbnail element that is being clicked
    //videoid: Id of video to switch to
    SwitchVideo(thumbnail, videoid) {
        // Change CSS of previous selected thumbnail to default, then make the current selected thumbnail the recently clicked one.
        this.currentselectedthumbnail.className = this.thumbnailclass
        this.currentselectedthumbnail = thumbnail
        thumbnail.focus()
        thumbnail.scrollIntoView({block: "nearest", inline: "nearest"})
        thumbnail.className = this.selectedthumbnailclass
        // Reset gallery to be empty
        this.gallery.innerHTML = '';
        const framewrap = document.createElement('div')
        framewrap.className = "relative pb-[56.25%] h-0 overflow-hidden"
        const youtubeframe = document.createElement('iframe');
        youtubeframe.className = "absolute top-0 left-0 w-full h-full"
        youtubeframe.src = 'https://www.youtube-nocookie.com/embed/' + videoid
        youtubeframe.allowFullscreen = ''
        framewrap.appendChild(youtubeframe)
        this.gallery.appendChild(framewrap)
    }

    //Function to switch currently viewed media to image
    //thumbnail: Thumbnail element that is being clicked
    //imageobj: Image struct of image that will be presented
    SwitchImage(thumbnail, imageobj) {
        this.gallery.innerHTML = '';
        this.currentselectedthumbnail.className = this.thumbnailclass
        thumbnail.focus()
        thumbnail.scrollIntoView({block: "nearest", inline: "nearest"})
        this.currentselectedthumbnail = thumbnail
        thumbnail.className = this.selectedthumbnailclass 
        const imageelement = document.createElement('img');
        imageelement.src = imageobj.ImageLink;
        imageelement.className = "h-100 min-w-100 object-contain"
        this.gallery.appendChild(imageelement);
        
    }

    //Creates the carousel
    Create() {
        const carouselbase = document.createElement('div')
        carouselbase.className = "block pb-20"
        const tophalf = document.createElement('div')
        tophalf.className = "flex justify-center items-center"
        this.gallery = document.createElement('div')
        this.gallery.className = "min-h-[66vw] min-w-lvw lg:min-h-100 lg:min-w-150 overflow-hidden"
        const thumbnailselection = document.createElement('div')
        thumbnailselection.className = "flex h-50 min-w-50 overflow-x-scroll snap-mandatory snap-x"
        let imagesobj = JSON.parse(this.images)
        for (let i = 0; i < this.videoids.length; i++)
        {
            let thumbnail = document.createElement('img')
            thumbnail.src = 'https://img.youtube.com/vi/' + this.videoids[i] + '/maxresdefault.jpg'
            thumbnail.className = this.thumbnailclass
            thumbnail.onclick = () => this.SwitchVideo(thumbnail, this.videoids[i])
            thumbnailselection.appendChild(thumbnail)
        }
        for (let i = 0; i < imagesobj.length; i++)
        {
            let thumbnail = document.createElement('img')
            thumbnail.src = imagesobj[i].ImageThumbnail
            thumbnail.className = this.thumbnailclass
            thumbnail.onclick = () => this.SwitchImage(thumbnail, imagesobj[i])
            thumbnailselection.appendChild(thumbnail)
        }
        tophalf.appendChild(this.gallery)
        carouselbase.appendChild(tophalf)
        carouselbase.appendChild(thumbnailselection)
        const parent = document.getElementById(this.id)
        parent.appendChild(carouselbase)
        this.currentselectedthumbnail = thumbnailselection.firstElementChild
        thumbnailselection.firstElementChild.click()
    }
}


