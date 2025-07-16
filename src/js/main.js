
// Do not include import statements, as those do not get pruned by minify

// This will search for key elements in page to determine if certain javascript objects get created and rendered in the DOM.
function main() {
    // Find if there is an element requesting this:
    potentialElement = document.getElementById("carousel")
    if (potentialElement != null) {
        const mediaURL = window.location.href + "/media.json"
        const mediaRequest = new Request(mediaURL)
        fetch(mediaRequest)
        .then((Response) => Response.json())
        .then((data) => {

            let images = []
            for (let i = 0; i < data.ThumbnailUrls.length; i++)
            {
                images.push({"ImageLink": data.ImageUrls[i], "ImageThumbnail": data.ThumbnailUrls[i]})
            }

            const car = new Carousel("carousel", images, data.VideoUrls)
        })
    }
    
    if (document.getElementById("nav")) {
        const menuhandler = new MobileNavMenu("nav")
        const containerdiv = document.createElement('div')
        containerdiv.className = "lg:hidden h-auto w-10 ml-auto mr-auto"
        const imageicon = document.createElement('img')
        imageicon.src = "/icons/hamburgermenuicon.svg"
        imageicon.alt = "Menu"
        imageicon.className = "border-transparent border-b-2 box-border hover:underline hover:underline-offset-2 hover:decoration-2 hover:cursor-pointer"
        imageicon.onclick = () => menuhandler.ToggleMenu()
        containerdiv.appendChild(imageicon)
        document.getElementById('nav').appendChild(containerdiv)           
    }
}

main()