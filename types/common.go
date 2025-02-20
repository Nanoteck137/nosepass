package types

import "path"

type MediaType string

const (
	MediaTypeFlac      MediaType = "flac"
	MediaTypeOggOpus   MediaType = "ogg-opus"
	MediaTypeOggVorbis MediaType = "ogg-vorbis"
	MediaTypeMp3       MediaType = "mp3"
	MediaTypeAcc       MediaType = "acc"
)

type Map map[string]any

type WorkDir string

func (d WorkDir) String() string {
	return string(d)
}

func (d WorkDir) DatabaseFile() string {
	return path.Join(d.String(), "data.db")
}

func (d WorkDir) SetupFile() string {
	return path.Join(d.String(), "setup")
}

type Change[T any] struct {
	Value   T
	Changed bool
}
