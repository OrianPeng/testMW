# 时间格式支持说明

## 概述

系统已经优化了时间解析功能，现在支持多种常用的日期时间格式。当您创建或更新采购请求时，可以使用以下任意一种格式来指定日期时间字段。

## 支持的日期时间格式

### 1. 标准日期格式
- `2026-05-19` - 标准日期格式 (YYYY-MM-DD)
- `2026/05/19` - 斜杠分隔日期格式 (YYYY/MM/DD)

### 2. 标准日期时间格式
- `2026-05-19T15:30:00` - ISO 8601 基本格式
- `2026-05-19 15:30:00` - 空格分隔格式
- `2026-05-19T15:30:00Z` - UTC 时间格式
- `2026-05-19T15:30:00.000Z` - 带毫秒的 UTC 格式

### 3. 带时区的格式
- `2026-05-19T15:30:00+08:00` - 带正时区偏移
- `2026-05-19T15:30:00-07:00` - 带负时区偏移
- `2026-05-19T15:30:00Z07:00` - RFC3339 标准格式

### 4. 简化时间格式
- `2026-05-19 15:30` - 只包含小时和分钟
- `2026-05-19T15:30` - T分隔符格式，只包含小时和分钟

### 5. 美式日期格式
- `05/19/2026` - 月/日/年格式 (MM/DD/YYYY)
- `19/05/2026` - 日/月/年格式 (DD/MM/YYYY)

## 使用示例

### 创建采购请求
```json
{
  "request_id": "PR-TEST-001",
  "doc_type": "PR",
  "plant": "1000",
  "quantity": 10,
  "unit_price": 25.50,
  "material": "TEST-MATERIAL",
  "delivery_date": "2026-05-19",
  "vendor_code": "VENDOR-TEST",
  "short_text": "Test purchase request",
  "material_group": "TEST-GROUP",
  "unit_type": "PCS",
  "requester": "Test User",
  "purchase_organization": "TEST-ORG",
  "currency": "CNY",
  "priority": 2,
  "urgency": "normal",
  "comments": "Test request"
}
```

### 更新采购请求
```json
{
  "delivery_date": "2026-05-19T15:30:00",
  "processed_at": "2026-05-20 10:00:00",
  "quantity": 20,
  "comments": "Updated via API"
}
```

### 使用不同格式的示例
```json
{
  "delivery_date": "2026/05/19",
  "processed_at": "2026-05-20T15:30:00Z"
}
```

```json
{
  "delivery_date": "05/19/2026",
  "processed_at": "2026-05-20T15:30:00+08:00"
}
```

## 特殊值处理

### 空字符串
- `"delivery_date": ""` - 将被解析为 null

### null 值
- `"processed_at": null` - 保持为 null

## 错误处理

如果提供的时间格式不在支持列表中，系统会返回详细的错误信息：

```json
{
  "error": "Invalid request format",
  "details": "failed to parse delivery_date: unable to parse time string 'invalid-date' with any supported format"
}
```

## 技术实现

系统使用 `parseTimeFlexible` 函数来处理时间解析，该函数会按顺序尝试所有支持的格式，直到找到匹配的格式。这确保了最大的兼容性，同时保持了良好的性能。

## 注意事项

1. **时区处理**: 如果时间字符串不包含时区信息，系统会使用本地时区
2. **性能考虑**: 解析过程会按顺序尝试格式，建议使用最常用的格式以获得最佳性能
3. **向后兼容**: 所有现有的 RFC3339 格式仍然完全支持
4. **数据验证**: 解析后的时间会进行有效性验证

## 推荐的格式

为了获得最佳性能和兼容性，建议使用以下格式：

- **日期**: `2026-05-19`
- **日期时间**: `2026-05-19T15:30:00`
- **带时区**: `2026-05-19T15:30:00Z` (UTC) 或 `2026-05-19T15:30:00+08:00` (特定时区) 