# 时间格式优化总结

## 问题描述

用户在使用采购请求API时遇到了时间解析错误：
```
{
    "details": "parsing time \"2026-05-19\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \"T\"",
    "error": "Invalid request format"
}
```

这个错误表明系统只支持RFC3339格式的时间字符串，而用户使用的是简单的日期格式 `"2026-05-19"`。

## 解决方案

### 1. 创建灵活的时间解析函数

在 `internal/models/purchase_request.go` 中添加了 `parseTimeFlexible` 函数，支持多种常用的日期时间格式：

```go
func parseTimeFlexible(timeStr string) (time.Time, error) {
    // 支持的日期时间格式列表
    formats := []string{
        time.RFC3339,                    // "2006-01-02T15:04:05Z07:00"
        "2006-01-02T15:04:05.000Z",     // "2006-01-02T15:04:05.000Z"
        "2006-01-02T15:04:05",          // "2006-01-02T15:04:05"
        "2006-01-02 15:04:05",          // "2006-01-02 15:04:05"
        "2006-01-02",                   // "2006-01-02"
        "2006/01/02",                   // "2006/01/02"
        "2006-01-02T15:04:05Z",         // "2006-01-02T15:04:05Z"
        "2006-01-02T15:04:05-07:00",    // "2006-01-02T15:04:05-07:00"
        "2006-01-02T15:04:05+07:00",    // "2006-01-02T15:04:05+07:00"
        "01/02/2006",                   // "01/02/2006" (MM/DD/YYYY)
        "02/01/2006",                   // "02/01/2006" (DD/MM/YYYY)
        "2006-01-02 15:04",             // "2006-01-02 15:04"
        "2006-01-02T15:04",             // "2006-01-02T15:04"
    }
    
    // 尝试每种格式
    for _, format := range formats {
        if t, err := time.Parse(format, timeStr); err == nil {
            return t, nil
        }
    }
    
    return time.Time{}, fmt.Errorf("unable to parse time string '%s' with any supported format", timeStr)
}
```

### 2. 更新CreatePurchaseRequestRequest的JSON反序列化

修改了 `CreatePurchaseRequestRequest` 的 `UnmarshalJSON` 方法，使用新的时间解析函数：

```go
func (r *CreatePurchaseRequestRequest) UnmarshalJSON(data []byte) error {
    // ... 现有代码 ...
    t, err := parseTimeFlexible(v)
    if err != nil {
        return fmt.Errorf("failed to parse delivery_date: %w", err)
    }
    r.DeliveryDate = &t
    // ... 现有代码 ...
}
```

### 3. 为UpdatePurchaseRequestRequest添加JSON反序列化

为 `UpdatePurchaseRequestRequest` 结构体添加了自定义的 `UnmarshalJSON` 方法，支持 `delivery_date` 和 `processed_at` 字段的灵活时间解析：

```go
func (r *UpdatePurchaseRequestRequest) UnmarshalJSON(data []byte) error {
    type Alias UpdatePurchaseRequestRequest
    tmp := &struct {
        DeliveryDate interface{} `json:"delivery_date"`
        ProcessedAt  interface{} `json:"processed_at"`
        *Alias
    }{
        Alias: (*Alias)(r),
    }
    
    // 处理delivery_date和processed_at字段
    // 使用parseTimeFlexible函数解析时间
    // ...
}
```

## 支持的格式列表

### 标准日期格式
- `2026-05-19` - 标准日期格式 (YYYY-MM-DD)
- `2026/05/19` - 斜杠分隔日期格式 (YYYY/MM/DD)

### 标准日期时间格式
- `2026-05-19T15:30:00` - ISO 8601 基本格式
- `2026-05-19 15:30:00` - 空格分隔格式
- `2026-05-19T15:30:00Z` - UTC 时间格式
- `2026-05-19T15:30:00.000Z` - 带毫秒的 UTC 格式

### 带时区的格式
- `2026-05-19T15:30:00+08:00` - 带正时区偏移
- `2026-05-19T15:30:00-07:00` - 带负时区偏移
- `2026-05-19T15:30:00Z07:00` - RFC3339 标准格式

### 简化时间格式
- `2026-05-19 15:30` - 只包含小时和分钟
- `2026-05-19T15:30` - T分隔符格式，只包含小时和分钟

### 美式日期格式
- `05/19/2026` - 月/日/年格式 (MM/DD/YYYY)
- `19/05/2026` - 日/月/年格式 (DD/MM/YYYY)

## 特殊值处理

- **空字符串**: `"delivery_date": ""` 将被解析为 null
- **null 值**: `"processed_at": null` 保持为 null

## 错误处理改进

现在当时间格式无法解析时，会返回更详细的错误信息：

```json
{
  "error": "Invalid request format",
  "details": "failed to parse delivery_date: unable to parse time string 'invalid-date' with any supported format"
}
```

## 测试验证

### 1. 更新了测试文件
在 `cmd/middle_test/main.go` 中添加了时间格式测试：

```go
// 测试不同的时间格式
timeFormats := []string{
    "2026-05-19",                    // 日期格式
    "2026-05-19T15:30:00",           // 日期时间格式
    "2026-05-19 15:30:00",           // 空格分隔格式
    "2026/05/19",                    // 斜杠分隔格式
    "2026-05-19T15:30:00Z",          // UTC格式
    "2026-05-19T15:30:00+08:00",     // 带时区格式
}
```

### 2. 创建了测试脚本
创建了 `test_time_formats.bat` 批处理脚本，用于快速测试不同的时间格式。

### 3. 创建了文档
- `TIME_FORMAT_SUPPORT.md` - 详细的时间格式支持说明
- `TIME_FORMAT_OPTIMIZATION_SUMMARY.md` - 本总结文档

## 向后兼容性

- 所有现有的 RFC3339 格式仍然完全支持
- 现有的API接口保持不变
- 数据库存储格式保持不变

## 性能考虑

- 解析过程会按顺序尝试格式，建议使用最常用的格式以获得最佳性能
- 推荐的格式：
  - **日期**: `2026-05-19`
  - **日期时间**: `2026-05-19T15:30:00`
  - **带时区**: `2026-05-19T15:30:00Z` (UTC)

## 使用示例

### 创建采购请求
```json
{
  "request_id": "PR-TEST-001",
  "delivery_date": "2026-05-19",
  "material": "TEST-MATERIAL",
  // ... 其他字段
}
```

### 更新采购请求
```json
{
  "delivery_date": "2026-05-19T15:30:00",
  "processed_at": "2026-05-20 10:00:00",
  "quantity": 20
}
```

## 总结

通过这次优化，系统现在支持13种不同的时间格式，大大提高了API的易用性和兼容性。用户可以使用他们习惯的日期格式，而不需要强制转换为特定的RFC3339格式。同时保持了向后兼容性，确保现有功能不受影响。 