package models

import "time"

// UiPathConfig UiPath配置
type UiPathConfig struct {
	OrchBaseURL string        `json:"orch_base_url" yaml:"orch_base_url"`
	TenancyName string        `json:"tenancy_name" yaml:"tenancy_name"`
	Username    string        `json:"username" yaml:"username"`
	Password    string        `json:"password" yaml:"password"`
	FolderID    int           `json:"folder_id" yaml:"folder_id"`
	QueueName   string        `json:"queue_name" yaml:"queue_name"`
	VerifySSL   bool          `json:"verify_ssl" yaml:"verify_ssl"`
	Timeout     time.Duration `json:"timeout" yaml:"timeout"`
}

// UiPathAuthRequest UiPath认证请求
type UiPathAuthRequest struct {
	TenancyName            string `json:"tenancyName"`
	UsernameOrEmailAddress string `json:"usernameOrEmailAddress"`
	Password               string `json:"password"`
}

// UiPathAuthResponse UiPath认证响应
type UiPathAuthResponse struct {
	Result string `json:"result"`
}

// UiPathQueueItemRequest UiPath队列项目请求
type UiPathQueueItemRequest struct {
	ItemData UiPathItemData `json:"itemData"`
}

// UiPathItemData UiPath队列项目数据
type UiPathItemData struct {
	Name            string                 `json:"Name"`
	Priority        string                 `json:"Priority"`
	SpecificContent map[string]interface{} `json:"SpecificContent"`
	Reference       string                 `json:"Reference"`
}

// UiPathQueueItemResponse UiPath队列项目响应
type UiPathQueueItemResponse struct {
	ID int `json:"Id"`
}

// UiPathQueueItemRequest UiPath队列项目请求（用于API）
type UiPathAddQueueItemRequest struct {
	QueueName       string                 `json:"queue_name" binding:"required"`
	Priority        string                 `json:"priority"` // Normal, High, Critical
	SpecificContent map[string]interface{} `json:"specific_content" binding:"required"`
	Reference       string                 `json:"reference"`
}

// UiPathQueueItemResponse UiPath队列项目响应（用于API）
type UiPathAddQueueItemResponse struct {
	ID        int       `json:"id"`
	QueueName string    `json:"queue_name"`
	Priority  string    `json:"priority"`
	Reference string    `json:"reference"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// UiPathStatusResponse UiPath系统状态响应
type UiPathStatusResponse struct {
	IsAvailable bool      `json:"is_available"`
	LastCheck   time.Time `json:"last_check"`
	Message     string    `json:"message,omitempty"`
}
