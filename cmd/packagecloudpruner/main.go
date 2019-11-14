package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	resty "github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

const (
	baseURL    = "https://packagecloud.io"
	user       = "sensu"
	repository = "ci-builds"
	maxRetries = 3
	maxAge     = 30 // in days
)

type PackagecloudPackage struct {
	Name          string    `json:"name"`
	DistroVersion string    `json:"distro_version"`
	CreatedAt     time.Time `json:"created_at"`
	Version       string    `json:"version"`
	Release       string    `json:"release"`
	Epoch         int       `json:"epoch"`
	DestroyURL    string    `json:"destroy_url"`
}

type ResponseError struct {
	Status int
	Error  string
}

type PaginationRequest struct {
	client         *resty.Client
	currentPage    uint
	pagesRemaining uint
	maxRetries     uint
	errCount       uint
}

func (p *PaginationRequest) HasPagesRemaining() bool {
	return p.pagesRemaining > 0
}

func (p *PaginationRequest) CanRetry(err error) bool {
	log.Errorf("%v", err)

	retriesRemaining := p.maxRetries - p.errCount

	if retriesRemaining != 0 {
		log.Errorf("retrying (%d retries left)", retriesRemaining)
		p.errCount = p.errCount + 1
		return true
	}

	log.Errorf("max retries of %d reached, quitting", p.maxRetries)

	return false
}

func (p *PaginationRequest) RequestCurrentPage() error {
	log.WithFields(log.Fields{
		"currentPage":    p.currentPage,
		"pagesRemaining": p.pagesRemaining,
		"maxRetries":     p.maxRetries,
		"errCount":       p.errCount,
	}).Info("fetching page")

	// request the current page
	resource := fmt.Sprintf("/api/v1/repos/%s/%s/packages.json", user, repository)
	queryParams := map[string]string{
		"page": strconv.Itoa(int(p.currentPage)),
	}
	resp, err := p.client.R().
		SetQueryParams(queryParams).
		SetResult([]PackagecloudPackage{}).
		Get(resource)
	if err != nil {
		return err
	}
	if resp.IsError() {
		err := *resp.Error().(*ResponseError)
		return fmt.Errorf("failed to get package list (url: %s, status: %d, error: %s)", resp.Request.URL, err.Status, err.Error)
	}

	// loop through the packages and remove them if the time since creation
	// is 30 days or more
	for _, pkg := range *resp.Result().(*[]PackagecloudPackage) {
		age := int(time.Since(pkg.CreatedAt).Hours() / 24)
		if age >= maxAge {
			log.WithFields(log.Fields{
				"name":           pkg.Name,
				"distro_version": pkg.DistroVersion,
				"version":        pkg.Version,
				"release":        pkg.Release,
				"epoch":          pkg.Epoch,
				"age":            age,
				"age_max":        maxAge,
			}).Info("removing package as it exceeds max age")

			err := removePackage(p.client, pkg.DestroyURL)
			if err != nil {
				return err
			}
		}
	}

	// retrieve the pagination headers
	headers := resp.Header()
	perPageStr, ok := headers["Per-Page"]
	if !ok {
		return errors.New("could not read header \"Per-Page\"")
	}
	totalPackagesStr, ok := headers["Total"]
	if !ok {
		return errors.New("could not read header \"Total\"")
	}

	// convert headers to ints
	perPage, err := strconv.Atoi(perPageStr[0])
	if err != nil {
		return err
	}
	totalPackages, err := strconv.Atoi(totalPackagesStr[0])
	if err != nil {
		return err
	}

	// success, reset the errCount
	p.errCount = 0

	// calculate the pages remaining
	p.currentPage = p.currentPage + 1
	totalPages := math.Ceil(float64(totalPackages) / float64(perPage))
	p.pagesRemaining = uint(totalPages) - p.currentPage

	return nil
}

func removePackage(client *resty.Client, destroyURL string) error {
	resp, err := client.R().Delete(destroyURL)
	if err != nil {
		return err
	}
	if resp.IsError() {
		err := *resp.Error().(*ResponseError)
		return fmt.Errorf("error removing package (url: %s, status: %d, error: %s)", resp.Request.URL, err.Status, err.Error)
	}
	return nil
}

func main() {
	token := os.Getenv("PACKAGECLOUD_TOKEN")
	if token == "" {
		log.Fatal("the environment variable PACKAGECLOUD_TOKEN must be set")
	}

	globalHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	client := resty.New()
	client.SetBasicAuth(token, "")
	client.SetHostURL(baseURL)
	client.SetHeaders(globalHeaders)
	client.SetError(&ResponseError{})

	p := PaginationRequest{
		client:         client,
		currentPage:    0,
		pagesRemaining: 1,
		errCount:       0,
		maxRetries:     maxRetries,
	}

	for p.HasPagesRemaining() {
		if err := p.RequestCurrentPage(); err != nil {
			if p.CanRetry(err) {
				continue
			}
			break
		}
	}
}
