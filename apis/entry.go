package apis

import (
	"context"
	"net/http"

	"github.com/nanoteck137/nosepass/core"
	"github.com/nanoteck137/nosepass/database"
	"github.com/nanoteck137/nosepass/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
)

type EntryInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Entry struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetEntries struct {
	Entries []EntryInfo `json:"entries"`
}

type CreateEntry struct {
	Id string `json:"id"`
}

type CreateEntryBody struct {
	Name string `json:"name"`
}

func (b *CreateEntryBody) Transform() {
	b.Name = transform.String(b.Name)
}

func (b CreateEntryBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type EditEntryBody struct {
	Name *string `json:"name"`
}

func (b *EditEntryBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
}

func (b EditEntryBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
	)
}

func InstallEntryHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetEntries",
			Method:       http.MethodGet,
			Path:         "/entries",
			ResponseType: GetEntries{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				ctx := context.TODO()

				entries, err := app.DB().GetAllEntries(ctx)
				if err != nil {
					return nil, err
				}

				res := GetEntries{
					Entries: make([]EntryInfo, len(entries)),
				}

				for i, e := range entries {
					res.Entries[i] = EntryInfo{
						Id:   e.Id,
						Name: e.Name,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetEntryById",
			Method:       http.MethodGet,
			Path:         "/entries/:id",
			ResponseType: Entry{},
			BodyType:     nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				entry, err := app.DB().GetEntryById(ctx, id)
				if err != nil {
					// TODO(patrik): Handle not found error
					return nil, err
				}

				return Entry{
					Id:   entry.Id,
					Name: entry.Name,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateEntry",
			Method:       http.MethodPost,
			Path:         "/entries",
			ResponseType: CreateEntry{},
			BodyType:     CreateEntryBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[CreateEntryBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				entry, err := app.DB().CreateEntry(ctx, database.CreateEntryParams{
					Name: body.Name,
				})
				if err != nil {
					return nil, err
				}

				return CreateEntry{
					Id: entry.Id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "EditEntry",
			Method:       http.MethodPatch,
			Path:         "/entries/:id",
			ResponseType: nil,
			BodyType:     EditEntryBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[EditEntryBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				entry, err := app.DB().GetEntryById(ctx, id)
				if err != nil {
					// TODO(patrik): Handle not found error
					return nil, err
				}

				changes := database.EntryChanges{}

				if body.Name != nil {
					changes.Name = types.Change[string]{
						Value:   *body.Name,
						Changed: *body.Name != entry.Name,
					}
				}

				err = app.DB().UpdateEntry(ctx, id, changes)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "DeleteEntry",
			Method:       http.MethodDelete,
			Path:         "/entries/:id",
			ResponseType: nil,
			BodyType:     nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				err := app.DB().DeleteEntry(ctx, id)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
