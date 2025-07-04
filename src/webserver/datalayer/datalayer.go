/*
Package datalayer is the intermediate layer between the database and the rest of the server.

It maintains a connection with the database with PGX.
It contains functions representing specific queries to the database, along with specific structs to capture the results of the queries.
*/
package datalayer

import (
	"context"

	"fmt"

	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Type Datalayer is the main struct representing the database connection and acts as the receiver for all the querying functions.
// Do not create this directly, use Init() to create an initialized Datalayer
type Datalayer struct {
	pool         *pgxpool.Pool
	relativePath string
}

// The result struct for FetchProject()
type FetchProjectResult struct {
	Id          int    //Project ID
	Name        string //Project name
	Companyname string // Company name
	Description string // Project Description
	Reponame    string // Code repository name
	Repolink    string // Code repository link
	Sitename    string // Project Website Name
	Sitelink    string // Project Website Link
	Startyear   int    // Project Start Year
	Endyear     int    // Project End Year
}

// The result struct for FetchProjectImages()
type FetchProjectImagesResult struct {
	Id                 int    // Image ID
	Imagelink          string // Full image link
	Imagethumbnaillink string // Image thumbnail link
}

// The result struct for FetchProjectVideos()
type FetchProjectVideosResult struct {
	Id             int    // Video ID
	VideoYoutubeID string // Video Youtube ID
}

// The result struct for FetchProjectTags()
type FetchProjectTagsResult struct {
	Id   int    // Tag ID
	Name string // Tag Name
}

// The result struct for FetchProjectTools()
type FetchProjectToolsResult struct {
	Id   int    // Tool ID
	Name string // Tool Name
}

// The result struct for FetchCareer()
type FetchCareerResult struct {
	Id          int
	Title       string
	Companyname string
	Description string
	Startmonth  string
	Startyear   int
	Endmonth    string
	Endyear     int
	Current     bool
}

// The result struct for FetchEducation()
type FetchEducationResult struct {
	Name      string
	Link      string
	Title     string
	Major     string
	Gpa       string
	Startdate int
	Enddate   int
}

// The result struct for FetchProjectSummary()
type FetchProjectSummaryResult struct {
	Id            int
	Name          string
	Thumbnaillink string
	Description   string
	Endyear       int
	IsCareer      bool
}

// The result struct for FetchCareerSummaries()
type FetchCareerSummariesResult struct {
	Id          string
	Title       string
	Name        string
	Description string
}

// The result struct for FetchEducationSummaries()
type FetchEducationSummariesResult struct {
	Title     string
	Major     string
	StartDate int
	EndDate   int
}

// The result struct for FetchBiography()
type FetchBiographyResult struct {
	Firstname    string
	Lastname     string
	Description  string
	Email        string
	Linkedinlink string
	Githublink   string
	Websitelink  string
	PortaitLink  string
	ResumeLink   string
}

// This function properly creates a Datalayer struct with the provided parameters. Always create a Datalayer with this function.
func Init(path string, ip string, port string, databasename string, username string, password string) *Datalayer {
	dl := &Datalayer{}
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, ip, port, databasename)
	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {

		fmt.Println("Connection error!")
	}
	dl.pool = conn
	dl.relativePath = path
	return dl
}

// Tests the connection established by the pgx.pool to ensure that it is still connected or is connected.
func (dl *Datalayer) DBConnectionTest() bool {
	err := dl.pool.Ping(context.Background())
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return false
	}
	return true
}

// Closes the datalayer database connection
func (dl *Datalayer) Close() {
	if dl.pool == nil {
		fmt.Println("Error: Database connection does not exist and cannot be closed.")
	}
	dl.pool.Close()
}

// Reads a query from the file specified at the path.
func getQueryFromPath(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error with reading file: %s\n", err.Error())
		return ""
	}
	return string(data)
}

// A catchall function that will return either the value of the pointer or the equivalent zeroed value.
func resolveNil[T any](input *T) T {
	if input == nil {
		var empty T
		return empty
	} else {
		return *input
	}
}

// Writes a redundant error message format witht he specified function name
func errMessQuery(funcName string, err error) error {
	return fmt.Errorf("datalayer: query failed when executing function: %s: %w", funcName, err)
}

