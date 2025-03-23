package library

import (
	"encoding/json"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/kr/pretty"
)

type Media struct {
	Name    string
	Path    string
	ModTime time.Time
}

type Collection struct {
	Name string
	Path string

	Media []Media
}

type Library struct {
	Path        string
	Collections []Collection
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

type CollectionMetadataMediaVariant struct {
	VideoTrackIndex int `json:"videoTrackIndex"`
	AudioTrackIndex int `json:"audioTrackIndex"`
	SubtitleIndex   int `json:"subtitleIndex"`
}

type CollectionMetadataMedia struct {
	Name     string `json:"name"`
	Filename string `json:"filename"`

	Sub *CollectionMetadataMediaVariant `json:"sub"`
	Dub *CollectionMetadataMediaVariant `json:"dub"`
}

type CollectionMetadata struct {
	Name string `json:"name"`

	Items []CollectionMetadataMedia `json:"items"`
}

func ReadFromDisk(p string) (*Library, error) {
	var collectionPaths []string

	filepath.WalkDir(p, func(p string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		name := d.Name()
		if strings.HasPrefix(name, ".") {
			return nil
		}

		if name == "collection.json" {
			b := path.Dir(p)
			collectionPaths = append(collectionPaths, b)
		}

		return nil
	})

	pretty.Println(collectionPaths)

	// entries, err := os.ReadDir(p)
	// if err != nil {
	// 	return nil, err
	// }

	var collections []Collection

	for _, p := range collectionPaths {
		d, err := os.ReadFile(path.Join(p, "collection.json"))
		if err != nil {
			return nil, err
		}

		// TODO(patrik): Validate this
		var metadata CollectionMetadata
		err = json.Unmarshal(d, &metadata)
		if err != nil {
			return nil, err
		}

		media, err := readMedia(p)
		if err != nil {
			return nil, err
		}

		collections = append(collections, Collection{
			Name:  metadata.Name,
			Path:  p,
			Media: media,
		})
	}

	// for _, entry := range entries {
	// 	if !entry.IsDir() {
	// 		continue
	// 	}
	//
	// 	p := path.Join(p, entry.Name())
	//
	// 	// entries, err := os.ReadDir(p)
	// 	// if err != nil {
	// 	// 	return nil, err
	// 	// }
	//
	// 	mappedCollections := make(map[string]MediaCollection)
	//
	// 	err = filepath.WalkDir(p, func(p string, d fs.DirEntry, err error) error {
	// 		// pretty.Println(d)
	//
	// 		if d.IsDir() {
	// 			return nil
	// 		}
	//
	// 		name := d.Name()
	// 		if strings.HasPrefix(name, ".") {
	// 			return nil
	// 		}
	//
	// 		ext := path.Ext(p)
	// 		switch ext {
	// 		case ".mp4", ".mkv":
	// 			name := path.Base(p)
	//
	// 			info, err := entry.Info()
	// 			if err != nil {
	// 				return err
	// 			}
	//
	// 			item := Media{
	// 				Name:    name,
	// 				Path:    p,
	// 				ModTime: info.ModTime(),
	// 			}
	//
	// 			dir := path.Dir(p)
	// 			col, exists := mappedCollections[dir]
	// 			if !exists {
	// 				b := path.Base(dir)
	// 				mappedCollections[dir] = MediaCollection{
	// 					Name:       b,
	// 					Path:       dir,
	// 					MediaItems: []Media{item},
	// 				}
	// 			} else {
	// 				col.MediaItems = append(col.MediaItems, item)
	// 				mappedCollections[dir] = col
	// 			}
	// 		}
	//
	// 		return nil
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	//
	// 	collections := make([]MediaCollection, 0, len(mappedCollections))
	//
	// 	for _, col := range mappedCollections {
	// 		collections = append(collections, col)
	// 	}
	//
	// 	series = append(series, Collection{
	// 		Name:             entry.Name(),
	// 		Path:             p,
	// 		MediaCollections: collections,
	// 	})
	// }

	return &Library{
		Path:        p,
		Collections: collections,
	}, nil
}
