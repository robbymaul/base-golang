package repositories

import (
	"application/config"
	"time"
)

type Adapter struct {
	JakartaLoc *time.Location
}

func NewRepositoryAdapter(cfg *config.Config) (*Adapter, error) {
	adapter := new(Adapter)

	location, err := time.LoadLocation(cfg.ServerTimeZone)
	if err != nil {
		return nil, err
	}

	adapter.JakartaLoc = location

	return adapter, nil
}
