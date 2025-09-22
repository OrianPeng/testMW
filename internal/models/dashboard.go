package models

// DashboardMetrics 仪表板指标
type DashboardMetrics struct {
	TodayPR         int     `json:"today_pr"`          // 今日采购请求数
	TodayPO         int     `json:"today_po"`          // 今日采购订单数
	PendingApproval int     `json:"pending_approval"`  // 待审批数
	AvgApprovalTime float64 `json:"avg_approval_time"` // 平均审批时间
}

// PRTrendData PR趋势数据
type PRTrendData struct {
	Dates  []string `json:"dates"`
	Values []int    `json:"values"`
}

// POStatusDistribution PO状态分布
type POStatusDistribution struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// SupplierRanking 供应商排名
type SupplierRanking struct {
	Supplier    string  `json:"supplier"`
	OrderCount  int     `json:"order_count"`
	TotalAmount float64 `json:"total_amount"`
}

// DepartmentStats 部门统计
type DepartmentStats struct {
	Department string `json:"department"`
	PRCount    int    `json:"pr_count"`
}

// DashboardQuery 仪表板查询条件
type DashboardQuery struct {
	Period string `json:"period"` // 时间周期：7d, 30d, 90d
}
