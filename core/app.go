package core

import (
	"github.com/nanoteck137/nosepass/config"
	"github.com/nanoteck137/nosepass/database"
	"github.com/nanoteck137/nosepass/types"
)

// Inspiration from Pocketbase: https://github.com/pocketbase/pocketbase
// File: https://github.com/pocketbase/pocketbase/blob/master/core/app.go
type App interface {
	DB() *database.Database
	Config() *config.Config

	WorkDir() types.WorkDir

	Bootstrap() error
}
