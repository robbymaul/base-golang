package repositories

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"paymentserviceklink/config"
	"paymentserviceklink/pkg/database"
	"paymentserviceklink/pkg/util"
)

type Repository struct {
	db      *database.Database
	cfg     *config.Config
	Adapter *Adapter
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	maxIdleConn := util.ParseIntFallback(cfg.DatabaseMaxIdleConn, 10)
	maxOpenConn := util.ParseIntFallback(cfg.DatabaseOpenConn, 100)
	maxConnLifetime := util.ParseIntFallback(cfg.DatabaseMaxConnLifetime, 10)

	database, err := database.NewDatabase(database.Config{
		Host:            cfg.DatabaseHost,
		Port:            cfg.DatabasePort,
		Username:        cfg.DatabaseUser,
		Password:        cfg.DatabasePass,
		Database:        cfg.DatabaseName,
		Driver:          database.DriverPostgresSql,
		SslMode:         cfg.DatabaseSslMode,
		TimeZone:        cfg.DatabaseTimezone,
		MaxIdleConn:     &maxIdleConn,
		MaxOpenConn:     &maxOpenConn,
		MaxConnLifetime: &maxConnLifetime,
		AppMode:         cfg.AppMode,
	})
	if err != nil {
		log.Error().Msg(fmt.Sprintf("failed init database with error = [%v]", err))
		return nil, err
	}

	adapter, err := NewRepositoryAdapter(cfg)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("failed init adapter with error = [%v]", err))
		return nil, err
	}

	return &Repository{db: database, cfg: cfg, Adapter: adapter}, nil
}

func (r *Repository) Connect(ctx context.Context) *RepositoryContext {
	rc, err := r.Connected(ctx)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("failed connect with error = [%v]", err))
		panic(err)
	}
	return rc
}

func (r *Repository) Connected(ctx context.Context) (*RepositoryContext, error) {
	err := r.db.Init()
	if err != nil {
		panic(fmt.Errorf("failed init database with error \n > [%v]", err))
	}

	log.Info().Msg(fmt.Sprint("database : [connected to database]", ctx))

	return &RepositoryContext{
		ctx:     ctx,
		db:      r.db.DB,
		cfg:     r.cfg,
		Adapter: r.Adapter,
	}, nil
}

func (r *Repository) Close() error {
	log.Info().Msg("close database connection")

	if err := r.db.Close(); err != nil {
		log.Error().Msg(fmt.Sprintf("failed close database connection with error = [%v]", err))
		return err
	}

	return nil
}
