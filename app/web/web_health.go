package web

type HealthResponseWeb struct {
	AppName    string       `json:"appName"`
	Uptime     string       `json:"uptime"`
	AppVersion string       `json:"appVersion"`
	Resource   ResourcesWeb `json:"resource"`
}

type ResourcesWeb struct {
	Database string `json:"database"`
}
