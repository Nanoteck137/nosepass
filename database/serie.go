package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/nosepass/tools/utils"
	"github.com/nanoteck137/nosepass/types"
)

type Serie struct {
	RowId int `db:"rowid"`

	Id   string `db:"id"`
	Name string `db:"name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func SerieQuery() *goqu.SelectDataset {
	query := dialect.From("series").
		Select(
			"series.rowid",

			"series.id",
			"series.name",

			"series.updated",
			"series.created",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllSeries(ctx context.Context) ([]Serie, error) {
	query := SerieQuery()

	var items []Serie
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetSerieById(ctx context.Context, id string) (Serie, error) {
	query := SerieQuery().
		Where(goqu.I("series.id").Eq(id))

	return single[Serie](db, query)
}

func (db *Database) GetSerieByName(ctx context.Context, name string) (Serie, error) {
	query := SerieQuery().
		Where(goqu.I("series.name").Eq(name))

	return single[Serie](db, query)
}

type CreateSerieParams struct {
	Id   string
	Name string

	Created int64
	Updated int64
}

func (db *Database) CreateSerie(ctx context.Context, params CreateSerieParams) (Serie, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateSerieId()
	}

	query := dialect.Insert("series").Rows(goqu.Record{
		"id":   id,
		"name": params.Name,

		"created": created,
		"updated": updated,
	}).
		Returning(
			"series.rowid",

			"series.id",
			"series.name",

			"series.updated",
			"series.created",
		).
		Prepared(true)

	var item Serie
	err := db.Get(&item, query)
	if err != nil {
		return Serie{}, err
	}

	return item, nil
}

type SerieChanges struct {
	Name types.Change[string]

	Created types.Change[int64]
}

func (db *Database) UpdateSerie(ctx context.Context, id string, changes SerieChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("series").
		Set(record).
		Where(goqu.I("series.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteSerie(ctx context.Context, id string) error {
	query := dialect.Delete("series").
		Prepared(true).
		Where(goqu.I("series.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
