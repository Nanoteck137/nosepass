package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/nosepass/tools/utils"
	"github.com/nanoteck137/nosepass/types"
)

type MediaCollection struct {
	RowId int `db:"rowid"`

	Id   string `db:"id"`
	Path string `db:"path"`

	CollectionId string `db:"collection_id"`

	Name string `db:"name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func MediaCollectionQuery() *goqu.SelectDataset {
	query := dialect.From("media_collections").
		Select(
			"media_collections.rowid",

			"media_collections.id",
			"media_collections.path",

			"media_collections.collection_id",

			"media_collections.name",

			"media_collections.updated",
			"media_collections.created",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllMediaCollections(ctx context.Context) ([]MediaCollection, error) {
	query := MediaCollectionQuery()

	var items []MediaCollection
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetMediaCollectionById(ctx context.Context, id string) (MediaCollection, error) {
	query := MediaCollectionQuery().
		Where(goqu.I("media_collections.id").Eq(id))

	return single[MediaCollection](db, query)
}

func (db *Database) GetMediaCollectionByPath(ctx context.Context, path string) (MediaCollection, error) {
	query := MediaCollectionQuery().
		Where(goqu.I("media_collections.path").Eq(path))

	return single[MediaCollection](db, query)
}

type CreateMediaCollectionParams struct {
	Id   string
	Path string

	CollectionId string

	Name string

	Created int64
	Updated int64
}

func (db *Database) CreateMediaCollection(ctx context.Context, params CreateMediaCollectionParams) (MediaCollection, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateMediaCollectionId()
	}

	query := dialect.Insert("media_collections").Rows(goqu.Record{
		"id":   id,
		"path": params.Path,

		"collection_id": params.CollectionId,

		"name": params.Name,

		"created": created,
		"updated": updated,
	}).
		Returning(
			"media_collections.rowid",

			"media_collections.id",
			"media_collections.path",

			"media_collections.collection_id",

			"media_collections.name",

			"media_collections.updated",
			"media_collections.created",
		).
		Prepared(true)

	return single[MediaCollection](db, query)
}

type MediaCollectionChanges struct {
	Path types.Change[string]

	CollectionId types.Change[string]

	Name types.Change[string]

	Created types.Change[int64]
}

func (db *Database) UpdateMediaCollection(ctx context.Context, id string, changes MediaCollectionChanges) error {
	record := goqu.Record{}

	addToRecord(record, "path", changes.Path)

	addToRecord(record, "collection_id", changes.CollectionId)

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("media_collections").
		Set(record).
		Where(goqu.I("media_collections.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteMediaCollection(ctx context.Context, id string) error {
	query := dialect.Delete("media_collections").
		Prepared(true).
		Where(goqu.I("media_collections.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
