package apis

import (
	"context"
	"net/http"
	"strings"

	"github.com/nanoteck137/nosepass/core"
	"github.com/nanoteck137/pyrin"
)

type Media struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

type MediaAudio struct {
	Index    int    `json:"index"`
	Language string `json:"language"`
}

type MediaSubtitle struct {
	Index    int    `json:"index"`
	Language string `json:"language"`
	Title    string `json:"title"`
}

type FullMedia struct {
	Id   string `json:"id"`
	Path string `json:"path"`

	AudioTracks []MediaAudio    `json:"audioTracks"`
	Subtitles   []MediaSubtitle `json:"subtitles"`
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

				audioTracks := media.AudioTracks.GetOrEmpty()
				subtitles := media.Subtitles.GetOrEmpty()

				res := FullMedia{
					Id:          media.Id,
					Path:        media.Path,
					AudioTracks: make([]MediaAudio, len(audioTracks)),
					Subtitles:   make([]MediaSubtitle, len(subtitles)),
				}

				for i, track := range audioTracks {
					res.AudioTracks[i] = MediaAudio{
						Index:    track.AudioIndex,
						Language: track.Language,
					}
				}

				for i, subtitle := range subtitles {
					title := strings.TrimSpace(subtitle.Title)
					if title == "" {
						title = "Unknown"
						switch subtitle.Language {
						case "eng":
							title = "English"
						}
					}

					language := "unknown"
					switch subtitle.Language {
					case "eng":
						language = "en"
					}

					res.Subtitles[i] = MediaSubtitle{
						Index:    subtitle.SubtitleIndex,
						Language: language,
						Title:    title,
					}
				}

				return res, nil
			},
		},
	)
}
