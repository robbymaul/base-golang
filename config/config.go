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

	// SENANGPAY
	SenangpayIsProduction bool   `envconfig:"SENANGPAY_IS_PRODUCTION"`
	SenangpayUrl          string `envconfig:"SENANGPAY_URL"`
	SenangpaySecretKey    string `envconfig:"SENANGPAY_SECRET_KEY"`
	SenangpayMerchantId   string `envconfig:"SENANGPAY_MERCHANT_ID"`

	// MIDTRANS
	MidtransIsProduction        bool   `envconfig:"MIDTRANS_IS_PRODUCTION"`
	MidtransBaseUrlProduction   string `envconfig:"MIDTRANS_BASE_URL_PRODUCTION"`
	MidtransServerKeyProduction string `envconfig:"MIDTRANS_SERVER_KEY_PRODUCTION"`
	MidtransClientKeyProduction string `envconfig:"MIDTRANS_CLIENT_KEY_PRODUCTION"`
	MidtransBaseUrlSandbox      string `envconfig:"MIDTRANS_BASE_URL_SANDBOX"`
	MidtransServerKeySandbox    string `envconfig:"MIDTRANS_SERVER_KEY_SANDBOX"`
	MidtransClientKeySandbox    string `envconfig:"MIDTRANS_CLIENT_KEY_SANDBOX"`

	// ESPAY
	EspayIsProduction            bool   `envconfig:"ESPAY_IS_PRODUCTION"`
	EspayTopupBaseUrl            string `envconfig:"ESPAY_TOPUP_BASE_URL"`
	EspayTopupMerchantCode       string `envconfig:"ESPAY_TOPUP_MERCHANT_CODE"`
	EspayTopupMerchantName       string `envconfig:"ESPAY_TOPUP_MERCHANT_NAME"`
	EspayTopupApiKey             string `envconfig:"ESPAY_TOPUP_API_KEY"`
	EspayTopupSignatureKey       string `envconfig:"ESPAY_TOPUP_SIGNATURE_KEY"`
	EspayTopupCredentialPassword string `envconfig:"ESPAY_TOPUP_CREDENTIAL_PASSWORD"`
	EspayTopupPublicKey          string `envconfig:"ESPAY_TOPUP_PUBLIC_KEY"`
	EspayTopupPrivateKey         string `envconfig:"ESPAY_TOPUP_PRIVATE_KEY"`
	EspayTopupReturnUrl          string `envconfig:"ESPAY_TOPUP_RETURN_URL"`

	// S3
	S3AccessKeyId       string `envconfig:"S3_ACCESS_KEY_ID"`
	S3BucketName        string `envconfig:"S3_BUCKET_NAME"`
	S3Endpoint          string `envconfig:"S3_ENDPOINT_NAME"`
	S3Region            string `envconfig:"S3_REGION"`
	S3SecretAccessKey   string `envconfig:"S3_SECRET_ACCESS_KEY"`
	S3UploadUrlLifetime uint   `envconfig:"S3_UPLOAD_URL_LIFETIME"`
	S3UseSSL            bool   `envconfig:"S3_USE_SSL"`
	S3AssetLogo         string `envconfig:"S3_ASSET_LOGO"`
}
