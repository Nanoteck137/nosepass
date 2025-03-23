package database

import (
	"context"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/nosepass/tools/utils"
	"github.com/nanoteck137/nosepass/types"
)

type MediaChapter struct {
	Title string  `json:"title"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type MediaSubtitle struct {
	Index         int    `json:"index"`
	SubtitleIndex int    `json:"subtitleIndex"`
	Type          string `json:"type"`
	Title         string `json:"title"`
	Language      string `json:"language"`
	IsDefault     bool   `json:"isDefault"`
	Filename      string `json:"filename"`
}

type MediaAttachment struct {
	Index           int `json:"index"`
	AttachmentIndex int `json:"attachmentIndex"`
}

type MediaVideoTrack struct {
	Index      int     `json:"index"`
	VideoIndex int     `json:"videoIndex"`
	Duration   float64 `json:"duration"`
}

type MediaAudioTrack struct {
	Index      int    `json:"index"`
	AudioIndex int    `json:"audioIndex"`
	Language   string `json:"language"`
}

type MediaVariantInfo struct {
	Id string `json:"id"`

	Name       string `json:"name"`
	Language   string `json:"language"`
	VideoTrack int    `json:"video_track"`
	AudioTrack int    `json:"audio_track"`
	Subtitle   *int64 `json:"subtitle"`
}

type Media struct {
	RowId int `db:"rowid"`

	Id   string `db:"id"`
	Path string `db:"path"`

	CollectionId string `db:"collection_id"`

	FileModifiedTime int64 `db:"file_modified_time"`

	Chapters    JsonColumn[[]MediaChapter]    `db:"chapters"`
	Subtitles   JsonColumn[[]MediaSubtitle]   `db:"subtitles"`
	Attachments JsonColumn[[]MediaAttachment] `db:"attachments"`
	VideoTracks JsonColumn[[]MediaVideoTrack] `db:"video_tracks"`
	AudioTracks JsonColumn[[]MediaAudioTrack] `db:"audio_tracks"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Variants JsonColumn[[]MediaVariantInfo] `db:"variants"`
}

func MediaQuery() *goqu.SelectDataset {
	variants := dialect.From("media_variants").
		Select(
			goqu.I("media_variants.media_id").As("media_id"),

			goqu.Func(
				"json_group_array",
				goqu.Func(
					"json_object",

					"id",
					goqu.I("media_variants.id"),
					"name",
					goqu.I("media_variants.name"),
					"language",
					goqu.I("media_variants.language"),
					"video_track",
					goqu.I("media_variants.video_track"),
					"audio_track",
					goqu.I("media_variants.audio_track"),
					"subtitle",
					goqu.I("media_variants.subtitle"),
				),
			).As("variants"),
		).
		GroupBy(goqu.I("media_variants.media_id"))

	query := dialect.From("media").
		Select(
			"media.rowid",

			"media.id",
			"media.path",

			"media.collection_id",

			"media.file_modified_time",

			"media.chapters",
			"media.subtitles",
			"media.attachments",
			"media.video_tracks",
			"media.audio_tracks",

			"media.updated",
			"media.created",

			goqu.I("variants.variants").As("variants"),
		).
		Prepared(true).
		LeftJoin(
			variants.As("variants"),
			goqu.On(goqu.I("media.id").Eq(goqu.I("variants.media_id"))),
		)

	return query
}

func (db *Database) GetAllMedia(ctx context.Context) ([]Media, error) {
	query := MediaQuery()

	var items []Media
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetMediaByCollectionId(ctx context.Context, collectionId string) ([]Media, error) {
	query := MediaQuery().
		Where(goqu.I("media.collection_id").Eq(collectionId))

	var items []Media
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetMediaById(ctx context.Context, id string) (Media, error) {
	query := MediaQuery().
		Where(goqu.I("media.id").Eq(id))

	return single[Media](db, query)
}

func (db *Database) GetMediaByPath(ctx context.Context, path string) (Media, error) {
	query := MediaQuery().
		Where(goqu.I("media.path").Eq(path))

	return single[Media](db, query)
}

type CreateMediaParams struct {
	Id   string
	Path string

	CollectionId string

	FileModifiedTime int64

	Chapters    []MediaChapter
	Subtitles   []MediaSubtitle
	Attachments []MediaAttachment
	VideoTracks []MediaVideoTrack
	AudioTracks []MediaAudioTrack

	Created int64
	Updated int64
}

func (db *Database) CreateMedia(ctx context.Context, params CreateMediaParams) (Media, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateMediaId()
	}

	query := dialect.Insert("media").Rows(goqu.Record{
		"id":   id,
		"path": params.Path,

		"collection_id": params.CollectionId,

		"file_modified_time": params.FileModifiedTime,

		"chapters":     NewJsonColumn(params.Chapters),
		"subtitles":    NewJsonColumn(params.Subtitles),
		"attachments":  NewJsonColumn(params.Attachments),
		"video_tracks": NewJsonColumn(params.VideoTracks),
		"audio_tracks": NewJsonColumn(params.AudioTracks),

		"created": created,
		"updated": updated,
	}).
		Returning(
			"media.rowid",

			"media.id",
			"media.path",

			"media.collection_id",

			"media.file_modified_time",

			"media.chapters",
			"media.subtitles",
			"media.attachments",
			"media.video_tracks",
			"media.audio_tracks",

			"media.updated",
			"media.created",
		).
		Prepared(true)

	var item Media
	err := db.Get(&item, query)
	if err != nil {
		return Media{}, err
	}

	return item, nil
}

type MediaChanges struct {
	FileModifiedTime types.Change[int64]

	Chapters    types.Change[[]MediaChapter]
	Subtitles   types.Change[[]MediaSubtitle]
	Attachments types.Change[[]MediaAttachment]
	VideoTracks types.Change[[]MediaVideoTrack]
	AudioTracks types.Change[[]MediaAudioTrack]

	Created types.Change[int64]
}

func (db *Database) UpdateMedia(ctx context.Context, id string, changes MediaChanges) error {
	record := goqu.Record{}

	addToRecord(record, "file_modified_time", changes.FileModifiedTime)

	addToRecordJson(record, "chapters", changes.Chapters)
	addToRecordJson(record, "subtitles", changes.Subtitles)
	addToRecordJson(record, "attachments", changes.Attachments)
	addToRecordJson(record, "video_tracks", changes.VideoTracks)
	addToRecordJson(record, "audio_tracks", changes.AudioTracks)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("media").
		Set(record).
		Where(goqu.I("media.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteMedia(ctx context.Context, id string) error {
	query := dialect.Delete("media").
		Prepared(true).
		Where(goqu.I("media.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteAllMediaEpisodes(ctx context.Context, mediaId string) error {
	query := dialect.Delete("media_episodes").
		Prepared(true).
		Where(goqu.I("media_episodes.media_id").Eq(mediaId))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddMediaToEpisode(ctx context.Context, mediaId, episodeId string) error {
	query := dialect.Insert("media_episodes").
		Prepared(true).
		Rows(goqu.Record{
			"media_id":   mediaId,
			"episode_id": episodeId,
		})

	_, err := db.Exec(ctx, query)
	if err != nil {
		var sqlErr sqlite3.Error
		if errors.As(err, &sqlErr) {
			if sqlErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
				return ErrItemAlreadyExists
			}
		}

		return err
	}

	return nil
}