// Query that fetches a project (excluding images, tools, tags, videos)
func (dl *Datalayer) FetchProject(projectid string) (FetchProjectResult, error) {
	var id *int
	var name *string
	var companyname *string
	var description *string
	var reponame *string
	var repolink *string
	var sitename *string
	var sitelink *string
	var startyear *int
	var endyear *int
	query := getQueryFromPath("../sql/FetchProject.sql")
	err := dl.pool.QueryRow(context.Background(), query, projectid).Scan(&id, &name, &companyname, &description, &reponame, &repolink, &sitename, &sitelink, &startyear, &endyear)
	if err != nil {
		return FetchProjectResult{}, errMessQuery("FetchProject", err)
	}
	result := FetchProjectResult{
		Id:          resolveNil(id),
		Name:        resolveNil(name),
		Companyname: resolveNil(companyname),
		Description: resolveNil(description),
		Reponame:    resolveNil(reponame),
		Repolink:    resolveNil(repolink),
		Sitename:    resolveNil(sitename),
		Sitelink:    resolveNil(sitelink),
		Startyear:   resolveNil(startyear),
		Endyear:     resolveNil(endyear),
	}
	return result, nil
}

// Query that fetches project images
func (dl *Datalayer) FetchProjectImages(projectid string) ([]FetchProjectImagesResult, error) {
	var id *int
	var imagelink *string
	var imagethumbnaillink *string
	query := getQueryFromPath("../sql/FetchProjectImages.sql")
	rows, err := dl.pool.Query(context.Background(), query, projectid)
	if err != nil {
		errMessQuery("FetchProjectImages", err)
		return []FetchProjectImagesResult{}, err
	}
	var results []FetchProjectImagesResult
	pgx.ForEachRow(rows, []any{&id, &imagelink, &imagethumbnaillink}, func() error {
		result := FetchProjectImagesResult{
			Id:                 resolveNil(id),
			Imagelink:          resolveNil(imagelink),
			Imagethumbnaillink: resolveNil(imagethumbnaillink)}
		results = append(results, result)
		return nil
	})
	return results, nil
}

// Query that fetches project videos
func (dl *Datalayer) FetchProjectVideos(projectid string) ([]FetchProjectVideosResult, error) {
	var id *int
	var videoLink *string
	query := getQueryFromPath("../sql/FetchProjectVideos.sql")
	rows, err := dl.pool.Query(context.Background(), query, projectid)
	if err != nil {
		return []FetchProjectVideosResult{}, errMessQuery("FetchProjectVideos", err)
	}
	var results []FetchProjectVideosResult
	pgx.ForEachRow(rows, []any{&id, &videoLink}, func() error {
		result := FetchProjectVideosResult{
			Id:             resolveNil(id),
			VideoYoutubeID: resolveNil(videoLink)}
		results = append(results, result)
		return nil
	})
	return results, nil

}

// Query that fetches project tags
func (dl *Datalayer) FetchProjectTags(projectid int) ([]FetchProjectTagsResult, error) {
	var id *int
	var name *string
	query := getQueryFromPath("../sql/FetchProjectTags.sql")
	rows, err := dl.pool.Query(context.Background(), query, projectid)
	if err != nil {
		return []FetchProjectTagsResult{}, errMessQuery("FetchProjectTags", err)
	}
	var results []FetchProjectTagsResult
	pgx.ForEachRow(rows, []any{&id, &name}, func() error {
		result := FetchProjectTagsResult{
			Id:   resolveNil(id),
			Name: resolveNil(name)}
		results = append(results, result)
		return nil
	})
	return results, nil

}

// Query that fetches project tools
func (dl *Datalayer) FetchProjectTools(projectid string) ([]FetchProjectToolsResult, error) {

	var id *int
	var name *string
	query := getQueryFromPath("../sql/FetchProjectTools.sql")
	rows, err := dl.pool.Query(context.Background(), query, projectid)
	if err != nil {
		return []FetchProjectToolsResult{}, errMessQuery("FetchProjectTools", err)
	}
	var results []FetchProjectToolsResult
	pgx.ForEachRow(rows, []any{&id, &name}, func() error {
		result := FetchProjectToolsResult{
			Id:   resolveNil(id),
			Name: resolveNil(name)}
		results = append(results, result)
		return nil
	})
	return results, nil
}

