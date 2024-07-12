package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

const (
	SCRIPTPATH = "./script/rollback.lua"
)

var (
	hash   string
	script string
)

type RedisEXCmd interface {
	Commit(ctx context.Context, args ...string) error //提交具有修改性的redis命令，失败自动回滚
}

type redisEXCmd struct {
	cmd redis.Cmdable
}

func init() {
	//加载脚本
	if script == "" {
		t, _ := os.ReadFile(SCRIPTPATH)
		script = string(t)
	}
}
func NewRedisCmd(cmd redis.Cmdable) RedisEXCmd {

	return &redisEXCmd{
		cmd: cmd,
	}
}

// Commit 执行脚本
func (rc *redisEXCmd) Commit(ctx context.Context, args ...string) error {
	//构造redis命令
	cmd, _ := json.Marshal([]interface{}{args})

	//先检测脚本是否缓存
	if exist, _ := rc.cmd.ScriptExists(ctx, hash).Result(); exist[0] {
		//执行缓存脚本
		res, _ := rc.cmd.EvalSha(ctx, hash, []string{}, cmd).Result()
		fmt.Println(res)
		return nil
	}
	//缓存脚本
	hash = rc.cmd.ScriptLoad(ctx, script).Val()
	res, _ := rc.cmd.EvalSha(ctx, hash, []string{}, cmd).Result()
	fmt.Println(res)
	return nil
}
