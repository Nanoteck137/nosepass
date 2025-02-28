package library

import (
	"os"
	"path"
	"time"
)

type Media struct {
	Name    string
	Path    string
	ModTime time.Time
}

type Collection struct {
	Name string
	Path string

	MediaItems []Media
}

type Serie struct {
	Name string
	Path string

	Collections []Collection
}

type Library struct {
	Path   string
	Series []Serie
}

func readMedia(p string) ([]Media, error) {
	entries, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}

	var items []Media

	for _, entry := range entries {
		name := entry.Name()
		if name[0] == '.' {
			continue
		}

		p := path.Join(p, name)

		ext := path.Ext(p)
		switch ext {
		case ".mkv", ".mp4":
			info, err := entry.Info()
			if err != nil {
				return nil, err
			}

			items = append(items, Media{
				Name:    name,
				Path:    p,
				ModTime: info.ModTime(),
			})
		}
	}

	return items, nil
}

func ReadFromDisk(p string) (*Library, error) {
	entries, err := os.ReadDir(p)
	if err != nil {
		return nil, err
	}

	var series []Serie

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		p := path.Join(p, entry.Name())

		entries, err := os.ReadDir(p)
		if err != nil {
			return nil, err
		}

		var collections []Collection

		for _, entry := range entries {
			p := path.Join(p, entry.Name())

			if entry.IsDir() {
				mediaItems, err := readMedia(p)
				if err != nil {
					return nil, err
				}

				collections = append(collections, Collection{
					Name:       entry.Name(),
					Path:       p,
					MediaItems: mediaItems,
				})
			}
		}

		series = append(series, Serie{
			Name:        entry.Name(),
			Path:        p,
			Collections: collections,
		})
	}

	return &Library{
		Path:   p,
		Series: series,
	}, nil
}
