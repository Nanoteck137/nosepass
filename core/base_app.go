package core

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/kr/pretty"
	"github.com/nanoteck137/nosepass/config"
	"github.com/nanoteck137/nosepass/core/log"
	"github.com/nanoteck137/nosepass/database"
	"github.com/nanoteck137/nosepass/types"
)

var _ App = (*BaseApp)(nil)

type BaseApp struct {
	db     *database.Database
	config *config.Config
}

func (app *BaseApp) DB() *database.Database {
	return app.db
}

func (app *BaseApp) Config() *config.Config {
	return app.config
}

func (app *BaseApp) WorkDir() types.WorkDir {
	return app.config.WorkDir()
}

func (app *BaseApp) Bootstrap() error {
	var err error

	workDir := app.config.WorkDir()

	dirs := []string{
		workDir.MediaDir(),
	}

	for _, dir := range dirs {
		err = os.Mkdir(dir, 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	app.db, err = database.Open(workDir)
	if err != nil {
		return err
	}

	if app.config.RunMigrations {
		err = database.RunMigrateUp(app.db)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(workDir.SetupFile())
	if errors.Is(err, os.ErrNotExist) && app.config.Username != "" {
		log.Info("Server not setup, creating the initial user")

		ctx := context.Background()

		_, err := app.db.CreateUser(ctx, database.CreateUserParams{
			Username: app.config.Username,
			Password: app.config.InitialPassword,
			Role:     types.RoleSuperUser,
		})
		if err != nil {
			return err
		}

		f, err := os.Create(workDir.SetupFile())
		if err != nil {
			return err
		}
		f.Close()

		{
			serie, err := app.db.CreateSerie(ctx, database.CreateSerieParams{
				Name: "Solo Leveling",
			})
			if err != nil {
				return err
			}

			pretty.Println(serie)

			season, err := app.db.CreateSeason(ctx, database.CreateSeasonParams{
				SerieId: serie.Id,
				Name:    "Season 1",
				Number:  1,
				Type:    "season",
			})
			if err != nil {
				return err
			}

			pretty.Println(season)

			for i := 0; i < 12; i++ {
				episode, err := app.db.CreateEpisode(ctx, database.CreateEpisodeParams{
					SeasonId: season.Id,
					Name:     fmt.Sprintf("Episode %d", i + 1),
					SeasonNumber: sql.NullInt64{
						Int64: int64(i + 1),
						Valid: true,
					},
					SerieNumber: sql.NullInt64{
						Int64: int64(i + 1),
						Valid: true,
					},
				})
				if err != nil {
					return err
				}

				pretty.Println(episode)
			}
		}
	}

	return nil
}

func NewBaseApp(config *config.Config) *BaseApp {
	return &BaseApp{
		config: config,
	}
}
