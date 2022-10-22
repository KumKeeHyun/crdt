package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
)

func RangeFromTo(ctx context.Context, cli *redis.Client, list string, from, to int) ([]string, error) {
	return rangeFromTo.
		Eval(ctx, cli, []string{list, list + "-deleted"}, from, to).
		StringSlice()
}

func RangeFrom(ctx context.Context, cli *redis.Client, list string, offset int) ([]string, error) {
	return RangeFromTo(ctx, cli, list, offset, -1)
}

var rangeFromTo = redis.NewScript(`
local list = KEYS[1]
local deleted = KEYS[2]
local offset = ARGV[1]
local to = ARGV[2]

local deleted = redis.call("GET", deleted)
if not deleted then
  deleted = 0
end

offset = offset - deleted
if offset < 0 then
	return redis.error_reply('offset is deleted')
end

return redis.call("LRANGE", list, offset, to)
`)

func Push(ctx context.Context, cli *redis.Client, list string, v ...any) error {
	return cli.RPush(ctx, list, v...).Err()
}