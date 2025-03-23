package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/nosepass/tools/utils"
	"github.com/nanoteck137/nosepass/types"
)

type Collection struct {
	RowId int `db:"rowid"`

	Id   string `db:"id"`
	Path string `db:"path"`

	Name string `db:"name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func CollectionQuery() *goqu.SelectDataset {
	query := dialect.From("collections").
		Select(
			"collections.rowid",

			"collections.id",
			"collections.path",

			"collections.name",

			"collections.updated",
			"collections.created",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllCollections(ctx context.Context) ([]Collection, error) {
	query := CollectionQuery()

	var items []Collection
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetCollectionById(ctx context.Context, id string) (Collection, error) {
	query := CollectionQuery().
		Where(goqu.I("collections.id").Eq(id))

	return single[Collection](db, query)
}

func (db *Database) GetCollectionByPath(ctx context.Context, path string) (Collection, error) {
	query := CollectionQuery().
		Where(goqu.I("collections.path").Eq(path))

	return single[Collection](db, query)
}

type CreateCollectionParams struct {
	Id   string
	Path string

	Name string

	Created int64
	Updated int64
}

func (db *Database) CreateCollection(ctx context.Context, params CreateCollectionParams) (Collection, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateCollectionId()
	}

	query := dialect.Insert("collections").Rows(goqu.Record{
		"id":   id,
		"path": params.Path,

		"name": params.Name,

		"created": created,
		"updated": updated,
	}).
		Returning(
			"collections.rowid",

			"collections.id",
			"collections.path",

			"collections.name",

			"collections.updated",
			"collections.created",
		).
		Prepared(true)

	return single[Collection](db, query)
}

type CollectionChanges struct {
	Path types.Change[string]

	Name types.Change[string]

	Created types.Change[int64]
}

func (db *Database) UpdateCollection(ctx context.Context, id string, changes CollectionChanges) error {
	record := goqu.Record{}

	addToRecord(record, "path", changes.Path)

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("collections").
		Set(record).
		Where(goqu.I("collections.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteCollection(ctx context.Context, id string) error {
	query := dialect.Delete("collections").
		Prepared(true).
		Where(goqu.I("collections.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
