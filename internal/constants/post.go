package constants

const (
	PostInternalServerError = 502001                      // 系统错误
	PostEditSuccess         = "Post edit success"         // 帖子编辑成功
	PostEditError           = "Post edit failed"          // 帖子编辑失败
	PostUpdateSuccess       = "Post update success"       // 帖子更新成功
	PostUpdateError         = "Post update failed"        // 帖子更新失败
	PostPublishSuccess      = "Post publish success"      // 帖子发布成功
	PostPublishError        = "Post publish failed"       // 帖子发布失败
	PostWithdrawSuccess     = "Post withdraw success"     // 帖子撤销成功
	PostWithdrawError       = "Post withdraw failed"      // 帖子撤销失败
	PostListPubSuccess      = "Public post query success" // 公开帖子查询成功
	PostListSuccess         = "post query success"        // 公开帖子查询成功
	PostListPubError        = "Public post query failed"  // 公开帖子查询失败
	PostListError           = "post query failed"         // 公开帖子查询失败
	PostDeleteSuccess       = "Post delete success"       // 帖子删除成功
	PostDeleteError         = "Post delete failed"        // 帖子删除失败
	PostGetInteractiveERROR = "get interactive failed"    // 互动信息获取失败
	PostLikedSuccess        = "liked success"
	PostCollectSuccess      = "collect success"
	PostCollectError        = "collect failed"
	PostCanceCollectError   = "cance collect failed"
	PostLikedError          = "liked failed"
	PostCanceLikedSuccess   = "cance liked success"
	PostCanceLikedError     = "cance liked failed"
	PostServerError         = "Post server error"
	PostGetDetailERROR      = "Post get detail failed"   //获取帖子详情失败
	PostGetLikedERROR       = "Post get liked failed"    //获取帖子点赞失败
	PostGetCollectERROR     = "Post get collectd failed" //获取帖子收藏失败
	PostGetPostERROR        = "Post get post failed"     //获取帖子失败
	PostGetDetailSuccess    = "Post get detail success"
	PostGetPubDetailERROR   = "Post get pub detail failed"
	PostGetPubDetailSuccess = "Post get pub detail success"
	PostGetIdsERROR         = "GetByIds failed"
	PostGetIdsSuccess       = "post get ids success"
)
