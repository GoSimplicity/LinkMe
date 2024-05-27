# LinkMe API接口清单

## 用户模块

### 注册接口
- **路径**: `/users/signup`
- **方法**: POST
- **描述**: 用户注册，使用邮箱格式用户名，强密码格式。
- **请求参数**:
  - `email`
  - `password`
  - `confirmPassword`

### 登录接口
- **路径**: `/users/login`
- **方法**: POST
- **描述**: 用户登录，使用邮箱和密码。
- **请求参数**:
  - `email`
  - `password`

### 测试接口
- **路径**: `/users/hello`
- **方法**: GET
- **描述**: 用于测试JWT校验。

### 刷新令牌接口
- **路径**: `/users/refresh_token`
- **方法**: PUT
- **描述**: 刷新X-Jwt-Token，需传入X-Refresh-Token。

### 登出接口
- **路径**: `/users/logout`
- **方法**: POST
- **描述**: 退出系统登陆，需传入X-Refresh-Token（请求头）、X-Jwt-Token（Authorization）。

## 帖子模块

### 编辑帖子接口
- **路径**: `/posts/edit`
- **方法**: POST
- **描述**: 编辑帖子，需传入帖子标题和内容。
- **请求参数**:
  - `title`
  - `content`

### 发布帖子接口
- **路径**: `/posts/publish`
- **方法**: PUT
- **描述**: 发布帖子，需传入帖子ID。
- **请求参数**:
  - `postId`

### 撤销帖子接口
- **路径**: `/posts/withdraw`
- **方法**: PUT
- **描述**: 撤销帖子，需传入帖子ID。
- **请求参数**:
  - `postId`

### 更新帖子接口
- **路径**: `/posts/update`
- **方法**: PUT
- **描述**: 更新帖子，需传入帖子ID、标题和内容。
- **请求参数**:
  - `postId`
  - `title`
  - `content`

### 获取公开帖子列表接口
- **路径**: `/posts/list_pub`
- **方法**: GET
- **描述**: 获取公开帖子列表，可传入分页参数。
- **请求参数**:
  - `page`
  - `size`

### 删除帖子接口
- **路径**: `/posts/1`
- **方法**: DELETE
- **描述**: 删除指定ID的帖子。

### 获取作者帖子详情接口
- **路径**: `/posts/detail/1`
- **方法**: GET
- **描述**: 获取指定ID的作者帖子详情。

### 获取公开帖子详情接口
- **路径**: `/posts/detail_pub/2`
- **方法**: GET
- **描述**: 获取指定ID的公开帖子详情。

### 获取作者帖子列表接口
- **路径**: `/posts/list`
- **方法**: GET
- **描述**: 获取作者帖子列表，可传入分页参数。
- **请求参数**:
  - `page`
  - `size`
