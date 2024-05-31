local key = KEYS[1]
local cntKey = ARGV[1]
local delta = tonumber(ARGV[2])

-- 使用HINCRBY命令进行递增/递减，如果键不存在则不执行操作
local result = redis.call("HINCRBY", key, cntKey, delta)

if result then
    return 1
else
    return 0
end
