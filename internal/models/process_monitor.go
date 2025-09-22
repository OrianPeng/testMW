package models

// ProcessStepStatus 流程步骤状态
type ProcessStepStatus string

const (
	ProcessStepCompleted ProcessStepStatus = "completed" // 已完成
	ProcessStepCurrent   ProcessStepStatus = "current"   // 当前步骤
	ProcessStepPending   ProcessStepStatus = "pending"   // 待处理
	ProcessStepBlocked   ProcessStepStatus = "blocked"   // 被阻塞
)

// ProcessType 流程类型
type ProcessType string

const (
	ProcessTypePR ProcessType = "pr" // 采购请求
	ProcessTypePO ProcessType = "po" // 采购订单
)

// ProcessOverview 流程概览
type ProcessOverview struct {
	TodayPending    int     `json:"today_pending"`     // 今日待处理
	Overdue         int     `json:"overdue"`           // 逾期
	AvgApprovalTime float64 `json:"avg_approval_time"` // 平均审批时间（小时）
	CompletionRate  float64 `json:"completion_rate"`   // 完成率
}

// BlockedProcess 被阻塞的流程
type BlockedProcess struct {
	ID              string `json:"id"`
	RequestID       string `json:"request_id"`
	CurrentStep     string `json:"current_step"`
	Responsible     string `json:"responsible"`
	Department      string `json:"department"`
	BlockedDuration string `json:"blocked_duration"`
	Priority        int    `json:"priority"`
	Reason          string `json:"reason"`
}

// ProcessStep 流程步骤
type ProcessStep struct {
	Title       string            `json:"title"`
	Status      ProcessStepStatus `json:"status"`
	Responsible string            `json:"responsible"`
	Time        string            `json:"time"`
	Count       int               `json:"count"`
}

// ProcessStepsResponse 流程步骤响应
type ProcessStepsResponse struct {
	Steps []ProcessStep `json:"steps"`
}
