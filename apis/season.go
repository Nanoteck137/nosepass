package apis

import (
	"context"
	"net/http"

	"github.com/kr/pretty"
	"github.com/nanoteck137/nosepass/core"
	"github.com/nanoteck137/pyrin"
)

type SeasonEpisodeVariant struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Language string `json:"language"`
}

type SeasonEpisode struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Variants []SeasonEpisodeVariant `json:"variants"`
}

type GetSeasonById struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Episodes []SeasonEpisode `json:"episodes"`
}

func InstallSeasonHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetSeasonById",
			Method:       http.MethodGet,
			Path:         "/seasons/:id",
			ResponseType: GetSeasonById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()
				season, err := app.DB().GetSeasonById(ctx, id)
				if err != nil {
					// TODO(patrik): Handle error
					return nil, err
				}

				episodes, err := app.DB().GetEpisodes(ctx, season.Id)
				if err != nil {
					return nil, err
				}

				pretty.Println(episodes)

				res := GetSeasonById{
					Id:       season.Id,
					Name:     season.Name,
					Episodes: make([]SeasonEpisode, len(episodes)),
				}

				for i, episode := range episodes {
					variants := episode.Variants.GetOrEmpty()

					e := SeasonEpisode{
						Id:       episode.Id,
						Name:     episode.Name,
						Variants: make([]SeasonEpisodeVariant, len(variants)),
					}

					for i, variant := range variants {
						e.Variants[i] = SeasonEpisodeVariant{
							Id:       variant.Id,
							Name:     variant.Name,
							Language: variant.Language,
						}
					}

					res.Episodes[i] = e
				}

				return res, nil
			},
		},
	)
}
