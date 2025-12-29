package models

import (
	"encoding/json"
	"gorm.io/datatypes"
	"paymentserviceklink/app/enums"
)

type AdminActivityLogs struct {
	Id           int64                          `gorm:"column:id;primaryKey;autoIncrement"`
	AdminUserId  int64                          `gorm:"column:admin_user_id;foreignKey"`
	Action       enums.ActionAdminActivityLog   `gorm:"column:action"`
	ResourceType enums.ResourceAdminActivityLog `gorm:"column:resource_type"`
	ResourceId   string                         `gorm:"column:resource_id"`
	Description  string                         `gorm:"column:description"`
	IpAddress    enums.StringEnum               `gorm:"column:ip_address"`
	UserAgent    string                         `gorm:"column:user_agent"`
	RequestData  datatypes.JSON                 `gorm:"column:request_data"`
	ResponseData datatypes.JSON                 `gorm:"column:response_data"`
	BaseField
}

func NewAdminActivityLogs(
	adminUserId int64,
	action enums.ActionAdminActivityLog,
	resourceType enums.ResourceAdminActivityLog,
	resourceId string,
	description string,
	ipAddress enums.StringEnum,
	userAgent string,
	requestData any,
	responseData any,
) *AdminActivityLogs {
	request, _ := json.Marshal(requestData)
	response, _ := json.Marshal(responseData)

	return &AdminActivityLogs{
		AdminUserId:  adminUserId,
		Action:       action,
		ResourceType: resourceType,
		ResourceId:   resourceId,
		Description:  description,
		IpAddress:    ipAddress,
		UserAgent:    userAgent,
		RequestData:  request,
		ResponseData: response,
	}
}

func (*AdminActivityLogs) TableName() string {
	return "admin_activity_logs"
}
