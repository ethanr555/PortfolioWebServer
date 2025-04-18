package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/a-h/templ"
)

type Application struct {
	dl Datalayer
}

func (app Application) CreateBase(content templ.Component, bio FetchBiographyResult) templ.Component {
	header := comp_header(bio.Firstname, bio.Lastname, bio.PortaitLink)
	footer := comp_footer(bio.ResumeLink, bio.Email, bio.Linkedinlink, bio.Githublink)
	base := comp_base(header, content, footer)
	return base
}

func (app Application) home_page(w http.ResponseWriter, r *http.Request) {

	Bio := app.dl.FetchBio()
	proj_summaries := app.dl.FetchProjectSummaries(5)
	car_summaries := app.dl.FetchCareerSummaries(3)
	edu_summaries := app.dl.FetchEducationSummaries(3)

	//Create Bio snippet
	biocomp := comp_biosnippet(Bio.Description)

	var proj_items []templ.Component
	//Create items for project summary
	for _, item := range proj_summaries {
		proj_items = append(proj_items, comp_summarysnippet_project(strconv.Itoa(item.Id), item.Name, item.Description, item.Thumbnaillink, 200, false))
	}
	var car_items []templ.Component
	for _, item := range car_summaries {
		car_items = append(car_items, comp_summarysnippet_career(item.Id, item.Title, item.Name, item.Description, 200))
	}
	var edu_items []templ.Component
	for _, item := range edu_summaries {
		edu_items = append(edu_items, comp_summarysnippet_education(item.Title, item.Major, fmt.Sprintf("%d - %d", item.StartDate, item.EndDate)))
	}

	proj_list := comp_summaryverticalcontainer(proj_items)
	car_list := comp_summaryverticalcontainer(car_items)
	edu_list := comp_summaryverticalcontainer(edu_items)

	content := comp_home(biocomp, proj_list, car_list, edu_list)
	base := app.CreateBase(content, Bio)
	base.Render(r.Context(), w)
}

func (app Application) projectPage(w http.ResponseWriter, r *http.Request) {

	Bio := app.dl.FetchBio()
	id := r.PathValue("id")
	project := app.dl.FetchProject(id)
	projectimages := app.dl.FetchProjectImages(id)
	projecttools := app.dl.FetchProjectTools(id)
	var toolnames []string
	for _, tool := range projecttools {
		toolnames = append(toolnames, tool.Name)
	}
	projectvideos := app.dl.FetchProjectVideos(id)
	var videos []string
	var images []ImageInfo
	for _, video := range projectvideos {
		videos = append(videos, video.VideoYoutubeID)
	}
	for _, image := range projectimages {
		images = append(images, ImageInfo{ImageLink: image.Imagelink, ImageThumbnail: image.Imagethumbnaillink})
	}

	content := comp_project(project.Name, project.Repolink, project.Sitelink, project.Companyname,
		"TODO", toolnames, fmt.Sprintf("%d - %d", project.Startyear, project.Endyear),
		project.Description, images, videos)

	base := app.CreateBase(content, Bio)
	base.Render(r.Context(), w)
}

func (app Application) careerPage(w http.ResponseWriter, r *http.Request) {
	Bio := app.dl.FetchBio()
	id := r.PathValue("id")
	career := app.dl.FetchCareer(id)
	page := comp_career(career.Title, career.Companyname, career.Description, fmt.Sprintf("%s %d - %s %d", career.Startmonth, career.Startyear, career.Endmonth, career.Endyear))
	base := app.CreateBase(page, Bio)
	base.Render(r.Context(), w)
}

func (app Application) educationPage(w http.ResponseWriter, r *http.Request) {

}

