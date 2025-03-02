package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"

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

	lib, err := library.ReadFromDisk("/Volumes/media2/anime")
	if err != nil {
		return err
	}

	pretty.Println(lib)

	ctx := context.TODO()

	wg := sync.WaitGroup{}
	count := 0

	for _, serie := range lib.Series {
		for _, collection := range serie.Collections {
			for _, mediaItem := range collection.MediaItems {
				if count > 5 {
					wg.Wait()
					count = 0
				}

				wg.Add(1)
				count++

				mediaItem := mediaItem
				go func() {
					defer wg.Done()

					fmt.Printf("media.Name: %v/%v/%v\n", serie.Name, collection.Name, mediaItem.Name)

					media, err := app.db.GetMediaByPath(ctx, mediaItem.Path)
					if err != nil {
						if errors.Is(err, database.ErrItemNotFound) {
							fmt.Println("Create new media item")

							probeData, err := ffprobe.ProbeURL(ctx, mediaItem.Path)
							if err != nil {
								log.Fatal("Failed", "err", err)
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
										// TODO(patrik): Move this
										Duration: probeData.Format.DurationSeconds,
									})
								case "audio":
									audioTracks = append(audioTracks, database.MediaAudioTrack{
										Index:      i,
										AudioIndex: len(audioTracks),
										Language:   s.Tags.Language,
									})
								case "subtitle":
									typ := s.CodecName

									switch typ {
									case "ass":
										ext := ""
										switch typ {
										case "ass":
											ext = ".ass"
										default:
											log.Error("Missing subtitle type", "type", typ)
											continue
										}

										index := len(subtitles)
										subtitles = append(subtitles, database.MediaSubtitle{
											Index:         i,
											SubtitleIndex: index,
											Type:          typ,
											Title:         s.Tags.Title,
											Language:      s.Tags.Language,
											IsDefault:     s.Disposition.Default == 1,
											Filename:      strconv.Itoa(index) + ext,
										})
									default:
										log.Warn("Unsupported subtitle type (skipping)", "path", mediaItem.Path, "type", typ)
										continue
									}

								case "attachment":
									attachments = append(attachments, database.MediaAttachment{
										Index:           i,
										AttachmentIndex: len(attachments) + 1,
									})
								}
							}

							modTime := mediaItem.ModTime
							res, err := app.db.CreateMedia(ctx, database.CreateMediaParams{
								Path:             mediaItem.Path,
								FileModifiedTime: 0,
								Chapters:         chapters,
								Subtitles:        subtitles,
								Attachments:      attachments,
								VideoTracks:      videoTracks,
								AudioTracks:      audioTracks,
							})
							if err != nil {
								log.Fatal("Failed", "err", err)
							}

							media, err := app.db.GetMediaById(ctx, res.Id)
							if err != nil {
								log.Fatal("Failed", "err", err)
							}

							mediaDir := workDir.MediaIdDir(media.Id)

							dirs := []string{
								mediaDir.String(),
								mediaDir.Subtitles(),
								mediaDir.Attachments(),
							}

							for _, dir := range dirs {
								err = os.Mkdir(dir, 0755)
								if err != nil && !os.IsExist(err) {
									log.Fatal("Failed", "err", err)
								}
							}

							subtitles = media.Subtitles.GetOrEmpty()
							if len(subtitles) > 0 {
								args := []string{
									"-nostats",
									"-hide_banner",
									"-loglevel", "error",
									"-i", media.Path,
								}

								// ffmpeg -i Movie.mkv -map 0:s:0 subs.srt
								for _, subtitle := range subtitles {
									ext := ""
									switch subtitle.Type {
									case "ass":
										ext = ".ass"
									default:
										log.Warn("Missing subtitle type", "type", subtitle.Type)
										continue
									}

									args = append(args,
										"-map", fmt.Sprintf("0:s:%d", subtitle.SubtitleIndex),
										"-c", "copy",
										path.Join(mediaDir.Subtitles(), strconv.Itoa(subtitle.SubtitleIndex)+ext),
									)
								}

								pretty.Println(args)

								cmd := exec.Command("ffmpeg", args...)
								cmd.Stderr = os.Stderr

								err = cmd.Run()
								if err != nil {
									log.Fatal("Failed", "err", err)
								}
							}

							for _, audio := range media.AudioTracks.GetOrEmpty() {
								switch audio.Language {
								case "eng":
									_, err := app.db.CreateMediaVariant(ctx, database.CreateMediaVariantParams{
										MediaId:    media.Id,
										Name:       "English",
										Language:   "en",
										VideoTrack: 0,
										AudioTrack: audio.AudioIndex,
										// Subtitle:   sql.NullString{},
									})
									if err != nil {
										log.Fatal("Failed", "err", err)
									}
								case "jpn":
									_, err := app.db.CreateMediaVariant(ctx, database.CreateMediaVariantParams{
										MediaId:    media.Id,
										Name:       "Japanese",
										Language:   "jp",
										VideoTrack: 0,
										AudioTrack: audio.AudioIndex,
										// Subtitle:   sql.NullString{},
									})
									if err != nil {
										log.Fatal("Failed", "err", err)
									}
								}
							}

							err = app.db.UpdateMedia(ctx, media.Id, database.MediaChanges{
								FileModifiedTime: types.Change[int64]{
									Value:   modTime.UnixMilli(),
									Changed: true,
								},
							})
							if err != nil {
								log.Fatal("Failed", "err", err)
							}
						} else {
							log.Fatal("Failed", "err", err)
						}
					} else {
						fmt.Println("Update media item")
						fmt.Println()

						if media.FileModifiedTime < mediaItem.ModTime.UnixMilli() {
							fmt.Printf("Needs Updating: %v\n", media.Path)

							mediaDir := workDir.MediaIdDir(media.Id)

							err := os.RemoveAll(mediaDir.String())
							if err != nil && !os.IsNotExist(err) {
								log.Fatal("Failed", "err", err)
							}

							dirs := []string{
								mediaDir.String(),
								mediaDir.Subtitles(),
								mediaDir.Attachments(),
							}

							for _, dir := range dirs {
								err = os.Mkdir(dir, 0755)
								if err != nil && !os.IsExist(err) {
									log.Fatal("Failed", "err", err)
								}
							}

							subtitles := media.Subtitles.GetOrEmpty()

							if len(subtitles) > 0 {
								args := []string{
									"-nostats",
									"-hide_banner",
									"-loglevel", "error",
									"-i", media.Path,
								}

								// ffmpeg -i Movie.mkv -map 0:s:0 subs.srt
								for _, subtitle := range subtitles {
									ext := ""
									switch subtitle.Type {
									case "ass":
										ext = ".ass"
									default:
										log.Warn("Missing subtitle type", "type", subtitle.Type)
										continue
									}

									args = append(args,
										"-map", fmt.Sprintf("0:s:%d", subtitle.SubtitleIndex),
										"-c", "copy",
										path.Join(mediaDir.Subtitles(), strconv.Itoa(subtitle.SubtitleIndex)+ext),
									)
								}

								pretty.Println(args)

								cmd := exec.Command("ffmpeg", args...)
								cmd.Stderr = os.Stderr

								err = cmd.Run()
								if err != nil {
									log.Fatal("Failed", "err", err)
								}
							}

							err = app.db.UpdateMedia(ctx, media.Id, database.MediaChanges{
								FileModifiedTime: types.Change[int64]{
									Value:   mediaItem.ModTime.UnixMilli(),
									Changed: true,
								},
							})
							if err != nil {
								log.Fatal("Failed", "err", err)
							}
						}
					}
				}()
			}
		}
	}

	wg.Wait()

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
