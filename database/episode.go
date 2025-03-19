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

type EpisodeVariant struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Language string `json:"language"`
}

type Episode struct {
	RowId int `db:"rowid"`

	Id       string `db:"id"`
	SeasonId string `db:"season_id"`

	Name         string        `db:"name"`
	SeasonNumber sql.NullInt64 `db:"season_number"`
	SerieNumber  sql.NullInt64 `db:"serie_number"`

	Created int64 `db:"created"`
	Updated int64 `db:"updated"`

	Variants JsonColumn[[]EpisodeVariant] `db:"variants"`
}

func EpisodeQuery() *goqu.SelectDataset {
	variants := dialect.From("media_episodes").
		Select(
			goqu.I("media_episodes.episode_id").As("episode_id"),

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
				),
			).As("variants"),
		).
		Join(
			goqu.T("media"),
			goqu.On(goqu.I("media_episodes.media_id").Eq(goqu.I("media.id"))),
		).
		Join(
			goqu.T("media_variants"),
			goqu.On(goqu.I("media.id").Eq(goqu.I("media_variants.media_id"))),
		).
		GroupBy(goqu.I("media_episodes.episode_id"))

	query := dialect.From("episodes").
		Select(
			"episodes.rowid",

			"episodes.id",
			"episodes.season_id",

			"episodes.name",
			"episodes.season_number",
			"episodes.serie_number",

			"episodes.updated",
			"episodes.created",

			goqu.I("variants.variants").As("variants"),
		).
		Prepared(true).
		LeftJoin(
			variants.As("variants"),
			goqu.On(goqu.I("episodes.id").Eq(goqu.I("variants.episode_id"))),
		)

	return query
}

func (db *Database) GetAllEpisodes(ctx context.Context, filterStr, sortStr string) ([]Episode, error) {
	query := EpisodeQuery()

	var items []Episode
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetEpisodes(ctx context.Context, seasonId string) ([]Episode, error) {
	query := EpisodeQuery().
		Where(goqu.I("episodes.season_id").Eq(seasonId))

	var items []Episode
	err := db.Select(&items, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (db *Database) GetEpisodeById(ctx context.Context, id string) (Episode, error) {
	query := EpisodeQuery().
		Where(goqu.I("episodes.id").Eq(id))

	var item Episode
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Episode{}, ErrItemNotFound
		}

		return Episode{}, err
	}

	return item, nil
}

type CreateEpisodeParams struct {
	Id       string
	SeasonId string

	Name         string
	SeasonNumber sql.NullInt64
	SerieNumber  sql.NullInt64

	Created int64
	Updated int64
}

func (db *Database) CreateEpisode(ctx context.Context, params CreateEpisodeParams) (Episode, error) {
	t := time.Now().UnixMilli()
	created := params.Created
	updated := params.Updated

	if created == 0 && updated == 0 {
		created = t
		updated = t
	}

	id := params.Id
	if id == "" {
		id = utils.CreateEpisodeId()
	}

	query := dialect.Insert("episodes").Rows(goqu.Record{
		"id":        id,
		"season_id": params.SeasonId,

		"name":          params.Name,
		"season_number": params.SeasonNumber,
		"serie_number":  params.SerieNumber,

		"created": created,
		"updated": updated,
	}).
		Returning(
			"episodes.rowid",

			"episodes.id",
			"episodes.season_id",

			"episodes.name",
			"episodes.season_number",
			"episodes.serie_number",

			"episodes.updated",
			"episodes.created",
		).
		Prepared(true)

	var item Episode
	err := db.Get(&item, query)
	if err != nil {
		return Episode{}, err
	}

	return item, nil
}

type EpisodeChanges struct {
	SeasonId types.Change[string]

	Name         types.Change[string]
	SeasonNumber types.Change[sql.NullInt64]
	SerieNumber  types.Change[sql.NullInt64]

	Created types.Change[int64]
}

func (db *Database) UpdateEpisode(ctx context.Context, id string, changes EpisodeChanges) error {
	record := goqu.Record{}

	addToRecord(record, "season_id", changes.SeasonId)

	addToRecord(record, "name", changes.Name)
	addToRecord(record, "season_number", changes.SeasonNumber)
	addToRecord(record, "serie_number", changes.SerieNumber)

	addToRecord(record, "created", changes.Created)

	if len(record) == 0 {
		return nil
	}

	record["updated"] = time.Now().UnixMilli()

	ds := dialect.Update("episodes").
		Set(record).
		Where(goqu.I("episodes.id").Eq(id)).
		Prepared(true)

	_, err := db.Exec(ctx, ds)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteEpisode(ctx context.Context, id string) error {
	query := dialect.Delete("episodes").
		Prepared(true).
		Where(goqu.I("episodes.id").Eq(id))

	_, err := db.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
