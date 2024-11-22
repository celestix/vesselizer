package common

type UpdateType int

const (
	UPDATE_ECHO UpdateType = iota
	UPDATE_NEW_VESSEL
	UPDATE_STOP_VESSEL
	UPDATE_START_VESSEL

	UPDATE_RELAY_STDOUT
	UPDATE_RELAY_STDERR
	UPDATE_RELAY_STDIN
)

var UpdateTypeToString = map[UpdateType]string{
	UPDATE_ECHO:         "echo",
	UPDATE_NEW_VESSEL:   "newVessel",
	UPDATE_STOP_VESSEL:  "stopVessel",
	UPDATE_START_VESSEL: "startVessel",
	UPDATE_RELAY_STDIN:  "relayStdin",
	UPDATE_RELAY_STDOUT: "relayStdout",
	UPDATE_RELAY_STDERR: "relayStderr",
}

type EchoUpdate struct {
	Message string `json:"message"`
}

type NewVesselRequest struct {
	Name         string   `json:"name"`
	BaseImage    string   `json:"base_image"`
	Entrypoint   []string `json:"entrypoint"`
	BuildFile    string   `json:"build_file"`
	AppDirectory string   `json:"app_directory"`
}

type NewVesselResponse struct {
	Id string `json:"id"`
}

type VesselControlRequest struct {
	Id string `json:"id"`
}

type VesselControlResponse struct {
	Success bool `json:"success"`
}

type RelayedData struct {
	Id   string `json:"id"`
	Data []byte `json:"data"`
}
