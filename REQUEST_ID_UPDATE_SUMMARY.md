# 采购请求API更新：从MySQL ID改为Request ID

## 概述

根据用户要求，将采购请求的更新API从使用MySQL的`id`字段改为使用`request_id`字段。这个修改使得API更加直观和易用，因为`request_id`通常包含业务含义，比数字ID更容易理解和记忆。

## 主要修改

### 1. 路由配置修改

**文件**: `internal/api/router.go`

**修改前**:
```go
purchase.GET("/:id", purchaseHandler.GetPurchaseRequest)
purchase.PUT("/:id", purchaseHandler.UpdatePurchaseRequest)
purchase.DELETE("/:id", purchaseHandler.DeletePurchaseRequest)
```

**修改后**:
```go
purchase.GET("/:request_id", purchaseHandler.GetPurchaseRequest)
purchase.PUT("/:request_id", purchaseHandler.UpdatePurchaseRequest)
purchase.DELETE("/:request_id", purchaseHandler.DeletePurchaseRequest)
```

### 2. API处理器修改

**文件**: `internal/api/purchase_request_handlers.go`

#### GetPurchaseRequest方法
- 将参数从`c.Param("id")`改为`c.Param("request_id")`
- 移除数字ID解析逻辑
- 直接使用`GetByRequestID`方法查询

#### UpdatePurchaseRequest方法
- 将参数从`c.Param("id")`改为`c.Param("request_id")`
- 移除数字ID解析和转换逻辑
- 使用`UpdateByRequestID`方法进行更新
- 更新日志记录，使用request_id而不是数字ID

#### DeletePurchaseRequest方法
- 将参数从`c.Param("id")`改为`c.Param("request_id")`
- 移除数字ID解析逻辑
- 使用`DeleteByRequestID`方法进行删除

### 3. Repository层新增方法

**文件**: `internal/repository/purchase_request_repository.go`

#### 新增UpdateByRequestID方法
```go
func (r *PurchaseRequestRepository) UpdateByRequestID(ctx context.Context, requestID string, updateReq *models.UpdatePurchaseRequestRequest) (*models.PurchaseRequest, error) {
    // 首先通过request_id获取记录
    purchaseReq, err := r.GetByRequestID(ctx, requestID)
    if err != nil {
        return nil, fmt.Errorf("purchase request not found with request_id: %s", requestID)
    }

    // 使用现有的Update方法进行更新
    return r.Update(ctx, purchaseReq.ID, updateReq)
}
```

#### 新增DeleteByRequestID方法
```go
func (r *PurchaseRequestRepository) DeleteByRequestID(ctx context.Context, requestID string) error {
    // 首先通过request_id获取记录
    purchaseReq, err := r.GetByRequestID(ctx, requestID)
    if err != nil {
        return fmt.Errorf("purchase request not found with request_id: %s", requestID)
    }

    // 使用现有的Delete方法进行删除
    return r.Delete(ctx, purchaseReq.ID)
}
```

### 4. 测试文件修改

**文件**: `cmd/middle_test/main.go`

- 将提取MySQL ID的逻辑改为提取request_id
- 更新所有测试用例，使用request_id而不是数字ID
- 更新错误测试用例，使用不存在的request_id字符串
- 重新编号测试用例（从测试11到测试21）

**文件**: `test_update.bat`

- 将所有curl命令中的数字ID改为request_id格式
- 更新测试用例说明

### 5. 文档更新

**文件**: `UPDATE_API_README.md`

- 更新API URL格式：从`{id}`改为`{request_id}`
- 更新参数说明：从数字ID改为request_id字符串
- 更新所有示例代码，使用request_id格式
- 更新特性说明，强调request_id的优势

**文件**: `UPDATE_IMPLEMENTATION_SUMMARY.md`

- 更新实现总结，反映使用request_id的变化
- 添加新增的Repository方法说明
- 更新主要变更列表
- 添加优势说明

## API使用示例

### 修改前（使用MySQL ID）
```bash
# 获取采购请求
curl -X GET http://localhost:8080/api/v1/purchase-requests/1

# 更新采购请求
curl -X PUT http://localhost:8080/api/v1/purchase-requests/1 \
  -H "Content-Type: application/json" \
  -d '{"quantity": 100, "comments": "Updated"}'

# 删除采购请求
curl -X DELETE http://localhost:8080/api/v1/purchase-requests/1
```

### 修改后（使用Request ID）
```bash
# 获取采购请求
curl -X GET http://localhost:8080/api/v1/purchase-requests/PR-123456

# 更新采购请求
curl -X PUT http://localhost:8080/api/v1/purchase-requests/PR-123456 \
  -H "Content-Type: application/json" \
  -d '{"quantity": 100, "comments": "Updated"}'

# 删除采购请求
curl -X DELETE http://localhost:8080/api/v1/purchase-requests/PR-123456
```

## 优势

1. **更直观**: request_id通常包含业务含义，比数字ID更容易理解和记忆
2. **更安全**: 避免数字ID可能带来的安全问题（如ID枚举攻击）
3. **更灵活**: request_id可以包含业务信息，如"PR-2024-001"表示2024年第1个采购请求
4. **更一致**: 与其他业务系统保持一致，通常业务系统都使用有意义的标识符
5. **更易调试**: 在日志和错误信息中，request_id比数字ID更容易识别和追踪

## 兼容性

这个修改是向后不兼容的，因为：
- API路径参数从`:id`改为`:request_id`
- 不再支持数字ID作为参数
- 客户端需要更新为使用request_id

## 测试验证

创建了测试文件`test_request_id_update.go`来验证修改的正确性，该测试：
1. 创建一个新的采购请求
2. 提取返回的request_id
3. 使用request_id进行更新操作
4. 验证更新结果

## 总结

这次修改成功将采购请求API从使用MySQL的`id`字段改为使用`request_id`字段，使得API更加直观、安全和易用。所有相关的代码、测试和文档都已更新，确保系统的一致性和完整性。 