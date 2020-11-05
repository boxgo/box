package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	Default = StdConfig("default").Build()
)

func Client() redis.UniversalClient {
	return Default.Client()
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return Client().Set(ctx, key, value, expiration)
}

func SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return Client().SetEX(ctx, key, value, expiration)
}

func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return Client().SetNX(ctx, key, value, expiration)
}

func SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return Client().SetXX(ctx, key, value, expiration)
}

func Get(ctx context.Context, key string) *redis.StringCmd {
	return Client().Get(ctx, key)
}

func MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	return Client().MGet(ctx, keys...)
}

func MSet(ctx context.Context, values ...interface{}) *redis.StatusCmd {
	return Client().MSet(ctx, values...)
}

func MSetNX(ctx context.Context, values ...interface{}) *redis.BoolCmd {
	return Client().MSetNX(ctx, values...)
}

func Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return Client().Del(ctx, keys...)
}

func Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return Client().Exists(ctx, keys...)
}

func Decr(ctx context.Context, key string) *redis.IntCmd {
	return Client().Decr(ctx, key)
}

func DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	return Client().DecrBy(ctx, key, decrement)
}

func Incr(ctx context.Context, key string) *redis.IntCmd {
	return Client().Incr(ctx, key)
}

func IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	return Client().IncrBy(ctx, key, value)
}

func IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	return Client().IncrByFloat(ctx, key, value)
}

func Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return Client().Expire(ctx, key, expiration)
}

func ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return Client().ExpireAt(ctx, key, tm)
}

func PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return Client().PExpire(ctx, key, expiration)
}

func PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return Client().PExpireAt(ctx, key, tm)
}

func TTL(ctx context.Context, key string) *redis.DurationCmd {
	return Client().TTL(ctx, key)
}

func PTTL(ctx context.Context, key string) *redis.DurationCmd {
	return Client().PTTL(ctx, key)
}

func HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	return Client().HDel(ctx, key, fields...)
}

func HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	return Client().HExists(ctx, key, field)
}

func HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return Client().HGet(ctx, key, field)
}

func HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	return Client().HGetAll(ctx, key)
}

func HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	return Client().HIncrBy(ctx, key, field, incr)
}

func HIncrByFloat(ctx context.Context, key, field string, incr float64) *redis.FloatCmd {
	return Client().HIncrByFloat(ctx, key, field, incr)
}

func HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	return Client().HKeys(ctx, key)
}

func HLen(ctx context.Context, key string) *redis.IntCmd {
	return Client().HLen(ctx, key)
}

func HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	return Client().HMGet(ctx, key, fields...)
}

func HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return Client().HSet(ctx, key, values...)
}

func HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	return Client().HMSet(ctx, key, values...)
}

func HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	return Client().HSetNX(ctx, key, field, value)
}

func HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	return Client().HVals(ctx, key)
}

func SCard(ctx context.Context, key string) *redis.IntCmd {
	return Client().SCard(ctx, key)
}

func SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return Client().SDiff(ctx, keys...)
}

func SDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return Client().SDiffStore(ctx, destination, keys...)
}

func SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return Client().SInter(ctx, keys...)
}

func SInterStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return Client().SInterStore(ctx, destination, keys...)
}

func SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	return Client().SIsMember(ctx, key, member)
}

func SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	return Client().SMembers(ctx, key)
}

func SMembersMap(ctx context.Context, key string) *redis.StringStructMapCmd {
	return Client().SMembersMap(ctx, key)
}

func SMove(ctx context.Context, source, destination string, member interface{}) *redis.BoolCmd {
	return Client().SMove(ctx, source, destination, member)
}

func SPop(ctx context.Context, key string) *redis.StringCmd {
	return Client().SPop(ctx, key)
}

func SPopN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return Client().SPopN(ctx, key, count)
}

func SRandMember(ctx context.Context, key string) *redis.StringCmd {
	return Client().SRandMember(ctx, key)
}

func SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return Client().SRandMemberN(ctx, key, count)
}

func SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return Client().SRem(ctx, key, members...)
}

func SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return Client().SUnion(ctx, keys...)
}

func SUnionStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return Client().SUnionStore(ctx, destination, keys...)
}

func ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return Client().ZAdd(ctx, key, members...)
}

func ZAddNX(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return Client().ZAddNX(ctx, key, members...)
}

func ZAddXX(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return Client().ZAddXX(ctx, key, members...)
}

func ZAddCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return Client().ZAddCh(ctx, key, members...)
}

func ZAddNXCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return Client().ZAddNXCh(ctx, key, members...)
}

func ZAddXXCh(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	return Client().ZAddXXCh(ctx, key, members...)
}

func ZIncr(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	return Client().ZIncr(ctx, key, member)
}

func ZIncrNX(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	return Client().ZIncrNX(ctx, key, member)
}

func ZIncrXX(ctx context.Context, key string, member *redis.Z) *redis.FloatCmd {
	return Client().ZIncrXX(ctx, key, member)
}

func ZCard(ctx context.Context, key string) *redis.IntCmd {
	return Client().ZCard(ctx, key)
}

func ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	return Client().ZCount(ctx, key, min, max)
}

func ZLexCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	return Client().ZLexCount(ctx, key, min, max)
}

func ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	return Client().ZIncrBy(ctx, key, increment, member)
}

func ZInterStore(ctx context.Context, destination string, store *redis.ZStore) *redis.IntCmd {
	return Client().ZInterStore(ctx, destination, store)
}

func ZPopMax(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	return Client().ZPopMax(ctx, key, count...)
}

func ZPopMin(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	return Client().ZPopMin(ctx, key, count...)
}

func ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return Client().ZRange(ctx, key, start, stop)
}

func ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return Client().ZRangeWithScores(ctx, key, start, stop)
}

func ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return Client().ZRangeByScore(ctx, key, opt)
}

func ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return Client().ZRangeByLex(ctx, key, opt)
}

func ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return Client().ZRangeByScoreWithScores(ctx, key, opt)
}

func ZRank(ctx context.Context, key, member string) *redis.IntCmd {
	return Client().ZRank(ctx, key, member)
}

func ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return Client().ZRem(ctx, key, members...)
}

func ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd {
	return Client().ZRemRangeByRank(ctx, key, start, stop)
}

func ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	return Client().ZRemRangeByScore(ctx, key, min, max)
}

func ZRemRangeByLex(ctx context.Context, key, min, max string) *redis.IntCmd {
	return Client().ZRemRangeByLex(ctx, key, min, max)
}

func ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return Client().ZRevRange(ctx, key, start, stop)
}

func ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return Client().ZRevRangeWithScores(ctx, key, start, stop)
}

func ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return Client().ZRevRangeByScore(ctx, key, opt)
}

func ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return Client().ZRevRangeByLex(ctx, key, opt)
}

func ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return Client().ZRevRangeByScoreWithScores(ctx, key, opt)
}

func ZRevRank(ctx context.Context, key, member string) *redis.IntCmd {
	return Client().ZRevRank(ctx, key, member)
}

func ZScore(ctx context.Context, key, member string) *redis.FloatCmd {
	return Client().ZScore(ctx, key, member)
}

func ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) *redis.IntCmd {
	return Client().ZUnionStore(ctx, dest, store)
}
