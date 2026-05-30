package commonmodel

import "github.com/redis/go-redis/v9"

func InitRedisScript() map[LuaOperate]*redis.Script {
	newScript := make(map[LuaOperate]*redis.Script)
	newScript[HSETEX] = redis.NewScript(HSETEXLua)
	newScript[DELBLURRY] = redis.NewScript(DELBLURRYLua)

	return newScript
}

const (
	HSETEX    LuaOperate = "HSETEX" //set多个个带有过期时间、hash类型的key，自带回滚。ARGV的第一个参数传入int类型的时间（ .Second转换,小于等于0则无过期时间），后续每连续两位为field的key和value。多个Key的field之间通过“###KEYS*GAP###”区分
	HSETEXLua string     = `local keyGap = "###KEYS*GAP###"
local expireTime = tonumber(ARGV[1])
local keyPos = 1
local keyLen = #KEYS
local argvPos = 2
local argvLen = #ARGV
if expireTime == nil or expireTime <= 0 then
	expireTime = nil
end
while keyPos <= keyLen do
	while argvPos <= argvLen do
		if ARGV[argvPos] ~= keyGap then
			local hsetResult = redis.pcall("HSET",KEYS[keyPos], ARGV[argvPos], ARGV[argvPos+1])
			if type(hsetResult) == "table" and hsetResult.err ~= nil then
				for k = 1,keyPos do
					redis.pcall("DEL",KEYS[k])
				end
				return redis.error_reply(hsetResult.err)
			end
			argvPos = argvPos + 2
		else
			argvPos = argvPos + 1
			break
		end
	end
	if expireTime and expireTime > 0 then
		local expireResult = redis.pcall("EXPIRE",KEYS[keyPos], expireTime)
		if type(expireResult) == "table" and expireResult.err ~= nil then
			for k = 1,keyPos do
				redis.pcall("DEL",KEYS[k])
			end
			return redis.error_reply(result.err)
		end
	end
	keyPos = keyPos + 1
end
return "OK"`
)
const (
	DELBLURRY    LuaOperate = "DELBlurry" //模糊删除匹配前缀的所有key，返回删除的数量
	DELBLURRYLua string     = `local rawKey=KEYS[1]
local allKey={}
local cursor="0"
repeat
	local result=redis.call("SCAN",cursor,"MATCH",rawKey,"COUNT",100)
	cursor=result[1]
	if #(result[2])>0 then
		for i,key in ipairs(result[2]) do
			table.insert(allKey,key)
		end
	end
until cursor=="0"
if #allKey>0 then
	local result=redis.call("DEL",unpack(allKey))
end
return result`
)
