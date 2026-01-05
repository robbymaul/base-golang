package init

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"paymentserviceklink/config"
	"paymentserviceklink/pkg/database"
)

var (
	AppName = "K-Link Address - Service"

	AppVersion = "0.0.1"

	Build = "local.build"
)

func UpgradeDB(config *config.Config) error {
	log.Info().Msg(fmt.Sprintf("Upgrade Database on Boot. Please wait... value = [%t]", config.DatabaseUpgradeOnBoot))

	if !config.DatabaseUpgradeOnBoot {
		return nil
	}

	// load database url
	log.Info().Msg(fmt.Sprintf("Load Database URL. value = [%s]", config.DatabaseDSN))
	uri := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", database.DriverPostgresSql, config.DatabaseUserMigration, config.DatabasePassMigration, config.DatabaseHost, config.DatabasePort, config.DatabaseNameMigration)

	log.Info().Str("uri", uri).Msg("uri migration database")
	// source dir
	path := filepath.Join(config.WorkDir, "/migrations")
	log.Info().Str("path", path).Msg("migrations path")

	source := "file://" + path
	log.Info().Str("source", source).Msg("migrations source")

	log.Info().Msg("migration: connecting database")
	m, err := migrate.New(source, uri)
	if err != nil {
		log.Error().Msg(fmt.Sprint("migrate: failed to connect db", err))
		return err
	}

	// migrate up
	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			log.Info().Msg("migration no change")
			return nil
		}

		log.Error().Msg(fmt.Sprint("migrate: failed to upgrade db migration script", err))

		return err
	}

	// get status
	version, dirty, err := m.Version()
	if err != nil {
		log.Error().Msg(fmt.Sprint("migrate: failed to get db version", err))
		return err
	}

	log.Info().Msg(fmt.Sprintf("migration: success to upgrade db. version = [%d], dirty = [%t]", version, dirty))

	return nil
}
