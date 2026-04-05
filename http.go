package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/colorrr34/status-detector/internal/database"
)

type Result struct{
	Name string
	Url string
	Status int
	Err error
}

func checkSite(site database.Site, ch chan<- Result, wg *sync.WaitGroup){
	defer wg.Done()

	client := &http.Client{Timeout: 10*time.Second}
	resp, err := client.Head(site.Url)
	if err != nil{
		ch <- Result{
			Name: site.Name,
			Url: site.Url,
			Err: err,
		}
	}
	defer resp.Body.Close()

	ch <- Result{
		Name: site.Name,
		Url: site.Url,
		Status: resp.StatusCode,
	}
}