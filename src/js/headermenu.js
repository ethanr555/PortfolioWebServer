
class MobileNavMenu {
    constructor(id) {
        this.id = id
        this.isExpanded = false
        this.menuinstance = null
    }  
   
    DestroyMenu() {
        if (this.menuinstance != null) {
            this.menuinstance.remove()
            this.menuinstance = null
            this.isExpanded = false
            console.log("ran 2")
        }
    }

    CreateMenu() {
        if (this.menuinstance != null) {
            this.DestroyMenu()
        }
        let parentElement = document.getElementById(this.id)
        let container = document.createElement('div')
        let homeelement = document.createElement('a')
        let projectselement = document.createElement('a')
        let careerelement = document.createElement('a')
        let educationelement = document.createElement('a')

        container.className = "flex flex-col items-center "
        homeelement.href = "/"
        homeelement.innerText = "Home"
        projectselement.href = "/projects"
        projectselement.innerText = "Projects"
        careerelement.href = "/careers"
        careerelement.innerText = "Career"
        educationelement.href = "/education"
        educationelement.innerText = "Education"
        
        let sharedclassname = "text-center text-white hover:underline"
        homeelement.className = projectselement.className = careerelement.className = educationelement.className = sharedclassname

        container.appendChild(homeelement)
        container.appendChild(projectselement)
        container.appendChild(careerelement)
        container.appendChild(educationelement)
        parentElement.appendChild(container)
        this.menuinstance = container
        this.isExpanded = true
        console.log("ran")
    }

    ToggleMenu() {
        if (!this.isExpanded) {
            this.CreateMenu()
        }
        else {
            this.DestroyMenu()
            this.isExpanded = false
        }
    }

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
