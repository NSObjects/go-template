# 错误码文档

> 系统错误码列表，由 `codegen -type=int -doc` 命令生成，不要对此文件做任何更改。

## 功能说明

如果返回结果中存在 `code` 字段，则表示调用 API 接口失败。例如：

```json
{
  "code": 100101,
  "message": "Database error",
  "request_id": "20231201120000-abc12345",
  "timestamp": 1701234567
}
```

上述返回中：
- `code` 表示业务错误码
- `message` 表示该错误的具体信息
- `request_id` 表示请求ID，用于错误追踪
- `timestamp` 表示错误发生的时间戳

每个错误同时也对应一个 HTTP 状态码，比如上述错误码对应了 HTTP 状态码 500(Internal Server Error)。

## 错误码分类

- **1xxxxx**: 通用错误
- **100001-100099**: 基本错误
- **100101-100199**: 数据库错误
- **100201-100299**: 认证授权错误
- **100301-100399**: 编解码错误
- **100401-100499**: 验证错误
- **100501-100599**: 菜单相关错误
- **100601-100699**: 用户相关错误
- **100701-100799**: 角色相关错误

## 错误码列表

系统支持的错误码列表如下：

| Identifier | Code | HTTP Code | Description |
| ---------- | ---- | --------- | ----------- |
| ErrAuth1020201 | 1020201 | 401 | 刷新令牌无效或已过期 |
| ErrAuth1020102 | 1020102 | 400 | 刷新令牌格式错误 |
| ErrAuth1010201 | 1010201 | 401 | 用户名或密码错误 |
| ErrAuth1010602 | 1010602 | 400 | 账号已禁用 |
| ErrAuth1010603 | 1010603 | 400 | 租户已停用，禁止登录 |
| ErrAuth1010804 | 1010804 | 400 | 登录失败次数过多，请稍后重试 |
| ErrAuth1010105 | 1010105 | 400 | 缺少 device_id |
| ErrAuth1030201 | 1030201 | 401 | 未登录或会话已失效 |
| ErrSuccess | 100001 | 200 | OK |
| ErrUnknown | 100002 | 500 | Internal server error |
| ErrBind | 100003 | 400 | Error occurred while binding the request body to the struct |
| ErrValidation | 100004 | 400 | Validation failed |
| ErrTokenInvalid | 100005 | 401 | Token invalid |
| ErrDatabase | 100101 | 500 | Database error |
| ErrRedis | 100102 | 500 | Redis error |
| ErrKafka | 100103 | 500 | Kafka error |
| ErrExternalService | 100104 | 500 | External service error |
| ErrBadRequest | 100400 | 400 | Bad request |
| ErrUnauthorized | 100401 | 401 | Unauthorized |
| ErrForbidden | 100402 | 403 | Forbidden |
| ErrNotFound | 100403 | 404 | Not found |
| ErrInternalServer | 100404 | 500 | Internal server error |
| ErrEncrypt | 100201 | 401 | Error occurred while encrypting the user password |
| ErrSignatureInvalid | 100202 | 401 | Signature is invalid |
| ErrExpired | 100203 | 401 | Token expired |
| ErrInvalidAuthHeader | 100204 | 401 | Invalid authorization header |
| ErrMissingHeader | 100205 | 401 | The `Authorization` header was empty |
| ErrPasswordIncorrect | 100206 | 401 | Password was incorrect |
| ErrPermissionDenied | 100207 | 403 | Permission denied |
| ErrAccountLocked | 100208 | 403 | Account is locked |
| ErrAccountDisabled | 100209 | 403 | Account is disabled |
| ErrTooManyAttempts | 100210 | 403 | Too many login attempts |
| ErrEncodingFailed | 100301 | 500 | Encoding failed due to an error with the data |
| ErrDecodingFailed | 100302 | 500 | Decoding failed due to an error with the data |
| ErrInvalidJSON | 100303 | 500 | Data is not valid JSON |
| ErrEncodingJSON | 100304 | 500 | JSON data could not be encoded |
| ErrDecodingJSON | 100305 | 500 | JSON data could not be decoded |
| ErrInvalidYaml | 100306 | 500 | Data is not valid Yaml |
| ErrEncodingYaml | 100307 | 500 | Yaml data could not be encoded |
| ErrDecodingYaml | 100308 | 500 | Yaml data could not be decoded |
| ErrBilling15010601 | 15010601 | 400 | 该账期账单已生成 |
| ErrInbound8010101 | 8010101 | 400 | 运单号列表为空 |
| ErrInbound8010102 | 8010102 | 400 | 入库明细为空 |
| ErrInbound8010403 | 8010403 | 404 | SKU 不存在 |
| ErrInbound8030604 | 8030604 | 400 | 入库单状态非 receiving，不可新增批次 |
| ErrInbound8030405 | 8030405 | 404 | 库位不存在 |
| ErrInbound8030106 | 8030106 | 400 | 收货件数必须为正数 |
| ErrIntegrations13010501 | 13010501 | 400 | 该店铺已绑定 |
| ErrIntegrations13010702 | 13010702 | 500 | 授权码无效或过期 |
| ErrIntegrations13030503 | 13030503 | 400 | SKU 映射已存在 |
| ErrIntegrations13030404 | 13030404 | 404 | 本地SKU不存在 |
| ErrInventory9010101 | 9010101 | 400 | 分页参数不合法 |
| ErrLocations6020501 | 6020501 | 400 | 库位编码在仓库内已存在 |
| ErrLocations6020104 | 6020104 | 400 | 库位体积填写不合法 |
| ErrMe2020101 | 2020101 | 400 | 新密码不符合安全策略 |
| ErrMe2020202 | 2020202 | 401 | 旧密码错误 |
| ErrOrders10010501 | 10010501 | 400 | 平台订单号已存在 |
| ErrOrders10010102 | 10010102 | 400 | 订单明细为空 |
| ErrOrders10010903 | 10010903 | 400 | 平台SKU未映射到本地SKU |
| ErrRates14010601 | 14010601 | 400 | 存在已生效/待生效版本，禁止重叠 |
| ErrRates14010102 | 14010102 | 400 | 费率参数不合法 |
| ErrReconNotFound | 100000 | 404 | Recon not found |
| ErrReconAlreadyExists | 100001 | 400 | Recon already exists |
| ErrReconInvalidData | 100002 | 400 | Recon invalid data |
| ErrReconPermissionDenied | 100003 | 403 | Recon permission denied |
| ErrReconInUse | 100004 | 400 | Recon is in use |
| ErrReconCreateFailed | 100005 | 500 | Recon create failed |
| ErrReconUpdateFailed | 100006 | 500 | Recon update failed |
| ErrReconDeleteFailed | 100007 | 500 | Recon delete failed |
| ErrReservations11010601 | 11010601 | 400 | 库存不足，无法预占 |
| ErrReservations11010602 | 11010602 | 400 | 订单状态不允许预占 |
| ErrReservations11010503 | 11010503 | 400 | 该订单与SKU已存在预占记录 |
| ErrRoles4020501 | 4020501 | 400 | 角色编码已存在 |
| ErrShipments12010501 | 12010501 | 400 | 运单号已存在 |
| ErrShipments12010602 | 12010602 | 400 | 订单未完成预占或预占不足 |
| ErrSku7020501 | 7020501 | 400 | 商品编码已存在 |
| ErrSku7020103 | 7020103 | 400 | 包装规格不合法 |
| ErrSystem17010101 | 17010101 | 400 | IP/CIDR 格式不合法 |
| ErrSystem17010502 | 17010502 | 400 | IP 白名单已存在 |
| ErrUsers3020501 | 3020501 | 400 | 登录名已存在 |
| ErrUsers3020502 | 3020502 | 400 | 邮箱已存在 |
| ErrUsers3020503 | 3020503 | 400 | 手机号已存在 |
| ErrUsers3020104 | 3020104 | 400 | 角色代码无效 |
| ErrUsers4030303 | 4030303 | 403 | 无权绑定该角色范围 |
| ErrUsers3010401 | 3010401 | 404 | 用户不存在 |
| ErrWarehouses5020501 | 5020501 | 400 | 仓库名称已存在 |

