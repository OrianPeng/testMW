package models

// AIResponseType AI响应类型
type AIResponseType string

const (
	AIResponseTypeText    AIResponseType = "text"    // 文本响应
	AIResponseTypeTable   AIResponseType = "table"   // 表格响应
	AIResponseTypeCard    AIResponseType = "card"    // 卡片响应
	AIResponseTypeProcess AIResponseType = "process" // 流程响应
)

// AIMessageRequest AI消息请求
type AIMessageRequest struct {
	Message        string                 `json:"message" binding:"required"`
	ConversationID string                 `json:"conversation_id"`
	Context        map[string]interface{} `json:"context"`
}

// AITextResponse AI文本响应
type AITextResponse struct {
	Type    AIResponseType `json:"type"`
	Content string         `json:"content"`
}

// AITableColumn AI表格列
type AITableColumn struct {
	Title     string `json:"title"`
	DataIndex string `json:"dataIndex"`
	Key       string `json:"key"`
}

// AITableData AI表格数据
type AITableData struct {
	Key  string                 `json:"key"`
	Data map[string]interface{} `json:"data"`
}

// AITableResponse AI表格响应
type AITableResponse struct {
	Type         AIResponseType  `json:"type"`
	TableColumns []AITableColumn `json:"table_columns"`
	TableData    []AITableData   `json:"table_data"`
}

// AICardResponse AI卡片响应
type AICardResponse struct {
	Type      AIResponseType    `json:"type"`
	CardTitle string            `json:"card_title"`
	CardData  map[string]string `json:"card_data"`
}

// AIProcessStep AI流程步骤
type AIProcessStep struct {
	Title       string `json:"title"`
	Status      string `json:"status"`
	Responsible string `json:"responsible"`
}

// AIProcessResponse AI流程响应
type AIProcessResponse struct {
	Type         AIResponseType  `json:"type"`
	ProcessSteps []AIProcessStep `json:"process_steps"`
}

// AIMessageResponse AI消息响应
type AIMessageResponse struct {
	Response       interface{}    `json:"response"`
	Type           AIResponseType `json:"type"`
	ConversationID string         `json:"conversation_id"`
	Suggestions    []string       `json:"suggestions"`
}
