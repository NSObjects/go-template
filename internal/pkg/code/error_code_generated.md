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
| ErrUserNotFound | 200000 | 404 | User not found |
| ErrUserAlreadyExists | 200001 | 400 | User already exists |
| ErrUserInvalidData | 200002 | 400 | User invalid data |
| ErrUserPermissionDenied | 200003 | 403 | User permission denied |
| ErrUserInUse | 200004 | 400 | User is in use |
| ErrUserCreateFailed | 200005 | 500 | User create failed |
| ErrUserUpdateFailed | 200006 | 500 | User update failed |
| ErrUserDeleteFailed | 200007 | 500 | User delete failed |

