-- 限流对象
local key = KEYS[1]
-- 窗口大小
local window = tonumber(ARGV[1])
-- 阈值
local threshold = tonumber( ARGV[2])
local now = tonumber(ARGV[3])
-- 窗口的起始时间
local min = now - window
redis.call('ZREMRANGEBYSCORE', key, '-inf', min)
local cnt = redis.call('ZCOUNT', key, '-inf', '+inf')
if cnt >= threshold then
    -- 执行限流
    return "true"
else
    redis.call('ZADD', key, now, now)
    redis.call('PEXPIRE', key, window)
    return "false"
end