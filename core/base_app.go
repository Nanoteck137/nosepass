package core

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/kr/pretty"
	"github.com/nanoteck137/nosepass/config"
	"github.com/nanoteck137/nosepass/core/log"
	"github.com/nanoteck137/nosepass/database"
	"github.com/nanoteck137/nosepass/library"
	"github.com/nanoteck137/nosepass/types"
	"gopkg.in/vansante/go-ffprobe.v2"
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

	// dirs := []string{
	// 	workDir.Artists(),
	// 	workDir.Albums(),
	// 	workDir.Tracks(),
	// 	workDir.Trash(),
	// }
	//
	// for _, dir := range dirs {
	// 	err = os.Mkdir(dir, 0755)
	// 	if err != nil && !os.IsExist(err) {
	// 		return err
	// 	}
	// }

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

	lib, err := library.ReadFromDisk("/Volumes/media2/anime")
	if err != nil {
		return err
	}

	pretty.Println(lib)

	ctx := context.TODO()

	for _, serie := range lib.Series {
		for _, collection := range serie.Collections {
			for _, media := range collection.MediaItems {
				fmt.Printf("media.Name: %v/%v/%v\n", serie.Name, collection.Name, media.Name)

				_, err := app.db.GetMediaByPath(ctx, media.Path)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						fmt.Println("Create new media item")

						probeData, err := ffprobe.ProbeURL(ctx, media.Path)
						if err != nil {
							return err
						}

						var chapters []database.MediaChapter

						for _, chapter := range probeData.Chapters {
							chapters = append(chapters, database.MediaChapter{
								Title: chapter.Title(),
								Start: chapter.StartTimeSeconds,
								End:   chapter.EndTimeSeconds,
							})
						}

						var subtitles []database.MediaSubtitle
						var attachments []database.MediaAttachment
						var videoTracks []database.MediaVideoTrack
						var audioTracks []database.MediaAudioTrack

						for i, s := range probeData.Streams {
							switch s.CodecType {
							case "video":
								videoTracks = append(videoTracks, database.MediaVideoTrack{
									Index:      i,
									VideoIndex: len(videoTracks),
								})
							case "audio":
								audioTracks = append(audioTracks, database.MediaAudioTrack{
									Index:      i,
									AudioIndex: len(audioTracks),
									Language:   s.Tags.Language,
								})
							case "subtitle":
								subtitles = append(subtitles, database.MediaSubtitle{
									Index:         i,
									SubtitleIndex: len(subtitles),
								})
							case "attachment":
								attachments = append(attachments, database.MediaAttachment{
									Index:           i,
									AttachmentIndex: len(attachments) + 1,
								})
							}
						}

						res, err := app.db.CreateMedia(ctx, database.CreateMediaParams{
							Path:             media.Path,
							FileModifiedTime: media.ModTime.UnixMilli(),
							Chapters:         chapters,
							Subtitles:        subtitles,
							Attachments:      attachments,
							VideoTracks:      videoTracks,
							AudioTracks:      audioTracks,
						})
						if err != nil {
							return err
						}

						pretty.Println(res)

						media, err := app.db.GetMediaById(ctx, res.Id)
						if err != nil {
							return err
						}

						for _, audio := range media.AudioTracks.GetOrEmpty() {
							pretty.Println(audio)

							name := "Unknown"
							language := "unknown"

							switch audio.Language {
							case "eng":
								name = "English"
								language = "en"
							case "jpn":
								name = "Japanese"
								language = "jp"
							}

							res, err := app.db.CreateMediaVariant(ctx, database.CreateMediaVariantParams{
								MediaId:    res.Id,
								Name:       name,
								Language:   language,
								VideoTrack: 0,
								AudioTrack: audio.AudioIndex,
								// Subtitle:   sql.NullString{},
							})
							if err != nil {
								return err
							}

							pretty.Println(res)
						}
					} else {
						return err
					}
				}

				fmt.Println("Update media item")
				fmt.Println()
			}
		}
	}

	// media, err := app.db.GetAllMedia(ctx)
	// if err != nil {
	// 	return err
	// }
	//
	// for _, m := range media {
	// 	pretty.Println(m.Chapters.GetOrEmpty())
	// }

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
	}

	return nil
}

func NewBaseApp(config *config.Config) *BaseApp {
	return &BaseApp{
		config: config,
	}
}
