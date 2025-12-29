package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"paymentserviceklink/app/enums"
)

type Configuration struct {
	Id           int64                   `gorm:"column:id;primaryKey;autoIncrement"`
	AggregatorId int64                   `gorm:"column:aggregator_id;foreignKey"`
	ConfigKey    string                  `gorm:"column:config_key"`
	ConfigValue  enums.SandboxProduction `gorm:"column:config_value"`
	ConfigName   enums.ConfigName        `gorm:"column:config_name"`
	ConfigJson   ConfigJson              `gorm:"column:config_json"`
	IsActive     bool                    `gorm:"column:is_active"`
	Aggregator   *Aggregator             `gorm:"foreignKey:AggregatorId;references:id;"`
	BaseField
}

func (*Configuration) TableName() string {
	return "configurations"
}

type ConfigJson struct {
	SandboxBaseUrl               string `json:"sandbox_base_url"`
	ProductionBaseUrl            string `json:"production_base_url"`
	SandboxMerchantId            string `json:"sandbox_merchant_id"`
	ProductionMerchantId         string `json:"production_merchant_id"`
	SandboxMerchantCode          string `json:"sandbox_merchant_code"`
	ProductionMerchantCode       string `json:"production_merchant_code"`
	SandboxMerchantName          string `json:"sandbox_merchant_name"`
	ProductionMerchantName       string `json:"production_merchant_name"`
	SandboxApiKey                string `json:"sandbox_api_key"`
	ProductionApiKey             string `json:"production_api_key"`
	SandboxServerKey             string `json:"sandbox_server_key"`
	ProductionServerKey          string `json:"production_server_key"`
	SandboxSecretKey             string `json:"sandbox_secret_key"`
	ProductionSecretKey          string `json:"production_secret_key"`
	SandboxClientKey             string `json:"sandbox_client_key"`
	ProductionClientKey          string `json:"production_client_key"`
	SandboxSignatureKey          string `json:"sandbox_signature_key"`
	ProductionSignatureKey       string `json:"production_signature_key"`
	SandboxCredentialPassword    string `json:"sandbox_credential_password"`
	ProductionCredentialPassword string `json:"production_credential_password"`
	PublicKey                    string `json:"public_key"`
	PrivateKey                   string `json:"private_key"`
	ReturnUrl                    string `json:"return_url"`
	//Subscription           bool   `json:"subscription"`
}

// Implement driver.Valuer
func (c *ConfigJson) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Implement sql.Scanner
func (c *ConfigJson) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, c)
}

func AllowedFilterColumnConfiguration() map[string]FilterColumn {
	return map[string]FilterColumn{
		"platform_id": {
			Operator: []string{"eq"},
			Variant:  "number",
			Table:    "platform_configuration",
		},
		"slug": {
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "aggregators",
		},
		"currency": {
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "aggregators",
		},
		"config_name": {
			Operator: []string{"eq", "like", "notLike"},
			Variant:  "string",
			Table:    "configurations",
		},
		"is_active": {
			Operator: []string{"eq"},
			Variant:  "boolean",
			Table:    "configurations",
		},
	}
}
