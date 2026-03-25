package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/colorrr34/status-detector/internal/database"
	"github.com/urfave/cli/v3"
)

func listHandler(ctx context.Context, _ *cli.Command, s *state)error{
	sites,err := s.db.GetSites(ctx)
	if err !=nil{
		return err
	}
	fmt.Println(sites)
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