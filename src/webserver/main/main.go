package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	components_content "webserver.ethanrandolph.com/components/content"
	components_core "webserver.ethanrandolph.com/components/core"
	"webserver.ethanrandolph.com/datalayer"
)

type Application struct {
	dl *datalayer.Datalayer
}

func DBErrorHandle(err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		if e := new(pgconn.ConnectError); errors.As(err, &e) {
			w.WriteHeader(503)
		} else if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(404)
		}
	}
}

func (app *Application) CreateBase(content templ.Component, bio datalayer.FetchBiographyResult) templ.Component {
	header := components_core.Header(bio.Firstname, bio.Lastname, bio.PortaitLink)
	footer := components_core.Footer(bio.ResumeLink, bio.Email, bio.Linkedinlink, bio.Githublink)
	base := components_core.Base(header, content, footer)
	return base
}

func (app *Application) home_page(w http.ResponseWriter, r *http.Request) {

	Bio, err := app.dl.FetchBio()
	DBErrorHandle(err, w, r)
	proj_summaries, err := app.dl.FetchProjectSummaries(5)
	DBErrorHandle(err, w, r)
	car_summaries, err := app.dl.FetchCareerSummaries(3)
	DBErrorHandle(err, w, r)
	edu_summaries, err := app.dl.FetchEducationSummaries(3)
	DBErrorHandle(err, w, r)

	//Create Bio snippet
	biocomp := components_core.Biosnippet(Bio.Description)

	var proj_items []templ.Component
	//Create items for project summary
	for _, item := range proj_summaries {
		proj_items = append(proj_items, components_core.Summarysnippet_project(strconv.Itoa(item.Id), item.Name, item.Description, item.Thumbnaillink, 200, false))
	}
	var car_items []templ.Component
	for _, item := range car_summaries {
		car_items = append(car_items, components_core.Summarysnippet_career(item.Id, item.Title, item.Name, item.Description, 200))
	}
	var edu_items []templ.Component
	for _, item := range edu_summaries {
		edu_items = append(edu_items, components_core.Summarysnippet_education(item.Title, item.Major, fmt.Sprintf("%d - %d", item.StartDate, item.EndDate)))
	}

	proj_list := components_core.Summaryverticalcontainer(proj_items)
	car_list := components_core.Summaryverticalcontainer(car_items)
	edu_list := components_core.Summaryverticalcontainer(edu_items)

	content := components_content.Home(biocomp, proj_list, car_list, edu_list)
	base := app.CreateBase(content, Bio)
	base.Render(r.Context(), w)
}

func (app *Application) projectPage(w http.ResponseWriter, r *http.Request) {

	Bio, err := app.dl.FetchBio()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	id := r.PathValue("id")
	project, err := app.dl.FetchProject(id)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	projectimages, err := app.dl.FetchProjectImages(id)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	projecttools, err := app.dl.FetchProjectTools(id)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	var toolnames []string
	for _, tool := range projecttools {
		toolnames = append(toolnames, tool.Name)
	}
	projectvideos, err := app.dl.FetchProjectVideos(id)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	var videos []string
	var images []components_core.ImageInfo
	for _, video := range projectvideos {
		videos = append(videos, video.VideoYoutubeID)
	}
	for _, image := range projectimages {
		images = append(images, components_core.ImageInfo{ImageLink: image.Imagelink, ImageThumbnail: image.Imagethumbnaillink})
	}

	content := components_content.Project(project.Name, project.Repolink, project.Sitelink, project.Companyname,
		"TODO", toolnames, fmt.Sprintf("%d - %d", project.Startyear, project.Endyear),
		project.Description, images, videos)

	base := app.CreateBase(content, Bio)
	base.Render(r.Context(), w)
}

func (app *Application) careerPage(w http.ResponseWriter, r *http.Request) {
	Bio, err := app.dl.FetchBio()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	id := r.PathValue("id")
	career, err := app.dl.FetchCareer(id)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	page := components_content.Career(career.Title, career.Companyname, career.Description, fmt.Sprintf("%s %d - %s %d", career.Startmonth, career.Startyear, career.Endmonth, career.Endyear))
	base := app.CreateBase(page, Bio)
	base.Render(r.Context(), w)
}

func (app *Application) projectSummariesPage(w http.ResponseWriter, r *http.Request) {

	Bio, err := app.dl.FetchBio()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	summaries, err := app.dl.FetchProjectSummariesExtra(-1)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	//Finish
	var topitems []templ.Component
	var botitems []templ.Component
	var i int
	var summary datalayer.FetchProjectSummaryResult
	for i, summary = range summaries {
		if summary.IsCareer {
			break
		}
		item := components_core.Summarysnippet_project(strconv.Itoa(summary.Id), summary.Name, summary.Description, summary.Thumbnaillink, 500, true)
		topitems = append(topitems, item)
	}
	for j := i; j < len(summaries); j++ {
		item := components_core.Summarysnippet_project(strconv.Itoa(summaries[j].Id), summaries[j].Name, summaries[j].Description, summaries[j].Thumbnaillink, 500, true)
		botitems = append(botitems, item)
	}

	content := components_content.Summarypage_split("Projects", "Click on any entry below for detailed information", "Career Projects", topitems, botitems)
	base := app.CreateBase(content, Bio)
	base.Render(r.Context(), w)
}

