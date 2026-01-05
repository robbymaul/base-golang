package web

type ResponseWeb struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
}

type Metadata struct {
	PerPage    int    `json:"perPage"`
	PageCount  int    `json:"pageCount"`
	TotalCount int    `json:"totalCount"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	OrderBy    string `json:"orderBy"`
}

type ListResponse struct {
	Items    any       `json:"data"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

type Session struct {
	Token     string `json:"token"`
	ExpiredAt int64  `json:"expiredAt"`
}

type Health struct {
	AppName    string `json:"appName"`
	Uptime     string `json:"uptime"`
	AppVersion string `json:"appVersion"`
}
