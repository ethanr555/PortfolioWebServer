# Portfolio Web Server

This is a personal portfolio web server that I wrote. I created this to show off my portfolio and demonstrate knowledge creating and operating full stack websites. My previous portfolio website was designed with React.js, with Github Pages handling the brunt of the back-end operations. This is not the case with this web server. This webserver is programmed with Go(lang). It uses the Templ Go module to handle HTML templating, and this server communicates with a PostgreSQL database for data access. I use the PGX module as the database driver for the server. The website is pre-rendered with templ and unlike my previous portfolio website, barely uses Javascript. It was determined that while it is good practice to redesign the website with a more traditional approach, the portfolio website does not need Javascript to render everything. The website is designed for both users with Javascript functionality and those without. In other words, it is designed to be accessible for various devices. Nonetheless, it still uses Javascript in select parts in order to have more dynamic website functionality (such as a video/image carousel). CSS is also created with a library called TailwindCSS. While raw CSS probably is enough to do the job, TailwindCSS does make it more intuitive to work with th design and helps with the REPL. 

# Initial Requirements

- Go 1.24
- CMake
- Ubuntu (or another Linux distro)
- wget
- Git

# Building

Once the initial requirements are meant, clone this project with

    git clone https://github.com/ethanr555/PortfolioWebServer
    cd PortfolioWebServer

Now, running CMake will setup the environment:

    make configure

This will setup dependencies for the Go project. It also installs minify as a Go binary. When you are ready to build it, run

    make

This will produce an executable in ./build/cmd/. Do note that this generates self-signed certificate files that the next command below will use. In addition, it will download a local copy of the standalone TailwindCSS executable into the tools file. This will only download once, unless the file is deleted. 

In order to run the server, ensure you run the executable in its directory, or from the root directory run

    make run

When the project needs to be purged from build/generated files, run

    make clean

# Running

Run the executable in build/cmd. It uses several command line arguments and/or environment variables:

- PORTFOLIOSERVER_DBIP: IP address of Postgresql database instance
- PORTFOLIOSERVER_DBPORT: Port of Postgresql database instance
- PORTFOLIOSERVER_DBNAME: Name of Postgresql database
- PORTFOLIOSERVER_DBUSER: User of Postgresql database to login with
- PORTFOLIOSERVER_DBPASS: Password of Postgresql database to login with
- PORTFOLIOSERVER_PORT: Port that the webserver should be hosted on
- PORTFOLIOSERVER_CERT: The file location for the certification file for HTTPS connection.
- PORTFOLIOSERVER_KEY: The file location for the certification key file for HTTPS connection.