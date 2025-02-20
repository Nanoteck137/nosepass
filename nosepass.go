package nosepass

import "fmt"

var AppName = "nosepass"
var CliAppName = AppName + "-cli"

var Version = "no-version"
var Commit = "no-commit"

func VersionTemplate(appName string) string {
	return fmt.Sprintf(
		"%s: %s (%s)\n",
		appName, Version, Commit)
}
