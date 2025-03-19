# Portfolio Web Server

This is a personal portfolio web server that I wrote. I created this to show off my portfolio and demonstrate knowledge creating and operating full stack websites. My previous portfolio website was designed with React.js, with Github Pages handling the brunt of the back-end operations. This is not the case with this web server. This webserver is programmed with Go(lang). It uses the Templ Go module to handle HTML templating, and this server communicates with a PostgreSQL database for data access. I use the PGX module as the database driver for the server. The website is pre-rendered with templ and unlike my previous portfolio website, barely uses Javascript. It was determined that while it is good practice to redesign the website with a more traditional approach, the portfolio website does not need Javascript to render everything. The website is designed for both users with Javascript functionality and those without. In other words, it is designed to be accessible for various devices. Nonetheless, it still uses Javascript in select parts in order to have more dynamic website functionality (such as a video/image carousel). CSS is also created with a library called TailwindCSS. While raw CSS probably is enough to do the job, TailwindCSS does make it more intuitive to work with th design and helps with the REPL. 

# How to install

TBD

- TailwindCSS-CLI 4.1 (with NPM)
- Go 1.24
- Templ
- CMake (for Makefile)
- Debian-based Linux

# Building

TBD

Use the Makefile and run "make" to build/copy everything into a build folder. To individually build CSS file and templ file, use "make tailwindcss" and "make templ" respectively.

To build just one file, make sure to target the build location instead of the src location.

# Running

Run the executable in build/bin. It uses several command line arguments and/or environment variables:

- PORTFOLIOSERVER_DBIP: IP address of Postgresql database instance
- PORTFOLIOSERVER_DBPORT: Port of Postgresql database instance
- PORTFOLIOSERVER_DBNAME: Name of Postgresql database
- PORTFOLIOSERVER_DBUSER: User of Postgresql database to login with
- PORTFOLIOSERVER_DBPASS: Password of Postgresql database to login with
- PORTFOLIOSERVER_PORT: Port that the webserver should be hosted on