package database

import "fmt"

type Config struct {
	Driver          string
	Host            string
	Port            string
	Username        string
	Password        string
	Database        string
	SslMode         string
	TimeZone        string
	MaxIdleConn     *int
	MaxOpenConn     *int
	MaxConnLifetime *int
	AppMode         string
}

func (c *Config) NormalizeValue() {
	// check if max idle conn is nil
	if c.MaxIdleConn == nil {
		c.MaxIdleConn = Integer(10)
	}

	// check if max open conn is nil
	if c.MaxOpenConn == nil {
		c.MaxOpenConn = Integer(100)
	}

	// check if max conn lifetime is nil
	if c.MaxConnLifetime == nil {
		c.MaxConnLifetime = Integer(10)
	}

	switch c.Driver {
	case "postgresql", "pg":
		c.Driver = DriverPostgresSql
	}
}

func (c *Config) DSN() (dsn string, err error) {
	switch c.Driver {
	case DriverPostgresSql:
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%v TimeZone=%s", c.Host, c.Username, c.Password, c.Database, c.Port, c.SslMode, c.TimeZone)
	default:
		err = fmt.Errorf("connection database unsupported driver '%s'", c.Driver)
	}

	return
}
