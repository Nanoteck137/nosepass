package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/nosepass/tools/utils"
	"github.com/nanoteck137/nosepass/types"
)

type Season struct {
	RowId int `db:"rowid"`

	Id      string `db:"id"`
	SerieId string `db:"serie_id"`

	Name   string `db:"name"`
	Number int64  `db:"number"`
	Type   string `db:"type"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func SeasonQuery() *goqu.SelectDataset {
	query := dialect.From("seasons").
		Select(
			"seasons.rowid",

			"seasons.id",
			"seasons.serie_id",

			"seasons.name",
			"seasons.number",
			"seasons.type",

			"seasons.updated",
			"seasons.created",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllSeasons(ctx context.Context, filterStr, sortStr string) ([]Season, error) {
	query := SeasonQuery()

	var items []Season
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetSeasonById(ctx context.Context, id string) (Season, error) {
	query := SeasonQuery().
		Where(goqu.I("seasons.id").Eq(id))

	var item Season
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Season{}, ErrItemNotFound
		}

		return Season{}, err
	}

	return item, nil
}

func (db *Database) GetSeason(ctx context.Context, serieId string, seasonNumber int) (Season, error) {
	query := SeasonQuery().
		Where(
			goqu.I("seasons.serie_id").Eq(serieId),
			goqu.I("seasons.number").Eq(seasonNumber),
		)

	return single[Season](db, query)
}

type CreateSeasonParams struct {
	Id      string
	SerieId string

	Name   string
	Number int64
	Type   string

	Created int64
	Updated int64
}

func (db *Database) CreateSeason(ctx context.Context, params CreateSeasonParams) (Season, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateSeasonId()
	}

	query := dialect.Insert("seasons").Rows(goqu.Record{
		"id":   id,
		"serie_id":   params.SerieId,

		"name": params.Name,
		"number": params.Number,
		"type": params.Type,

		"created": created,
		"updated": updated,
	}).
		Returning(
			"seasons.rowid",

			"seasons.id",
			"seasons.serie_id",

			"seasons.name",
			"seasons.number",
			"seasons.type",

			"seasons.updated",
			"seasons.created",
		).
		Prepared(true)

	var item Season
	err := db.Get(&item, query)
	if err != nil {
		return Season{}, err
	}

	return item, nil
}

type SeasonChanges struct {
	SerieId types.Change[string]

	Name   types.Change[string]
	Number types.Change[int64]
	Type   types.Change[string]

	Created types.Change[int64]
}

func (db *Database) UpdateSeason(ctx context.Context, id string, changes SeasonChanges) error {
	record := goqu.Record{}

	addToRecord(record, "serie_id", changes.SerieId)

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "number", changes.Number)
	addToRecord(record, "type", changes.Type)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("seasons").
		Set(record).
		Where(goqu.I("seasons.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteSeason(ctx context.Context, id string) error {
	query := dialect.Delete("seasons").
		Prepared(true).
		Where(goqu.I("seasons.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
