package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/nosepass/tools/utils"
	"github.com/nanoteck137/nosepass/types"
)

type MediaVariant struct {
	RowId int `db:"rowid"`

	Id      string `db:"id"`
	MediaId string `db:"media_id"`

	Name       string         `db:"name"`
	Language   string         `db:"language"`
	VideoTrack int            `db:"video_track"`
	AudioTrack int            `db:"audio_track"`
	Subtitle   sql.NullString `db:"subtitle"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func MediaVariantQuery() *goqu.SelectDataset {
	query := dialect.From("media_variants").
		Select(
			"media_variants.rowid",

			"media_variants.id",
			"media_variants.media_id",

			"media_variants.name",
			"media_variants.language",
			"media_variants.video_track",
			"media_variants.audio_track",
			"media_variants.subtitle",

			"media_variants.updated",
			"media_variants.created",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllMediaVariants(ctx context.Context) ([]MediaVariant, error) {
	query := MediaVariantQuery()

	var items []MediaVariant
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetMediaVariantById(ctx context.Context, id string) (MediaVariant, error) {
	query := MediaVariantQuery().
		Where(goqu.I("media_variants.id").Eq(id))

	return single[MediaVariant](db, query)
}

type CreateMediaVariantParams struct {
	Id      string
	MediaId string

	Name       string
	Language   string
	VideoTrack int
	AudioTrack int
	Subtitle   sql.NullString

	Created int64
	Updated int64
}

func (db *Database) CreateMediaVariant(ctx context.Context, params CreateMediaVariantParams) (MediaVariant, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateMediaVariantId()
	}

	query := dialect.Insert("media_variants").Rows(goqu.Record{
		"id":       id,
		"media_id": params.MediaId,

		"name":        params.Name,
		"language":    params.Language,
		"video_track": params.VideoTrack,
		"audio_track": params.AudioTrack,
		"subtitle":    params.Subtitle,

		"created": created,
		"updated": updated,
	}).
		Returning(
			"media_variants.rowid",

			"media_variants.id",
			"media_variants.media_id",

			"media_variants.name",
			"media_variants.language",
			"media_variants.video_track",
			"media_variants.audio_track",
			"media_variants.subtitle",

			"media_variants.updated",
			"media_variants.created",
		).
		Prepared(true)

	var item MediaVariant
	err := db.Get(&item, query)
	if err != nil {
		return MediaVariant{}, err
	}

	return item, nil
}

type MediaVariantChanges struct {
	MediaId types.Change[string]

	Name       types.Change[string]
	Language   types.Change[string]
	VideoTrack types.Change[int]
	AudioTrack types.Change[int]
	Subtitle   types.Change[sql.NullString]

	Created types.Change[int64]
}

func (db *Database) UpdateMediaVariant(ctx context.Context, id string, changes MediaVariantChanges) error {
	record := goqu.Record{}

	addToRecord(record, "media_id", changes.MediaId)

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "language", changes.Language)
	addToRecord(record, "video_track", changes.VideoTrack)
	addToRecord(record, "audio_track", changes.AudioTrack)
	addToRecord(record, "subtitle", changes.Subtitle)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("media_variants").
		Set(record).
		Where(goqu.I("media_variants.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteMediaVariant(ctx context.Context, id string) error {
	query := dialect.Delete("media_variants").
		Prepared(true).
		Where(goqu.I("media_variants.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
