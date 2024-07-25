-- date: 2024-7-11
-- version: 1.1
local function rollback(rollbackTab)
	--回滚已经执行的命令 先进后出(后续优化同命令执行)
	local n = #rollbackTab
	for i = n - 1, 1, -1 do
		local v = rollbackTab[i]
		redis.call(unpack(v))
	end
end

local function contrast(cmd, key, args)
	-- set --> del
	if cmd == "SET" then
		return { "DEL", key }
	end

	-- incr | decr --> set | del

	if cmd == "INCR" then
		local res = redis.call("GET", key)
		if res then
			return { "SET", key, res }
		end
		return { "DEL", key }
	end

	if cmd == "DECR" then
		local res = redis.call("GET", key)
		if res then
			return { "SET", key, res }
		end
		return { "DEL", key }
	end

	-- incrBy arg --> set | del
	if cmd == "INCRBY" then
		local res = redis.call("GET", key)
		if res then
			return { "SET", key, res }
		end
		return { "DEL", key }
	end

	-- append --> get + set | del
	if cmd == "APPEND" then
		local oldVal = redis.call("GET", key)
		if oldVal then
			return { "SET", key, oldVal }
		end
		return { "DEL", key }
	end

	-- hmset --> hdel | hmset
	if cmd == "HMSET" and #args == 4 then
		--获取field
		local oldVal = redis.call("HMGET", key, args[3])
		if oldVal[1] then
			return {"HMSET",key,args[3],oldVal[1]}
		end
		return {"HDEL",key,args[3]}
	end
    -- hincrBy
    if cmd == "HINCRBY" and #args == 4 then
    		--获取field
    		local oldVal = redis.call("HMGET", key, args[3])
    		if oldVal[1] then
    			return {"HMSET",key,args[3],oldVal[1]}
    		end
    		return {"HDEL",key,args[3]}
    end
	-- end

	-- lpush | rpush --> lpop | rpop
	if cmd == "LPUSH" then
		return { "LPOP", key }
	end
	if cmd == "RPUSH" then
		return { "RPOP", key }
	end
    -- LPOP和RPOP执行后会返回移除的值

	-- srem --> sadd | nil
	if cmd == "SREM" and #args == 3 then
		local ok = redis.call("SISMEMBER", key, args[3])
		if ok == 1 then
			return { "SADD", key, args[3] }
		end
	end
	-- sadd --> srem | nil
	if cmd == "SADD" and #args == 3 then
		local res = redis.call("SISMEMBER", key, args[3])
		-- set中本就存在，添加无效,不需要回滚
		if res == 0 then
			return { "SREM", key, args[3] }
		end
	end
	-- zrem --> zadd | nil
	if cmd == "ZREM" and #args == 3 then
		local score = redis.call("ZMSCORE", key, args[3])

		if score[1]  then
			return { "ZADD", key, score[1], args[3] }
		end
	end
	-- zadd --> zrem
	if cmd == "ZADD" and #args == 4 then
		-- 检查成员是否存在
		local sc = redis.call("ZSCORE", key, args[4])
		if not sc then
			return { "ZREM", key, args[4] }
		end
	end
	-- zincrBy --> zincrBy
	if cmd == "ZINCRBY" and #args == 4  then
		local score = redis.call("ZMSCORE", key, args[4])
		if score[1] then
			return { "ZADD", key, score[1], args[4] }
		end
		return { "ZREM", key, args[4] }
	end

	-- expire --> expire
	if cmd == "EXPIRE"  and #args == 3 then
		local ttl = redis.call("TTL", key)
		if ttl ~= -2 then
			return {"EXPIRE",key,ttl}
		end
	end
end



local cmds = cjson.decode(ARGV[1])
local rollbackTab = {}
for _, args in ipairs(cmds) do
	if #args < 2 then
		table.insert(rollbackTab, {})
		rollback(rollbackTab)
		-- return { 0, "command error <\"" .. table.concat(args," ") .. "\">", rollbackTab }
		return { "command error <\"" .. table.concat(args," ") .. "\">"}
	end

	-- local cmd = args[1]
	local cmd = string.upper(args[1])
	local key = args[2]

	-- 预先获取回滚命令
	local res = contrast(cmd, key, args)
	if res then
		table.insert(rollbackTab, res)
	end

	-- 执行 Redis 命令，检查执行结果
	local ok = redis.pcall(unpack(args))

	if type(ok) == "table" and ok.err then
		rollback(rollbackTab)
		-- return { 0, "执行命令失败: <\"" .. table.concat(args," ") .. "\">:" .. ok.err, rollbackTab }
		return { "执行命令失败: <\"" .. table.concat(args," ") .. "\">:" .. ok.err}
	end
    if cmd == "LPOP" and ok then
        table.insert(rollbackTab,{"LPUSH",key,ok})
    elseif cmd == "RPOP" and ok then
        table.insert(rollbackTab,{"RPUSH",key,ok})
    end
end
-- return { 1, "OK", rollbackTab }

return { "OK", rollbackTab }