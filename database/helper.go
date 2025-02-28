package database

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/nosepass/types"
)

func addToRecord[T any](record goqu.Record, name string, change types.Change[T]) {
	if change.Changed {
		record[name] = change.Value
	}
}

func addToRecordJson[T any](record goqu.Record, name string, change types.Change[T]) {
	if change.Changed {
		record[name] = NewJsonColumn(change.Value)
	}
}

func single[T any](db *Database, query *goqu.SelectDataset) (T, error) {
	var item T
	err := db.Get(&item, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return item, ErrItemNotFound
		}

		return item, err
	}

	return item, nil
}

type JsonColumn[T any] struct {
	v *T
}

func NewJsonColumn[T any](value T) JsonColumn[T] {
	return JsonColumn[T]{
		v: &value,
	}
}

func (j *JsonColumn[T]) Scan(src any) error {
	if src == nil {
		j.v = nil
		return nil
	}
	j.v = new(T)

	switch value := src.(type) {
	case string:
		return json.Unmarshal([]byte(value), j.v)
	case []byte:
		return json.Unmarshal(value, j.v)
	default:
		return fmt.Errorf("unsupported type %T", src)
	}
}

func (j JsonColumn[T]) Value() (driver.Value, error) {
	raw, err := json.Marshal(j.v)
	return raw, err
}

func (j *JsonColumn[T]) Get() *T {
	return j.v
}

func (j *JsonColumn[T]) GetOrEmpty() T {
	if j.v == nil {
		var empty T
		return empty
	}

	return *j.v
}
