package apis

import (
	"context"
	"net/http"
	"sort"

	"github.com/maruel/natural"
	"github.com/nanoteck137/nosepass/core"
	"github.com/nanoteck137/pyrin"
)

type Collection struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type FullCollection struct {
	Collection

	Media []FullMedia `json:"media"`
}

type GetCollections struct {
	Collections []Collection `json:"collections"`
}

func InstallCollectionHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetCollections",
			Method:       http.MethodGet,
			Path:         "/collections",
			ResponseType: GetCollections{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				collections, err := app.DB().GetAllCollections(ctx)
				if err != nil {
					return nil, err
				}

				res := GetCollections{
					Collections: make([]Collection, len(collections)),
				}

				for i, col := range collections {
					res.Collections[i] = Collection{
						Id:   col.Id,
						Name: col.Name,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetCollectionById",
			Method:       http.MethodGet,
			Path:         "/collections/:id",
			ResponseType: FullCollection{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				collection, err := app.DB().GetCollectionById(ctx, id)
				if err != nil {
					// TODO(patrik): Handle error
					return nil, err
				}

				media, err := app.DB().GetMediaByCollectionId(ctx, collection.Id)
				if err != nil {
					return nil, err
				}

				res := FullCollection{
					Collection: Collection{
						Id:   collection.Id,
						Name: collection.Name,
					},
					Media: make([]FullMedia, len(media)),
				}

				for i, m := range media {

					variants := m.Variants.GetOrEmpty()

					me := FullMedia{
						Id:       m.Id,
						Path:     m.Path,
						Variants: make([]MediaVariant, len(variants)),
					}

					for i, variant := range variants {
						me.Variants[i] = MediaVariant{
							Id:       variant.Id,
							Name:     variant.Name,
							Language: variant.Language,
						}
					}

					res.Media[i] = me
				}

				sort.SliceStable(res.Media, func(i, j int) bool {
					return natural.Less(res.Media[i].Path, res.Media[j].Path)
				})

				return res, nil
			},
		},
	)
}
