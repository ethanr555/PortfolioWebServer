
// Creates a drop-down menu to display various webpage directories. Can be opened and closed with a click.
class MobileNavMenu {
    constructor(id) {
        this.id = id
        this.isExpanded = false
        this.menuinstance = null
    }  
  
    // This will remove the HTML elements and deleted the current menu instance object. 
    DestroyMenu() {
        if (this.menuinstance != null) {
            this.menuinstance.remove()
            this.menuinstance = null
            this.isExpanded = false
        }
    }

    // This will create a new menu instance object and attach it to the DOM
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
    }

    // This will toggle whether the menu is expanded or not.
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