// Query that fetches specific job information.
func (dl *Datalayer) FetchCareer(careerid string) (FetchCareerResult, error) {

	var id *int
	var title *string
	var companyname *string
	var description *string
	var startmonth *string
	var startyear *int
	var endmonth *string
	var endyear *int
	var current *bool
	query := getQueryFromPath("../sql/FetchCareer.sql")
	err := dl.pool.QueryRow(context.Background(), query, careerid).Scan(&id, &title, &companyname, &description, &startmonth, &startyear, &endmonth, &endyear, &current)
	if err != nil {
		return FetchCareerResult{}, errMessQuery("FetchCareer", err)
	}
	result := FetchCareerResult{
		Id:          resolveNil(id),
		Title:       resolveNil(title),
		Companyname: resolveNil(companyname),
		Description: resolveNil(description),
		Startmonth:  resolveNil(startmonth),
		Startyear:   resolveNil(startyear),
		Endmonth:    resolveNil(endmonth),
		Endyear:     resolveNil(endyear),
		Current:     resolveNil(current),
	}
	return result, nil
}

// Query that fetches education summaries.
func (dl *Datalayer) FetchEducation() ([]FetchEducationResult, error) {
	var name *string
	var link *string
	var title *string
	var major *string
	var gpa *string
	var startdate *int
	var enddate *int
	query := getQueryFromPath("../sql/FetchEducation.sql")

	rows, err := dl.pool.Query(context.Background(), query)
	if err != nil {
		return []FetchEducationResult{}, errMessQuery("FetchEducation", err)
	}
	var results []FetchEducationResult
	pgx.ForEachRow(rows, []any{&name, &link, &title, &major, &gpa, &startdate, &enddate}, func() error {
		result := FetchEducationResult{
			Name:      resolveNil(name),
			Link:      resolveNil(link),
			Title:     resolveNil(title),
			Major:     resolveNil(major),
			Gpa:       resolveNil(gpa),
			Startdate: resolveNil(startdate),
			Enddate:   resolveNil(enddate),
		}
		results = append(results, result)
		return nil
	})
	//fmt.Println(err.Error())
	return results, nil

}

// Query that fetches summaries for all projects.
func (dl *Datalayer) FetchAllProjects() ([]FetchProjectSummaryResult, error) {
	return dl.FetchProjectSummaries(-1)
}

// Query that fetches summaries for all careers.
func (dl *Datalayer) FetchAllCareers() ([]FetchCareerSummariesResult, error) {
	return dl.FetchCareerSummaries(-1)
}

// Query that fetches summaries for all education (degrees).
func (dl *Datalayer) FetchAllEducation() ([]FetchEducationSummariesResult, error) {
	return dl.FetchEducationSummaries(-1)
}

// Query that fetches biographical information.
func (dl *Datalayer) FetchBio() (FetchBiographyResult, error) {
	var firstname *string
	var lastname *string
	var description *string
	var email *string
	var linkedinlink *string
	var githublink *string
	var websitelink *string
	var portraitlink *string
	var resumelink *string
	query := getQueryFromPath("../sql/FetchBiography.sql")
	err := dl.pool.QueryRow(context.Background(), query).Scan(&firstname, &lastname, &description,
		&email, &linkedinlink, &githublink, &websitelink, &portraitlink, &resumelink)
	if err != nil {
		return FetchBiographyResult{}, errMessQuery("FetchBio", err)
	}
	result := FetchBiographyResult{
		Firstname:    resolveNil(firstname),
		Lastname:     resolveNil(lastname),
		Description:  resolveNil(description),
		Email:        resolveNil(email),
		Linkedinlink: resolveNil(linkedinlink),
		Githublink:   resolveNil(githublink),
		Websitelink:  resolveNil(websitelink),
		PortaitLink:  resolveNil(portraitlink),
		ResumeLink:   resolveNil(resumelink),
	}
	return result, nil
}

