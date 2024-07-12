package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"os"
)

const (
	SCRIPTPATH = "./script/commit.lua"
)

var (
	hash   string
	script string
)

type RedisEXCmd interface {
	AddCommand(args ...string) //添加redis命令
	Exec() error               //执行具有修改性的redis命令，失败自动回滚
	Rollback()                 //可调用完成回滚
}

type redisEXCmd struct {
	cmd         redis.Cmdable
	commandsTab []interface{} //存储命令
	rollbackTab []interface{} //存储回滚命令
	commandsCnt int           //记录命令数量
	rollbackCnt int           //记录回滚命令数量
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

func (rc *redisEXCmd) AddCommand(args ...string) {
	var v []interface{}
	for _, arg := range args {
		v = append(v, arg)
	}
	rc.commandsTab = append(rc.commandsTab, v)
	rc.commandsCnt++
}

// Exec 执行脚本
func (rc *redisEXCmd) Exec() error {
	//fmt.Println(rc.commandsTab[:rc.commandsCnt])
	cmd, _ := json.Marshal(rc.commandsTab[:rc.commandsCnt])
	//标记清除存储的命令
	rc.commandsCnt = 0

	ctx := context.Background()
	//先检测脚本是否缓存
	if exist, _ := rc.cmd.ScriptExists(ctx, hash).Result(); !exist[0] {
		//缓存脚本
		hash = rc.cmd.ScriptLoad(ctx, script).Val()
	}
	//执行缓存脚本
	res, err := rc.cmd.EvalSha(ctx, hash, []string{}, cmd).Result()
	if err != nil {
		//日志

		return err
	}
	//redis命令执行失败
	if res.([]interface{})[0] != "OK" {
		return errors.New(res.([]interface{})[0].(string))
	}
	//解析可能执行的回滚命令
	cmds := res.([]interface{})[1].([]interface{})
	rc.rollbackCnt = len(cmds)
	for _, v := range cmds {
		rc.rollbackTab = append(rc.rollbackTab, v.([]interface{}))
	}
	//todo 超时取消

	return nil
}

func (rc *redisEXCmd) Rollback() {
	//fmt.Println(rc.rollbackTab)
	luaScript := `
	local cmds = cjson.decode(ARGV[1])	
	for _, cmd in ipairs(cmds) do 
		redis.call(unpack(cmd))
	end
`

	cmd, _ := json.Marshal(rc.rollbackTab[:rc.rollbackCnt])
	rc.rollbackCnt = 0

	rc.cmd.Eval(context.Background(), luaScript, []string{}, cmd)

}
