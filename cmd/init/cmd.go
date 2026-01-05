package init

import (
	"application/api/routes"
	"application/app/repositories"
	"application/config"
	"application/pkg/middleware"
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	slug = "Application"
	name = "Application - Service"
)

type BootOptions struct {
	WorkDir   string
	EnvPrefix string
}

type Flags struct {
	CmdPrintVersion         *bool
	OptWorkDir              *string
	OptEnvPrefix            *string
	CreateUserAdmin         *bool
	UpdatePasswordUserAdmin *bool
}

type InitVariables struct {
	StartedAt  time.Time
	AppVersion string
	Signature  string
	Config     *config.Config
	Logger     *zerolog.Logger
}

type Command struct {
	Flags Flags
	Args  []string
	Err   error
}

func NewCLI(args []string) *Command {
	fmt.Println(args)
	return &Command{
		Flags{
			CmdPrintVersion:         flag.Bool("version", false, "Command: show version"),
			OptWorkDir:              flag.String("dir", "", "Option: set working directory"),
			OptEnvPrefix:            flag.String("env-prefix", "", "Option: set env prefix"),
			CreateUserAdmin:         flag.Bool("create-user-admin", false, "Option: create user admin"),
			UpdatePasswordUserAdmin: flag.Bool("update-password-user-admin", false, "Option: update user admin"),
		},
		args,
		nil,
	}
}

func (cmd *Command) Start(inv InitVariables) {
	if *cmd.Flags.CmdPrintVersion {
		log.Info().Str("event", "running service").Msg(fmt.Sprintf("usage: [%s] [options] <additional>\nOptions:\n", os.Args[0]))
		flag.PrintDefaults()
		flag.String(slug, "", "./main -start=true")
		flag.Parse()
		return
	}
	cmd.Run(inv)
}

func (cmd *Command) Run(inv InitVariables) {
	cfg := inv.Config

	logLevel := zerolog.InfoLevel // Default level
	if cfg.LogLevel == "debug" {
		logLevel = zerolog.DebugLevel
	} else if cfg.LogLevel == "warn" {
		logLevel = zerolog.WarnLevel
	} else if cfg.LogLevel == "error" {
		logLevel = zerolog.ErrorLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	repo, err := repositories.NewRepository(cfg)
	if err != nil {
		panic(err)
	}

	repositoryContext := repo.Connect(context.Background())

	handler, err := InitServer(inv.StartedAt, inv.AppVersion, inv.Signature, cfg, repositoryContext)
	if err != nil {
		panic(fmt.Errorf("failed to init controllers. Error = [%v]", err))
	}

	server := &http.Server{
		Addr:    GetListenPort(cfg.ServerPort),
		Handler: handler,
	}

	log.Warn().Str("service", name).Str("version", inv.AppVersion).Time("started_at", inv.StartedAt).Int("running_on", int(cfg.ServerPort)).Msg("service is starting")

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Msg(fmt.Sprintf("Failed to start server: %v", err))
		}
	}()

	cmd.GracefulShutdown(server, cfg, repo)
}

func GetListenPort(i uint16) string {
	if i > 0 {
		return fmt.Sprintf("0.0.0.0:%d", i)
	}

	log.Info().Msg("port is not set. server running to next free port")
	return "0.0.0.0:0"
}

func InitServer(startTime time.Time, appVersion string, signature string, cfg *config.Config, repo *repositories.RepositoryContext) (*gin.Engine, error) {

	if cfg.AppMode == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	} else if (cfg.AppMode == "DEVELOPMENT") || (cfg.AppMode == "TEST") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	e := gin.New()

	e.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"*"},
	}))

	e.Use(gin.Recovery())
	e.Use(gin.Logger())
	e.Use(middleware.ErrorHandler())

	route := routes.NewRoute(startTime, appVersion, cfg, repo, e.Group("/api/v1"))

	route.RegisterCoreServicesRoutes()

	return e, nil
}

func (cmd *Command) GracefulShutdown(server *http.Server, cfg *config.Config, repo *repositories.Repository) {
	quit := make(chan os.Signal, 1)

	defer func() {
		if err2 := repo.Close(); err2 != nil {
			log.Fatal().Msg(fmt.Sprintf("failed close repository %v", err2))
		}
	}()
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Print("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Msg(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	log.Print("Server exiting")
}

func (cmd *Command) Error() error {
	return cmd.Err
}

func HandleBootFlags(cmd *Command) *BootOptions {
	flags := cmd.Flags
	flag.Parse()

	// intercept version command
	if *flags.CmdPrintVersion {
		PrintVersion()
		os.Exit(0)
	}

	if *flags.CreateUserAdmin {
		load, err := Load(&BootOptions{
			WorkDir:   *flags.OptWorkDir,
			EnvPrefix: *flags.OptEnvPrefix,
		})
		if err != nil {
			panic(err)
		}

		cmd.CreateUserAdmin(load)
	}

	if *flags.UpdatePasswordUserAdmin {
		load, err := Load(&BootOptions{
			WorkDir:   *flags.OptWorkDir,
			EnvPrefix: *flags.OptEnvPrefix,
		})
		if err != nil {
			panic(err)
		}

		cmd.UpdatePasswordUserAdmin(load)
	}

	return &BootOptions{
		WorkDir:   *flags.OptWorkDir,
		EnvPrefix: *flags.OptEnvPrefix,
	}
}

func Load(bootOptions *BootOptions) (*config.Config, error) {
	cfg := new(config.Config)
	if err := envconfig.Process(bootOptions.EnvPrefix, cfg); err != nil {
		return nil, err
	}

	log.Warn().Int("port", int(cfg.ServerPort)).Msg("Running")

	if cfg.NodeId == "" {
		cfg.NodeId = uuid.New().String()
	}

	cfg.WorkDir = bootOptions.WorkDir

	return cfg, nil
}

func PrintVersion() {
	fmt.Printf("%s\n Version: %s\n Signature: %s\n", name, AppVersion, Build)
}

func (cmd *Command) CreateUserAdmin(cfg *config.Config) {
	repo, err := repositories.NewRepository(cfg)
	if err != nil {
		fmt.Printf("failed to read database configuration. Error = [%v]", err)
		os.Exit(0)
	}

	ctx := context.Background()

	rc, err := repo.Connected(ctx)
	if err != nil {
		fmt.Printf("failed to connect to database. Error = [%v]", err)
		os.Exit(0)
	}

	createdUserAdmin(rc, cfg)

	os.Exit(0)

}

func (cmd *Command) UpdatePasswordUserAdmin(cfg *config.Config) {
	repo, err := repositories.NewRepository(cfg)
	if err != nil {
		fmt.Printf("failed to read database configuration. Error = [%v]", err)
		os.Exit(0)
	}

	ctx := context.Background()

	rc, err := repo.Connected(ctx)
	if err != nil {
		fmt.Printf("failed to connect to database. Error = [%v]", err)
		os.Exit(0)
	}

	updatePasswordUserAdmin(rc, cfg)

	os.Exit(0)
}

func updatePasswordUserAdmin(rc *repositories.RepositoryContext, cfg *config.Config) string {

	return password

}

func createdUserAdmin(rc *repositories.RepositoryContext, cfg *config.Config) (string, string) {

	return "", ""
}
