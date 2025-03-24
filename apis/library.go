package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync/atomic"

	"github.com/kr/pretty"
	"github.com/nanoteck137/nosepass/core"
	"github.com/nanoteck137/nosepass/core/log"
	"github.com/nanoteck137/nosepass/database"
	"github.com/nanoteck137/nosepass/library"
	"github.com/nanoteck137/nosepass/types"
	"github.com/nanoteck137/pyrin"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type GetLibraryStatus struct {
	Syncing bool `json:"syncing"`
}

var syncing = atomic.Bool{}

type filedata struct {
	chapters    []database.MediaChapter
	subtitles   []database.MediaSubtitle
	attachments []database.MediaAttachment
	videoTracks []database.MediaVideoTrack
	audioTracks []database.MediaAudioTrack
}

func getFileData(ctx context.Context, p string) (filedata, error) {
	probeData, err := ffprobe.ProbeURL(ctx, p)
	if err != nil {
		return filedata{}, err
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
				log.Warn("Unsupported subtitle type (skipping)", "path", p, "type", typ)
				continue
			}

		case "attachment":
			attachments = append(attachments, database.MediaAttachment{
				Index:           i,
				AttachmentIndex: len(attachments) + 1,
			})
		}
	}

	return filedata{
		chapters:    chapters,
		subtitles:   subtitles,
		attachments: attachments,
		videoTracks: videoTracks,
		audioTracks: audioTracks,
	}, nil
}

func getOrCreateMedia(db *database.Database, ctx context.Context, collectionId, p string) (database.Media, error) {
	media, err := db.GetMediaByPath(ctx, p)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			res, err := db.CreateMedia(ctx, database.CreateMediaParams{
				Path:         p,
				CollectionId: collectionId,
			})
			if err != nil {
				return database.Media{}, err
			}

			media, err := db.GetMediaById(ctx, res.Id)
			if err != nil {
				return database.Media{}, err
			}

			return media, nil
		}

		return database.Media{}, err
	}

	return media, err
}

func getOrCreateCollection(ctx context.Context, db *database.Database, path, name string) (database.Collection, error) {
	col, err := db.GetCollectionByPath(ctx, path)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			col, err := db.CreateCollection(ctx, database.CreateCollectionParams{
				Path: path,
				Name: name,
			})
			if err != nil {
				return database.Collection{}, err
			}

			return col, nil
		}

		return database.Collection{}, err
	}

	return col, nil
}

