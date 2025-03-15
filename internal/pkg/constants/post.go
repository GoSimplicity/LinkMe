package constants

const (
	// 错误码定义 (4开头表示客户端错误,5开头表示服务端错误)
	PostServerErrorCode     = 502001 // 系统内部错误
	PostCreateErrorCode     = 402001 // 帖子创建错误
	PostEditErrorCode       = 402002 // 帖子编辑错误
	PostUpdateErrorCode     = 402003 // 帖子更新错误
	PostPublishErrorCode    = 402004 // 帖子发布错误
	PostWithdrawErrorCode   = 402005 // 帖子撤回错误
	PostListErrorCode       = 402006 // 帖子列表查询错误
	PostDeleteErrorCode     = 402007 // 帖子删除错误
	PostDetailErrorCode     = 402008 // 帖子详情查询错误
	PostLikeErrorCode       = 402009 // 帖子点赞错误
	PostCollectErrorCode    = 402010 // 帖子收藏错误
	PostValidationErrorCode = 402011 // 帖子参数验证错误
	PostPermissionErrorCode = 402012 // 帖子权限错误
	PostNotFoundErrorCode   = 402013 // 帖子不存在错误
	PostStateErrorCode      = 402014 // 帖子状态错误

	// 成功提示信息
	PostCreateSuccess    = "帖子创建成功"
	PostEditSuccess      = "帖子编辑成功"
	PostUpdateSuccess    = "帖子更新成功"
	PostPublishSuccess   = "帖子发布成功"
	PostWithdrawSuccess  = "帖子撤回成功"
	PostDeleteSuccess    = "帖子删除成功"
	PostListSuccess      = "帖子列表获取成功"
	PostDetailSuccess    = "帖子详情获取成功"
	PostLikeSuccess      = "帖子点赞成功"
	PostUnlikeSuccess    = "取消点赞成功"
	PostCollectSuccess   = "帖子收藏成功"
	PostUncollectSuccess = "取消收藏成功"

	// 错误提示信息
	PostServerError      = "系统内部错误"
	PostCreateError      = "帖子创建失败"
	PostEditError        = "帖子编辑失败"
	PostUpdateError      = "帖子更新失败"
	PostPublishError     = "帖子发布失败"
	PostWithdrawError    = "帖子撤回失败"
	PostDeleteError      = "帖子删除失败"
	PostListError        = "帖子列表获取失败"
	PostDetailEmptyError = "帖子内容或标题为空，请检查后提交"
	PostDetailError      = "帖子详情获取失败"
	PostLikeError        = "帖子点赞失败"
	PostUnlikeError      = "取消点赞失败"
	PostCollectError     = "帖子收藏失败"
	PostUncollectError   = "取消收藏失败"
	PostNotFoundError    = "帖子不存在"
	PostPermissionError  = "没有操作权限"
	PostValidationError  = "参数验证失败"
	PostStateError       = "帖子状态异常"

	// 参数校验相关
	PostInvalidParamsError  = "无效的请求参数"
	PostInvalidIdError      = "无效的帖子ID"
	PostInvalidTitleError   = "无效的帖子标题"
	PostInvalidContentError = "无效的帖子内容"
	PostInvalidPlateError   = "无效的板块ID"
)
