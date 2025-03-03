package types

import "path"

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

func (d WorkDir) MediaDir() string {
	return path.Join(d.String(), "media")
}

func (d WorkDir) MediaIdDir(id string) MediaDir {
	return MediaDir(path.Join(d.MediaDir(), id))
}

type MediaDir string

func (d MediaDir) String() string {
	return string(d)
}

func (d MediaDir) Subtitles() string {
	return path.Join(d.String(), "subtitles")
}

func (d MediaDir) Attachments() string {
	return path.Join(d.String(), "attachments")
}

type Change[T any] struct {
	Value   T
	Changed bool
}
