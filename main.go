package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/colorrr34/status-detector/config"
	"github.com/colorrr34/status-detector/internal/database"
	_ "github.com/lib/pq"
	"github.com/urfave/cli/v3"
)

type state struct{
	db *database.Queries
	cfg *config.Config
}

func main() {
	cfg := config.Read();
	db,err := sql.Open("postgres",cfg.DbUrl)
	if err != nil{
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	cfgState := state{
		cfg: &cfg,
		db: dbQueries,
	}
    cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name: "seed",
				Usage: "Seed the database with default sites",
				Action: func(ctx context.Context, _ *cli.Command) error {
					return seedHandler(ctx,&cfgState)
				},
			},
			{
				Name: "list",
				Usage: "List all sites",
				Action: func(ctx context.Context, c *cli.Command) error {
					return listHandler(ctx,c,&cfgState);
				},
			},
			{
				Name: "add",
				Usage: "add a site to the list, 1st arg for the name and 2nd for the url, flag d for description",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "description",
						Aliases: []string{"d"},
						Usage: "description for the site",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return addHandler(ctx,c,&cfgState)
				},
			},
			{
				Name: "ping",
				Usage: "Ping all sites and list the responses",
				Action: func(ctx context.Context, c *cli.Command) error {
					return pingHandler(ctx,c,&cfgState)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(),os.Args); err != nil{
		log.Fatal(err)
	}
}