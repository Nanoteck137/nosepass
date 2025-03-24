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

type MediaSubtitle struct {
	Index     int    `json:"index"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Language  string `json:"language"`
	// TODO(patrik): Add to audio track
	IsDefault bool   `json:"isDefault"`
}

type MediaAudioTrack struct {
	Index    int    `json:"index"`
	Language string `json:"language"`
}

type MediaVariant struct {
	AudioTrack int    `json:"audio_track"`
	Subtitle   *int64 `json:"subtitle"`
}

type FullMedia struct {
	Id   string `json:"id"`
	Path string `json:"path"`

	AudioTracks []MediaAudioTrack `json:"audioTracks"`
	Subtitles   []MediaSubtitle   `json:"subtitles"`

	SubVariant *MediaVariant `json:"subVariant"`
	DubVariant *MediaVariant `json:"dubVariant"`
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

				audioTracks := media.AudioTracks.GetOrEmpty()
				subtitles := media.Subtitles.GetOrEmpty()

				res := FullMedia{
					Id:          media.Id,
					Path:        media.Path,
					AudioTracks: make([]MediaAudioTrack, len(audioTracks)),
					Subtitles:   make([]MediaSubtitle, len(subtitles)),
				}

				for i, audio := range audioTracks {
					res.AudioTracks[i] = MediaAudioTrack{
						Index:    audio.AudioIndex,
						Language: audio.Language,
					}
				}

				for i, subtitle := range subtitles {
					res.Subtitles[i] = MediaSubtitle{
						Index:     subtitle.SubtitleIndex,
						Type:      subtitle.Type,
						Title:     subtitle.Title,
						Language:  subtitle.Language,
						IsDefault: subtitle.IsDefault,
					}
				}

				return res, nil
			},
		},
	)
}
