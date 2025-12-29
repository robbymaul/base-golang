package web

import (
	"paymentserviceklink/app/enums"
	"time"
)

type CreateConfigurationRequest struct {
	ConfigName  enums.ConfigName        `json:"configName"`
	ConfigValue enums.SandboxProduction `json:"configValue"`
	Aggregator  *AggregatorResponse     `json:"aggregator"`
	ConfigJson  *ConfigJson             `json:"configJson"`
}

type ConfigJson struct {
	SandboxBaseUrl               string `json:"sandboxBaseUrl,omitempty"`
	ProductionBaseUrl            string `json:"productionBaseUrl,omitempty"`
	SandboxMerchantId            string `json:"sandboxMerchantId,omitempty"`
	ProductionMerchantId         string `json:"productionMerchantId,omitempty"`
	SandboxMerchantCode          string `json:"sandboxMerchantCode,omitempty"`
	ProductionMerchantCode       string `json:"productionMerchantCode,omitempty"`
	SandboxMerchantName          string `json:"sandboxMerchantName,omitempty"`
	ProductionMerchantName       string `json:"productionMerchantName,omitempty"`
	SandboxApiKey                string `json:"sandboxApiKey,omitempty"`
	ProductionApiKey             string `json:"productionApiKey,omitempty"`
	SandboxServerKey             string `json:"sandboxServerKey,omitempty"`
	ProductionServerKey          string `json:"productionServerKey,omitempty"`
	SandboxSecretKey             string `json:"sandboxSecretKey,omitempty"`
	ProductionSecretKey          string `json:"productionSecretKey,omitempty"`
	SandboxClientKey             string `json:"sandboxClientKey,omitempty"`
	ProductionClientKey          string `json:"productionClientKey,omitempty"`
	SandboxSignatureKey          string `json:"sandboxSignatureKey,omitempty"`
	ProductionSignatureKey       string `json:"productionSignatureKey,omitempty"`
	SandboxCredentialPassword    string `json:"sandboxCredentialPassword,omitempty"`
	ProductionCredentialPassword string `json:"productionCredentialPassword,omitempty"`
	ReturnUrl                    string `json:"returnUrl"`
}

type ResponseConfiguration struct {
	Id           int64                   `json:"id"`
	AggregatorId int64                   `json:"aggregatorId"`
	Aggregator   *AggregatorResponse     `json:"aggregator,omitempty"`
	ConfigName   enums.ConfigName        `json:"configName"`
	ConfigValue  enums.SandboxProduction `json:"configValue"`
	IsActive     bool                    `json:"isActive"`
	ConfigJson   *ConfigJson             `json:"configJson"`
	CreatedAt    *time.Time              `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time              `json:"updatedAt,omitempty"`
}
