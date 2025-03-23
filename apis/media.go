package apis

import (
	"context"
	"net/http"
	"sort"

	"github.com/maruel/natural"
	"github.com/nanoteck137/nosepass/core"
	"github.com/nanoteck137/pyrin"
)

type Media struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

type MediaVariant struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Language string `json:"language"`
}

type FullMedia struct {
	Id   string `json:"id"`
	Path string `json:"path"`

	Variants []MediaVariant `json:"variants"`
}

type GetMedia struct {
	Media []Media `json:"media"`
}

func InstallMediaHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetMedia",
			Method:       http.MethodGet,
			Path:         "/media",
			ResponseType: GetMedia{},
			Errors:       []pyrin.ErrorType{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				media, err := app.DB().GetAllMedia(ctx)
				if err != nil {
					return nil, err
				}

				res := GetMedia{
					Media: make([]Media, len(media)),
				}

				for i, media := range media {
					res.Media[i] = Media{
						Id:   media.Id,
						Path: media.Path,
					}
				}

				sort.SliceStable(res.Media, func(i, j int) bool {
					return natural.Less(res.Media[i].Path, res.Media[j].Path)
				})

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetMediaById",
			Method:       http.MethodGet,
			Path:         "/media/:id",
			ResponseType: FullMedia{},
			Errors:       []pyrin.ErrorType{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				media, err := app.DB().GetMediaById(ctx, id)
				if err != nil {
					return nil, err
				}

				variants := media.Variants.GetOrEmpty()

				res := FullMedia{
					Id:          media.Id,
					Path:        media.Path,
					Variants: make([]MediaVariant, len(variants)),
				}

				for i, variant := range variants {
					res.Variants[i] = MediaVariant{
						Id:       variant.Id,
						Name:     variant.Name,
						Language: variant.Language,
					}
				}

				return res, nil
			},
		},
	)
}
