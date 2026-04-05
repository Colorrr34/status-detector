package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/colorrr34/status-detector/internal/database"
	"github.com/urfave/cli/v3"
)

func seedHandler(ctx context.Context,s *state)error{
	sites := []database.CreateSiteParams{
		database.CreateSiteParams{
			Name: "GitHub",
			Url: "https://github.com/",
			Description: sql.NullString{
				Valid: true,
				String: "Developer platform",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		database.CreateSiteParams{
			Name: "Youtube",
			Url: "https://www.youtube.com/",
			Description: sql.NullString{
				Valid: true,
				String: "Streaming platform",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	dbSites,err := s.db.GetSites(ctx)
	if err != nil{
		return err
	}else if dbSites!= nil{
		fmt.Println("Database already seeded")
		return nil
	}
	for _,site := range sites{
		if site,err:= s.db.CreateSite(ctx,site);err!=nil{
			fmt.Printf("%v added",site.Name)
		}
	}
	return nil
}

func listHandler(ctx context.Context, _ *cli.Command, s *state)error{
	type siteResponse struct{
		name string
		url string
		description string
	}
	sites,err := s.db.GetSites(ctx)
	if err !=nil{
		return err
	}
	for _,site := range sites{
		fmt.Println(siteResponse{
			name: site.Name,
			url: site.Url,
			description: site.Description.String,
		})
	}
	
	return nil
}

func addHandler(ctx context.Context, c *cli.Command, s *state)error{
	type response struct{
		message string
		name string
		url string
		description string
	}
	name := c.Args().Get(0)
	url := c.Args().Get(1)
	var description sql.NullString
	if name == "" || url == ""{
		return errors.New("missing arguments")
	}
	
	if c.String("description")!=""{
		description = sql.NullString{
			Valid: true,
			String: c.String("description"),
		}
	}
	site,err := s.db.CreateSite(ctx,database.CreateSiteParams{
		Name: name,
		Url: url,
		Description: description,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err!=nil{
		return err
	}

	fmt.Println(response{
		message: "site added",
		name: site.Name,
		url: site.Url,
		description: site.Description.String,
	})

	return nil
}

func pingHandler(ctx context.Context, c *cli.Command, s *state)error{
	type pingResponse struct{
		name string
		status int
	}
	sites,err := s.db.GetSites(ctx)

	if err != nil{
		return err
	}

	ch := make(chan Result,len(sites))
	var wg sync.WaitGroup
	for _,site := range sites{
		wg.Add(1)
		go checkSite(site, ch,&wg)
	}

	go func(){
		wg.Wait()
		close(ch)
	}()

	for result := range ch{
		if result.Err!=nil{
			fmt.Printf("%s - Down (%v)\n",result.Name,result.Err)
		}else{
			fmt.Printf("%s - Status: %d",result.Name,result.Status)
		}
	}
	return nil	
}