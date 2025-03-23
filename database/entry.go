package database

import (
	"context"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/nosepass/tools/utils"
	"github.com/nanoteck137/nosepass/types"
)

type Entry struct {
	RowId int `db:"rowid"`

	Id   string `db:"id"`

	Name string `db:"name"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`
}

func EntryQuery() *goqu.SelectDataset {
	query := dialect.From("entries").
		Select(
			"entries.rowid",

			"entries.id",
			"entries.name",

			"entries.updated",
			"entries.created",
		).
		Prepared(true)

	return query
}

func (db *Database) GetAllEntries(ctx context.Context) ([]Entry, error) {
	query := EntryQuery()

	var items []Entry
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetEntryById(ctx context.Context, id string) (Entry, error) {
	query := EntryQuery().
		Where(goqu.I("entries.id").Eq(id))

	return single[Entry](db, query)
}

func (db *Database) GetEntryByName(ctx context.Context, name string) (Entry, error) {
	query := EntryQuery().
		Where(goqu.I("entries.name").Eq(name))

	return single[Entry](db, query)
}

type CreateEntryParams struct {
	Id   string
	Name string

	Created int64
	Updated int64
}

func (db *Database) CreateEntry(ctx context.Context, params CreateEntryParams) (Entry, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateEntryId()
	}

	query := dialect.Insert("entries").Rows(goqu.Record{
		"id":   id,
		"name": params.Name,

		"created": created,
		"updated": updated,
	}).
		Returning(
			"entries.rowid",

			"entries.id",
			"entries.name",

			"entries.updated",
			"entries.created",
		).
		Prepared(true)

	var item Entry
	err := db.Get(&item, query)
	if err != nil {
		return Entry{}, err
	}

	return item, nil
}

type EntryChanges struct {
	Name types.Change[string]

	Created types.Change[int64]
}

func (db *Database) UpdateEntry(ctx context.Context, id string, changes EntryChanges) error {
	record := goqu.Record{}

	addToRecord(record, "name", changes.Name)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("entries").
		Set(record).
		Where(goqu.I("entries.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteEntry(ctx context.Context, id string) error {
	query := dialect.Delete("entries").
		Prepared(true).
		Where(goqu.I("entries.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