func syncLibrary(app core.App) error {
	syncing.Store(true)
	defer syncing.Store(false)

	lib, err := library.ReadFromDisk("/Volumes/media2/test")
	if err != nil {
		return err
	}

	ctx := context.TODO()

	for _, collection := range lib.Collections {
		dbCollection, err := getOrCreateCollection(ctx, app.DB(), collection.Path, collection.Name)
		if err != nil {
			return err
		}

		err = app.DB().UpdateCollection(ctx, dbCollection.Id, database.CollectionChanges{
			Name: types.Change[string]{
				Value:   collection.Name,
				Changed: collection.Name != dbCollection.Name,
			},
		})
		if err != nil {
			return err
		}

		for _, mediaItem := range collection.Media {
			media, err := getOrCreateMedia(app.DB(), ctx, dbCollection.Id, mediaItem.Path)
			if err != nil {
				return err
			}

			// err = app.DB().DeleteAllMediaEpisodes(ctx, media.Id)
			// if err != nil {
			// 	return err
			// }

			if media.FileModifiedTime < mediaItem.ModTime.UnixMilli() {
				mediaDir := app.WorkDir().MediaIdDir(media.Id)

				err := os.RemoveAll(mediaDir.String())
				if err != nil && !os.IsNotExist(err) {
					return err
				}

				dirs := []string{
					mediaDir.String(),
					mediaDir.Subtitles(),
					mediaDir.Attachments(),
				}

				for _, dir := range dirs {
					err = os.Mkdir(dir, 0755)
					if err != nil && !os.IsExist(err) {
						return err
					}
				}

				data, err := getFileData(ctx, media.Path)
				if err != nil {
					return err
				}

				if len(data.subtitles) > 0 {
					args := []string{
						"-nostats",
						"-hide_banner",
						"-loglevel", "error",
						"-i", media.Path,
					}

					// ffmpeg -i Movie.mkv -map 0:s:0 subs.srt
					for _, subtitle := range data.subtitles {
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
						return err
					}
				}

				err = app.DB().UpdateMedia(ctx, media.Id, database.MediaChanges{
					FileModifiedTime: types.Change[int64]{
						Value:   mediaItem.ModTime.UnixMilli(),
						Changed: true,
					},

					Chapters: types.Change[[]database.MediaChapter]{
						Value:   data.chapters,
						Changed: true,
					},
					Subtitles: types.Change[[]database.MediaSubtitle]{
						Value:   data.subtitles,
						Changed: true,
					},
					Attachments: types.Change[[]database.MediaAttachment]{
						Value:   data.attachments,
						Changed: true,
					},
					VideoTracks: types.Change[[]database.MediaVideoTrack]{
						Value:   data.videoTracks,
						Changed: true,
					},
					AudioTracks: types.Change[[]database.MediaAudioTrack]{
						Value:   data.audioTracks,
						Changed: true,
					},
				})
				if err != nil {
					return err
				}

			// 	variants := media.Variants.GetOrEmpty()
			// 	if len(variants) <= 0 {
			// 		for _, audio := range data.audioTracks {
			// 			switch audio.Language {
			// 			case "eng":
			// 				_, err := app.DB().CreateMediaVariant(ctx, database.CreateMediaVariantParams{
			// 					MediaId:    media.Id,
			// 					Name:       "English",
			// 					Language:   "en",
			// 					VideoTrack: 0,
			// 					AudioTrack: audio.AudioIndex,
			// 					// Subtitle:   sql.NullString{},
			// 				})
			// 				if err != nil {
			// 					return err
			// 				}
			// 			case "und":
			// 				subIndex := int64(-1)
			//
			// 				// TODO(patrik): Better selection
			// 				for _, subtitle := range data.subtitles {
			// 					if subtitle.IsDefault {
			// 						subIndex = int64(subtitle.SubtitleIndex)
			// 						break
			// 					}
			// 				}
			//
			// 				_, err := app.DB().CreateMediaVariant(ctx, database.CreateMediaVariantParams{
			// 					MediaId:    media.Id,
			// 					Name:       "Undefined",
			// 					Language:   "und",
			// 					VideoTrack: 0,
			// 					AudioTrack: audio.AudioIndex,
			// 					Subtitle: sql.NullInt64{
			// 						Int64: subIndex,
			// 						Valid: subIndex != -1,
			// 					},
			// 				})
			// 				if err != nil {
			// 					return err
			// 				}
			// 			case "jpn":
			// 				subIndex := int64(-1)
			//
			// 				// TODO(patrik): Better selection
			// 				for _, subtitle := range data.subtitles {
			// 					if subtitle.IsDefault {
			// 						subIndex = int64(subtitle.SubtitleIndex)
			// 						break
			// 					}
			// 				}
			//
			// 				_, err := app.DB().CreateMediaVariant(ctx, database.CreateMediaVariantParams{
			// 					MediaId:    media.Id,
			// 					Name:       "Japanese",
			// 					Language:   "jp",
			// 					VideoTrack: 0,
			// 					AudioTrack: audio.AudioIndex,
			// 					Subtitle: sql.NullInt64{
			// 						Int64: subIndex,
			// 						Valid: subIndex != -1,
			// 					},
			// 				})
			// 				if err != nil {
			// 					return err
			// 				}
			// 			default:
			// 				log.Warn("Unsupported audio language", "path", media.Path, "language", audio.Language)
			// 			}
			// 		}
			// 	}
			}

			// serie, err := app.DB().GetSerieByName(ctx, serie.Name)
			// if err != nil {
			// 	if errors.Is(err, database.ErrItemNotFound) {
			// 		continue
			// 	} else {
			// 		return err
			// 	}
			// }

			// n, err := strconv.ParseInt(strings.TrimLeft(collection.Name, "Season "), 10, 0)
			// if err != nil {
			// 	return err
			// }
			//
			// fmt.Printf("n: %v\n", n)

			// var r = regexp.MustCompile(`.+S(\d+)E(\d+)`)
			//
			// m := r.FindStringSubmatch(mediaItem.Name)
			//
			// s, err := strconv.ParseInt(m[1], 10, 0)
			// if err != nil {
			// 	return err
			// }
			//
			// e, err := strconv.ParseInt(m[2], 10, 0)
			// if err != nil {
			// 	return err
			// }
			//
			// season, err := app.DB().GetSeason(ctx, serie.Id, int(s))
			// if err != nil {
			// 	if errors.Is(err, database.ErrItemNotFound) {
			// 		continue
			// 	} else {
			// 		return err
			// 	}
			// }
			//
			// episodes, err := app.DB().GetEpisodes(ctx, season.Id)
			// if err != nil {
			// 	return err
			// }
			//
			// for _, episode := range episodes {
			// 	if episode.SeasonNumber.Int64 == e {
			// 		err := app.DB().AddMediaToEpisode(ctx, media.Id, episode.Id)
			// 		if err != nil {
			// 			return err
			// 		}
			//
			// 		fmt.Printf("%s -> %s/%s/%s\n", mediaItem.Path, serie.Name, season.Name, episode.Name)
			//
			// 		break
			// 	}
			// }
		}
	}

	return nil
}

func InstallLibraryHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetLibraryStatus",
			Method:       http.MethodGet,
			Path:         "/library",
			ResponseType: GetLibraryStatus{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				return GetLibraryStatus{
					Syncing: syncing.Load(),
				}, nil
			},
		},
		pyrin.ApiHandler{
			Name:   "SyncLibrary",
			Method: http.MethodPost,
			Path:   "/library",
			HandlerFunc: func(c pyrin.Context) (any, error) {
				go func() {
					err := syncLibrary(app)
					if err != nil {
						log.Error("Failed to sync library", "err", err)
					}
				}()

				return nil, nil
			},
		},
	)
}