// Query that fetches a specific amount of summaries of projects.
func (dl *Datalayer) FetchProjectSummaries(limit int) ([]FetchProjectSummaryResult, error) {
	var id *int
	var name *string
	var thumbnaillink *string
	var description *string
	query := getQueryFromPath("../sql/FetchProjectSummaries.sql")
	var rows pgx.Rows
	var err error
	if limit <= -1 {
		rows, err = dl.pool.Query(context.Background(), query, nil)
	} else {
		rows, err = dl.pool.Query(context.Background(), query, limit)
	}
	if err != nil {
		return []FetchProjectSummaryResult{}, errMessQuery("FetchProjectSummaries", err)
	}
	var results []FetchProjectSummaryResult
	pgx.ForEachRow(rows, []any{&id, &name, &thumbnaillink, &description}, func() error {
		result := FetchProjectSummaryResult{
			Id:            resolveNil(id),
			Name:          resolveNil(name),
			Thumbnaillink: resolveNil(thumbnaillink),
			Description:   resolveNil(description),
		}
		results = append(results, result)
		return nil
	})

	return results, nil
}

// Query that fetches differently organized summaries of projects.
//
// It organizes by whether or not the project is tied to a career, with the ones not associated with one being first.
func (dl *Datalayer) FetchProjectSummariesExtra(limit int) ([]FetchProjectSummaryResult, error) {
	var id *int
	var name *string
	var thumbnaillink *string
	var description *string
	var endyear *int
	var iscareer *bool
	query := getQueryFromPath("../sql/FetchProjectSummariesCareerFilter.sql")
	var rows pgx.Rows
	var err error
	if limit <= -1 {
		rows, err = dl.pool.Query(context.Background(), query, nil)
	} else {
		rows, err = dl.pool.Query(context.Background(), query, limit)
	}
	if err != nil {
		return []FetchProjectSummaryResult{}, errMessQuery("FetchProjectSummariesExtra", err)
	}
	var results []FetchProjectSummaryResult
	pgx.ForEachRow(rows, []any{&id, &name, &thumbnaillink, &description, &endyear, &iscareer}, func() error {
		result := FetchProjectSummaryResult{
			Id:            resolveNil(id),
			Name:          resolveNil(name),
			Thumbnaillink: resolveNil(thumbnaillink),
			Description:   resolveNil(description),
			Endyear:       resolveNil(endyear),
			IsCareer:      resolveNil(iscareer),
		}
		results = append(results, result)
		return nil
	})

	return results, nil
}

// Query that fetches a specific number of summaries of careers.
func (dl *Datalayer) FetchCareerSummaries(limit int) ([]FetchCareerSummariesResult, error) {

	var id *string
	var title *string
	var name *string
	var description *string
	query := getQueryFromPath("../sql/FetchCareerSummaries.sql")
	var rows pgx.Rows
	var err error
	if limit <= -1 {
		rows, err = dl.pool.Query(context.Background(), query, nil)
	} else {
		rows, err = dl.pool.Query(context.Background(), query, limit)
	}
	if err != nil {
		return []FetchCareerSummariesResult{}, errMessQuery("FetchCareerSummaries", err)
	}
	var results []FetchCareerSummariesResult
	pgx.ForEachRow(rows, []any{&id, &title, &name, &description}, func() error {
		result := FetchCareerSummariesResult{
			Id:          resolveNil(id),
			Title:       resolveNil(title),
			Name:        resolveNil(name),
			Description: resolveNil(description),
		}

		results = append(results, result)
		return nil
	})

	return results, nil
}

// Query that fetches a specific summaries of degrees.
func (dl *Datalayer) FetchEducationSummaries(limit int) ([]FetchEducationSummariesResult, error) {

	var title *string
	var major *string
	var startdate *int
	var enddate *int
	query := getQueryFromPath("../sql/FetchEducationSummaries.sql")
	var rows pgx.Rows
	var err error
	if limit <= -1 {
		rows, err = dl.pool.Query(context.Background(), query, nil)
	} else {
		rows, err = dl.pool.Query(context.Background(), query, limit)
	}
	if err != nil {
		return []FetchEducationSummariesResult{}, errMessQuery("FetchEducationSummaries", err)
	}
	var results []FetchEducationSummariesResult
	pgx.ForEachRow(rows, []any{&title, &major, &startdate, &enddate}, func() error {
		result := FetchEducationSummariesResult{
			Title:     resolveNil(title),
			Major:     resolveNil(major),
			StartDate: resolveNil(startdate),
			EndDate:   resolveNil(enddate),
		}

		results = append(results, result)
		return nil
	})

	return results, nil
}
