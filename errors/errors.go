package errors

import e "errors"

var (
	ErrBaseImageNotSpecified      = e.New("base image not specified")
	ErrBuildFileNotSpecified      = e.New("build file not specified")
	ErrEntrypointNotSpecified     = e.New("entrypoint not specified")
	ErrSpecifiedBaseImageNotFound = e.New("specified base image not found")
	ErrVesselNotFound             = e.New("vessel not found")
	ErrVesselNotRunning           = e.New("vessel not running")
)
