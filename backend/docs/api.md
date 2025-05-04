# ToRead API 文档

ligen131 [i@lg.ee](mailto:i@lg.ee)

## 总览

Link: <https://to-read.lg.gl/api/v1>

+ 1 [总览](#总览)
+ 2 [Health](#health)
  + 2.1 [[GET] `/health`](#get-health)
+ 3 [用户 User](#用户-user)
  + 3.1 [[GET] `/user`](#get-user)
  + 3.2 [[POST] `/user/register`](#post-userregister)
  + 3.3 [[POST] `/user/login`](#post-userlogin)
+ 4 [链接收藏 Collection](#链接收藏-collection)
  + 4.1 [*[GET] `/collection/list`](#get-collectionlist)
  + 4.2 [*[POST] `/collection/add`](#post-collectionadd)
  + 4.3 [*[GET] `/collection/summary`](#get-collectionsummary)
  + 4.4 [*[GET] `/collection/tag`](#get-collectiontag)

在标题带 `*` 标识的请求中，请在请求头中提供登录获取到的 JWT token。

**所有 GET 都使用 QueryString 格式而非 JSON Body。**

```yaml
Authorization: Bearer <token>
```

## Health

### [GET] `/health`

获取服务状态。

#### Request

无。

#### Response

```json
{
  "code": 200,
  "msg": null,
  "data": "ok",
}
```

## 用户 User

- `user_id`: 用户注册时由后端生成的用户 ID，递增整数。
- `user_name`: 用户自定义昵称，字符串。
- `role`: 用户角色，整数，目前只有一种角色。

### [GET] `/user`

#### Request

```json
{
  "user_id": 1,
  "user_name": "ligen131"
}
```

`user_id` 和 `user_name` 二选一，若都提供则优先使用 `user_id`。

#### Response

```json
{
  "code": 200,
  "msg": null,
  "data": {
    "user_id": 1,
    "user_name": "ligen131",
    "role": 1
  },
}
```

### [POST] `/user/register`

用户注册。

#### Request

```json
{
  "user_name": "ligen131",
  "password": "xxxxxx",
}
```

- `user_name`: 用户自定义昵称。
- `password`: 用户密码。

#### Response

```json
{
  "code": 200,
  "msg": null,
  "data": {
    "user_id": 1,
    "user_name": "ligen131",
    "role": 1
  }
}
```

若注册成功，返回用户信息。若失败，在 `data.msg` 中返回错误信息。

### [POST] `/user/login`

用户登录系统。

#### Request

```json
{
  "user_name": "ligen131",
  "password": "xxxxxx",
}
```

#### Response

```json
{
  "code": 200,
  "msg": null,
  "data": {
    "user_id": 1,
    "user_name": "ligen131",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxM...",
    "token_expiration_time": 1683561600,
  },
}
```

若登录成功，返回用户信息。若失败，在 `data.msg` 中返回错误信息。

- `token`: 登录时获取 JWT token，请在与用户权限相关的请求发送时在请求头中包含该 token。

  ```yaml
  Authorization: Bearer <token>
  ```
- `token_expiration_time`: token 过期时间，格式：Unix 时间戳。由于暂时不设置 `refresh_token` 接口，故过期时间可能会很长。

## 链接收藏 Collection

### *[GET] `/collection/list`

获取用户收藏列表。

#### Request

```json
{
  "search": "title or description including sth.",
  "tags": [
    "tag1",
    "tag2",
  ]
}
```

- `search`: 可选，搜索关键词，字符串，用于搜索收藏标题和描述。
- `tags`: 可选，标签，字符串数组，包含多个标签时是 and 关系。

#### Response

```json
{
  "code": 200,
  "msg": null,
  "data": {
    "collections": [
      {
        "collection_id": 1,
        "url": "https://example.com",
        "type": "text",
        "title": "Example",
        "description": "This is an example.",
        "tags": ["example", "test"],
        "created_at": 1683561600,
      },
      {
        "collection_id": 2,
        "url": "https://example2.com/xxx.jpg",
        "type": "image",
        "title": "Example 2",
        "description": "This is an example 2.",
        "tags": ["example", "test"],
        "created_at": 1683561600,
      },
    ]
  },
}
```

- `collection_id`: 收藏 ID，递增整数。
- `type`: 收藏类型，字符串，可能为 `text` 、 `image` 、 `video` 三种类型。
- `url`: 收藏链接，字符串。
- `title`: 收藏标题，字符串。
- `description`: 收藏描述，字符串。
- `tags`: 标签，字符串数组。
- `created_at`: 收藏创建时间，格式：Unix 时间戳。

### *[POST] `/collection/add`

添加收藏链接。

#### Request

```json
{
  "url": "https://example.com",
}
```

#### Response

```json
{
  "code": 200,
  "msg": null,
  "data": {
    "collection_id": 1,
    "url": "https://example.com",
    "type": "text",
    "title": "Example",
    "description": "This is an example.",
    "tags": ["example", "test"],
    "created_at": 1683561600,
  },
}
```

若收藏成功，返回 AI 总结。若失败，在 `data.msg` 中返回错误信息。

### *[GET] `/collection/summary`

获取已收藏的链接的 AI 总结。

#### Request

```json
{
  "search": "title or description including sth.",
  "tags": [
    "tag1",
    "tag2",
  ]
}
```

- `search`: 可选，搜索关键词，字符串，用于搜索收藏标题和描述。
- `tags`: 可选，标签，字符串数组，包含多个标签时是 and 关系。

同用户收藏列表接口。

#### Response

```json
{
  "code": 200,
  "msg": null,
  "data": {
    "summary": "AI summary here."
  },
}
```

### *[GET] `/collection/tag`

获取已收藏的链接的所有 tag。

#### Request

None.

#### Response

```json
{
  "code": 200,
  "msg": null,
  "data": {
    "tags": [
      "tag1",
      "tag2",
    ]
  },
}
```
