package apis

import (
	"context"
	"net/http"

	"github.com/nanoteck137/nosepass/core"
	"github.com/nanoteck137/pyrin"
)

type Serie struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetSeries struct {
	Series []Serie `json:"series"`
}

type SerieSeason struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetSerieById struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Seasons []SerieSeason `json:"seasons"`
}

func InstallSerieHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetSeries",
			Method:       http.MethodGet,
			Path:         "/series",
			ResponseType: GetSeries{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				series, err := app.DB().GetAllSeries(ctx)
				if err != nil {
					return nil, err
				}

				res := GetSeries{
					Series: make([]Serie, len(series)),
				}

				for i, serie := range series {
					res.Series[i] = Serie{
						Id:   serie.Id,
						Name: serie.Name,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetSerieById",
			Method:       http.MethodGet,
			Path:         "/series/:id",
			ResponseType: GetSerieById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				serie, err := app.DB().GetSerieById(ctx, id)
				if err != nil {
					// TODO(patrik): Handle error
					return nil, err
				}

				seasons, err := app.DB().GetSeasonBySerie(ctx, serie.Id)
				if err != nil {
					return nil, err
				}

				res := GetSerieById{
					Id:      serie.Id,
					Name:    serie.Name,
					Seasons: make([]SerieSeason, len(seasons)),
				}

				for i, season := range seasons {
					res.Seasons[i] = SerieSeason{
						Id:   season.Id,
						Name: season.Name,
					}
				}

				return res, nil
			},
		},
	)
}
