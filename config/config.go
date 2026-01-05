package config

type Config struct {
	/// Work directory path
	WorkDir     string `envconfig:"-"`
	Application string `envconfig:"APP_NAME"`

	// Gin Mode
	GinMode string `envconfig:"GIN_MODE"`
	AppMode string `envconfig:"APP_MODE"`

	// Basic Auth
	BasicAuthUsername string `envconfig:"BASIC_AUTH_USERNAME"`
	BasicAuthPassword string `envconfig:"BASIC_AUTH_PASSWORD"`

	// Database config migration
	DatabaseNameMigration string `envconfig:"MIGRATION_DB_NAME"`
	DatabaseUserMigration string `envconfig:"MIGRATION_DB_USER"`
	DatabasePassMigration string `envconfig:"MIGRATION_DB_PASSWORD"`
	DatabaseUpgradeOnBoot bool   `envconfig:"DB_BOOT_UPGRADE"`

	// Database config main
	DatabaseDSN             string `envconfig:"DATABASE_DSN"`
	DatabaseName            string `envconfig:"DATABASE_NAME"`
	DatabaseHost            string `envconfig:"DATABASE_HOST"`
	DatabasePort            string `envconfig:"DATABASE_PORT"`
	DatabaseUser            string `envconfig:"DATABASE_USER"`
	DatabasePass            string `envconfig:"DATABASE_PASS"`
	DatabaseTimezone        string `envconfig:"DATABASE_TIMEZONE"`
	DatabaseSslMode         string `envconfig:"DATABASE_SSL_MODE"`
	DatabaseMaxIdleConn     int    `envconfig:"DATABASE_MAX_IDLE_CONN"`
	DatabaseMaxConnLifetime int    `envconfig:"DATABASE_CONN_LIFETIME"`
	DatabaseOpenConn        int    `envconfig:"DATABASE_OPEN_CONN"`

	// SERVER
	ServerPort     uint16 `envconfig:"SERVER_PORT"`
	NodeId         string `envconfig:"NODE_ID"`
	ServerTimeZone string `envconfig:"TZ"`

	// LOGGING
	LogLevel string `envconfig:"LOG_LEVEL"`

	// JWT
	JwtSecret string `envconfig:"JWT_SECRET"`
	JwtExpire int64  `envconfig:"JWT_EXPIRE"`
	JwtIssuer string `envconfig:"JWT_ISSUER"`
}