func (app Application) projectSummariesPage(w http.ResponseWriter, r *http.Request) {

	Bio := app.dl.FetchBio()
	summaries := app.dl.FetchProjectSummariesExtra(-1)

	//Finish
	var topitems []templ.Component
	var botitems []templ.Component
	var i int
	var summary FetchProjectSummaryResult
	for i, summary = range summaries {
		if summary.IsCareer {
			break
		}
		item := comp_summarysnippet_project(strconv.Itoa(summary.Id), summary.Name, summary.Description, summary.Thumbnaillink, 500, true)
		topitems = append(topitems, item)
	}
	for j := i; j < len(summaries); j++ {
		item := comp_summarysnippet_project(strconv.Itoa(summaries[j].Id), summaries[j].Name, summaries[j].Description, summaries[j].Thumbnaillink, 500, true)
		botitems = append(botitems, item)
	}

	content := comp_summarypage_split("Projects", "Click on any entry below for detailed information", "Career Projects", topitems, botitems)
	base := app.CreateBase(content, Bio)
	base.Render(r.Context(), w)
}

func (app Application) careerSummariesPage(w http.ResponseWriter, r *http.Request) {

	Bio := app.dl.FetchBio()
	summaries := app.dl.FetchCareerSummaries(-1)

	//Finish
	var items []templ.Component
	for _, summary := range summaries {
		item := comp_summarysnippet_career(summary.Id, summary.Title, summary.Name, summary.Description, 500)
		items = append(items, item)
	}

	content := comp_summarypage("Career", "Click on any entry below for detailed information", items)
	base := app.CreateBase(content, Bio)
	base.Render(r.Context(), w)
}

func (app Application) educationSummariesPage(w http.ResponseWriter, r *http.Request) {
	Bio := app.dl.FetchBio()
	summaries := app.dl.FetchEducationSummaries(10)

	//Finish
	var items []templ.Component
	for _, summary := range summaries {
		item := comp_summarysnippet_education(summary.Title, summary.Major, fmt.Sprintf("%d - %d", summary.StartDate, summary.EndDate))
		items = append(items, item)
	}

	content := comp_summarypage("Education", "", items)
	base := app.CreateBase(content, Bio)
	base.Render(r.Context(), w)

}

func main() {
	mux := http.NewServeMux()
	app := Application{}
	dl := Datalayer{}
	dbip := flag.String("dbip", os.Getenv("PORTFOLIOSERVER_DBIP"), "Postgresql Database IP Address")
	dbport := flag.String("dbport", os.Getenv("PORTFOLIOSERVER_DBPORT"), "Postgresql Database Port")
	dbuser := flag.String("dbuser", os.Getenv("PORTFOLIOSERVER_DBUSER"), "Postgresql Database login username")
	dbpass := flag.String("dbpass", os.Getenv("PORTFOLIOSERVER_DBPASS"), "Postgresql Database login password")
	dbname := flag.String("dbname", os.Getenv("PORTFOLIOSERVER_DBNAME"), "Postgresql Database Name")
	port := os.Getenv("PORTFOLIOSERVER_PORT")
	if port == "" {
		port = *flag.String("port", "4000", "Portfolio server port")
	}
	flag.Parse()
	dl.Init("", *dbip, *dbport, *dbname, *dbuser, *dbpass)
	app.dl = dl
	mux.HandleFunc("/", app.home_page)
	mux.HandleFunc("/projects", app.projectSummariesPage)
	mux.HandleFunc("/careers", app.careerSummariesPage)
	mux.HandleFunc("/education", app.educationSummariesPage)
	mux.HandleFunc("/projects/{id}", app.projectPage)
	mux.HandleFunc("/career/{id}", app.careerPage)
	//mux.HandleFunc("/education/{id}", app.educationPage)
	fileServer := http.FileServer(http.Dir("../css/"))
	fontServer := http.FileServer(http.Dir("../fonts/"))
	jsServer := http.FileServer(http.Dir("../js/"))
	mux.Handle("GET /css/", http.StripPrefix("/css", fileServer))
	mux.Handle("GET /fonts/", http.StripPrefix("/fonts", fontServer))
	mux.Handle("GET /js/", http.StripPrefix("/js", jsServer))
	//http.ListenAndServe(fmt.Sprintf(":%s", port), mux)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println(err.Error())
	}
}