func (app *Application) careerSummariesPage(w http.ResponseWriter, r *http.Request) {

	Bio, err := app.dl.FetchBio()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	summaries, err := app.dl.FetchCareerSummaries(-1)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	//Finish
	var items []templ.Component
	for _, summary := range summaries {
		item := components_core.Summarysnippet_career(summary.Id, summary.Title, summary.Name, summary.Description, 500)
		items = append(items, item)
	}

	content := components_content.Summarypage("Career", "Click on any entry below for detailed information", items)
	base := app.CreateBase(content, Bio)
	base.Render(r.Context(), w)
}

func (app *Application) educationSummariesPage(w http.ResponseWriter, r *http.Request) {
	Bio, err := app.dl.FetchBio()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	summaries, err := app.dl.FetchEducationSummaries(10)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	//Finish
	var items []templ.Component
	for _, summary := range summaries {
		item := components_core.Summarysnippet_education(summary.Title, summary.Major, fmt.Sprintf("%d - %d", summary.StartDate, summary.EndDate))
		items = append(items, item)
	}

	content := components_content.Summarypage("Education", "", items)
	base := app.CreateBase(content, Bio)
	base.Render(r.Context(), w)

}

func (app *Application) menuPage(w http.ResponseWriter, r *http.Request) {
	Bio, err := app.dl.FetchBio()
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	content := components_content.Categories()
	base := app.CreateBase(content, Bio)
	base.Render(r.Context(), w)
}

type mediaJsonPayload struct {
	VideoUrls     []string
	ImageUrls     []string
	ThumbnailUrls []string
}

func (app *Application) projectMediaJson(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	gopayload := mediaJsonPayload{}
	videos, err := app.dl.FetchProjectVideos(id)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	images, err := app.dl.FetchProjectImages(id)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	for _, video := range videos {
		gopayload.VideoUrls = append(gopayload.VideoUrls, video.VideoYoutubeID)
	}
	for _, image := range images {
		gopayload.ImageUrls = append(gopayload.ImageUrls, image.Imagelink)
		gopayload.ThumbnailUrls = append(gopayload.ThumbnailUrls, image.Imagethumbnaillink)
	}
	payload, err := json.Marshal(gopayload)
	if err != nil {
		fmt.Printf("Error occurred when serializing JSON: %s", err.Error())
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(payload)
}

func (app *Application) Layer1MiddleMan(w http.ResponseWriter, r *http.Request) {

	isReachable := app.dl.DBConnectionTest()
	if !isReachable {
		errPage := components_core.ServiceUnavailable()
		w.WriteHeader(503)
		errPage.Render(r.Context(), w)
	} else {
		mux := http.NewServeMux()
		mux.HandleFunc("/", app.home_page)
		mux.HandleFunc("/projects", app.projectSummariesPage)
		mux.HandleFunc("/careers", app.careerSummariesPage)
		mux.HandleFunc("/education", app.educationSummariesPage)
		mux.HandleFunc("/projects/{id}", app.projectPage)
		mux.HandleFunc("/projects/{id}/media.json", app.projectMediaJson)
		mux.HandleFunc("/career/{id}", app.careerPage)
		mux.HandleFunc("/menu", app.menuPage)
		mux.ServeHTTP(w, r)
	}
}

func main() {
	dbip := flag.String("dbip", os.Getenv("PORTFOLIOSERVER_DBIP"), "Postgresql Database IP Address")
	dbport := flag.String("dbport", os.Getenv("PORTFOLIOSERVER_DBPORT"), "Postgresql Database Port")
	dbuser := flag.String("dbuser", os.Getenv("PORTFOLIOSERVER_DBUSER"), "Postgresql Database login username")
	dbpass := flag.String("dbpass", os.Getenv("PORTFOLIOSERVER_DBPASS"), "Postgresql Database login password")
	dbname := flag.String("dbname", os.Getenv("PORTFOLIOSERVER_DBNAME"), "Postgresql Database Name")
	port := os.Getenv("PORTFOLIOSERVER_PORT")
	if port == "" {
		port = *flag.String("port", "4000", "Portfolio server port")
	}
	certpath := flag.String("cert", os.Getenv("PORTFOLIOSERVER_CERT"), "Portfolio Server Certification File")
	keypath := flag.String("key", os.Getenv("PORTFOLIOSERVER_KEY"), "Portfolio Server Key File")
	flag.Parse()

	app := Application{}
	dl := datalayer.Init("", *dbip, *dbport, *dbname, *dbuser, *dbpass)
	app.dl = dl

	fileServer := http.FileServer(http.Dir("../css/"))
	fontServer := http.FileServer(http.Dir("../fonts/"))
	jsServer := http.FileServer(http.Dir("../js/"))
	iconServer := http.FileServer(http.Dir("../icons/"))
	// Create new http.ServeMux. If server is running, at minimum,
	// handle basic file fetching routes that are not reliant on DB
	mux := http.NewServeMux()
	mux.Handle("GET /css/", http.StripPrefix("/css", fileServer))
	mux.Handle("GET /fonts/", http.StripPrefix("/fonts", fontServer))
	mux.Handle("GET /js/", http.StripPrefix("/js", jsServer))
	mux.Handle("GET /icons/", http.StripPrefix("/icons", iconServer))
	// Catch all pattern if none of the other patterns match. This
	// will handle specific pages tied to DB results.
	mux.HandleFunc("/", app.Layer1MiddleMan)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	// Able to run with or without TLS encryption, depending if the certificate files can be found.
	var runerr error
	if *certpath != "" && *keypath != "" {
		fmt.Println("Running with tls certificate...")
		runerr = server.ListenAndServeTLS(*certpath, *keypath)
	} else {
		fmt.Println("Running without tls certificate...")
		runerr = server.ListenAndServe()
	}

	if runerr != nil {
		fmt.Println(runerr.Error())
	}
}
