package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type (
	Redis struct {
		opts   *redis.UniversalOptions
		client redis.UniversalClient
	}
)

func New(opts *redis.UniversalOptions) *Redis {
	client := redis.NewUniversalClient(opts)

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(client, redisotel.WithDBStatement(false)); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(client); err != nil {
		panic(err)
	}
	r := &Redis{
		opts:   opts,
		client: client,
	}

	return r
}

func (r *Redis) Name() string {
	return "redis"
}

func (r *Redis) Serve(ctx context.Context) error {
	if r.client != nil {
		return r.client.Ping(ctx).Err()
	}

	return errors.New("redis client not init")
}

func (r *Redis) Shutdown(ctx context.Context) error {
	if r.client != nil {
		return r.client.Close()
	}

	return errors.New("redis client not init")
}

func (r *Redis) Client() redis.UniversalClient {
	return r.client
}

func (r *Redis) NewScript(script string) *Script {
	return newScript(r.client, script)
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.Client().Set(ctx, key, value, expiration)
}

func (r *Redis) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.Client().SetEx(ctx, key, value, expiration)
}

func (r *Redis) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return r.Client().SetNX(ctx, key, value, expiration)
}

func (r *Redis) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return r.Client().SetXX(ctx, key, value, expiration)
}

func (r *Redis) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.Client().Get(ctx, key)
}

func (r *Redis) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	return r.Client().MGet(ctx, keys...)
}

func (r *Redis) MSet(ctx context.Context, values ...interface{}) *redis.StatusCmd {
	return r.Client().MSet(ctx, values...)
}

func (r *Redis) MSetNX(ctx context.Context, values ...interface{}) *redis.BoolCmd {
	return r.Client().MSetNX(ctx, values...)
}

func (r *Redis) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.Client().Del(ctx, keys...)
}

func (r *Redis) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.Client().Exists(ctx, keys...)
}

func (r *Redis) Decr(ctx context.Context, key string) *redis.IntCmd {
	return r.Client().Decr(ctx, key)
}

func (r *Redis) DecrBy(ctx context.Context, key string, decrement int64) *redis.IntCmd {
	return r.Client().DecrBy(ctx, key, decrement)
}

func (r *Redis) Incr(ctx context.Context, key string) *redis.IntCmd {
	return r.Client().Incr(ctx, key)
}

func (r *Redis) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	return r.Client().IncrBy(ctx, key, value)
}

func (r *Redis) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	return r.Client().IncrByFloat(ctx, key, value)
}

func (r *Redis) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.Client().Expire(ctx, key, expiration)
}

func (r *Redis) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return r.Client().ExpireAt(ctx, key, tm)
}

func (r *Redis) PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	return r.Client().PExpire(ctx, key, expiration)
}

func (r *Redis) PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	return r.Client().PExpireAt(ctx, key, tm)
}

func (r *Redis) TTL(ctx context.Context, key string) *redis.DurationCmd {
	return r.Client().TTL(ctx, key)
}

func (r *Redis) PTTL(ctx context.Context, key string) *redis.DurationCmd {
	return r.Client().PTTL(ctx, key)
}

func (r *Redis) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	return r.Client().HDel(ctx, key, fields...)
}

func (r *Redis) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	return r.Client().HExists(ctx, key, field)
}

func (r *Redis) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return r.Client().HGet(ctx, key, field)
}

func (r *Redis) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	return r.Client().HGetAll(ctx, key)
}

func (r *Redis) HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	return r.Client().HIncrBy(ctx, key, field, incr)
}

func (r *Redis) HIncrByFloat(ctx context.Context, key, field string, incr float64) *redis.FloatCmd {
	return r.Client().HIncrByFloat(ctx, key, field, incr)
}

func (r *Redis) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	return r.Client().HKeys(ctx, key)
}

func (r *Redis) HLen(ctx context.Context, key string) *redis.IntCmd {
	return r.Client().HLen(ctx, key)
}

func (r *Redis) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	return r.Client().HMGet(ctx, key, fields...)
}

func (r *Redis) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return r.Client().HSet(ctx, key, values...)
}

func (r *Redis) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	return r.Client().HMSet(ctx, key, values...)
}

func (r *Redis) HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	return r.Client().HSetNX(ctx, key, field, value)
}

func (r *Redis) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	return r.Client().HVals(ctx, key)
}

func (r *Redis) SCard(ctx context.Context, key string) *redis.IntCmd {
	return r.Client().SCard(ctx, key)
}

func (r *Redis) SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return r.Client().SDiff(ctx, keys...)
}

func (r *Redis) SDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return r.Client().SDiffStore(ctx, destination, keys...)
}

func (r *Redis) SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return r.Client().SInter(ctx, keys...)
}

func (r *Redis) SInterStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return r.Client().SInterStore(ctx, destination, keys...)
}

func (r *Redis) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	return r.Client().SIsMember(ctx, key, member)
}

func (r *Redis) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	return r.Client().SMembers(ctx, key)
}

func (r *Redis) SMembersMap(ctx context.Context, key string) *redis.StringStructMapCmd {
	return r.Client().SMembersMap(ctx, key)
}

func (r *Redis) SMove(ctx context.Context, source, destination string, member interface{}) *redis.BoolCmd {
	return r.Client().SMove(ctx, source, destination, member)
}

func (r *Redis) SPop(ctx context.Context, key string) *redis.StringCmd {
	return r.Client().SPop(ctx, key)
}

func (r *Redis) SPopN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return r.Client().SPopN(ctx, key, count)
}

func (r *Redis) SRandMember(ctx context.Context, key string) *redis.StringCmd {
	return r.Client().SRandMember(ctx, key)
}

func (r *Redis) SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	return r.Client().SRandMemberN(ctx, key, count)
}

func (r *Redis) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return r.Client().SRem(ctx, key, members...)
}

func (r *Redis) SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	return r.Client().SUnion(ctx, keys...)
}

func (r *Redis) SUnionStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	return r.Client().SUnionStore(ctx, destination, keys...)
}

func (r *Redis) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return r.Client().ZAdd(ctx, key, members...)
}

func (r *Redis) ZAddNX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return r.Client().ZAddNX(ctx, key, members...)
}

func (r *Redis) ZAddXX(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return r.Client().ZAddXX(ctx, key, members...)
}

func (r *Redis) ZCard(ctx context.Context, key string) *redis.IntCmd {
	return r.Client().ZCard(ctx, key)
}

func (r *Redis) ZCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	return r.Client().ZCount(ctx, key, min, max)
}

func (r *Redis) ZLexCount(ctx context.Context, key, min, max string) *redis.IntCmd {
	return r.Client().ZLexCount(ctx, key, min, max)
}

func (r *Redis) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	return r.Client().ZIncrBy(ctx, key, increment, member)
}

func (r *Redis) ZInterStore(ctx context.Context, destination string, store *redis.ZStore) *redis.IntCmd {
	return r.Client().ZInterStore(ctx, destination, store)
}

func (r *Redis) ZPopMax(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	return r.Client().ZPopMax(ctx, key, count...)
}

func (r *Redis) ZPopMin(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	return r.Client().ZPopMin(ctx, key, count...)
}

func (r *Redis) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return r.Client().ZRange(ctx, key, start, stop)
}

func (r *Redis) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return r.Client().ZRangeWithScores(ctx, key, start, stop)
}

func (r *Redis) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return r.Client().ZRangeByScore(ctx, key, opt)
}

func (r *Redis) ZRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return r.Client().ZRangeByLex(ctx, key, opt)
}

func (r *Redis) ZRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return r.Client().ZRangeByScoreWithScores(ctx, key, opt)
}

func (r *Redis) ZRank(ctx context.Context, key, member string) *redis.IntCmd {
	return r.Client().ZRank(ctx, key, member)
}

func (r *Redis) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	return r.Client().ZRem(ctx, key, members...)
}

func (r *Redis) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd {
	return r.Client().ZRemRangeByRank(ctx, key, start, stop)
}

func (r *Redis) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	return r.Client().ZRemRangeByScore(ctx, key, min, max)
}

func (r *Redis) ZRemRangeByLex(ctx context.Context, key, min, max string) *redis.IntCmd {
	return r.Client().ZRemRangeByLex(ctx, key, min, max)
}

func (r *Redis) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	return r.Client().ZRevRange(ctx, key, start, stop)
}

func (r *Redis) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	return r.Client().ZRevRangeWithScores(ctx, key, start, stop)
}

func (r *Redis) ZRevRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return r.Client().ZRevRangeByScore(ctx, key, opt)
}

func (r *Redis) ZRevRangeByLex(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return r.Client().ZRevRangeByLex(ctx, key, opt)
}

func (r *Redis) ZRevRangeByScoreWithScores(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.ZSliceCmd {
	return r.Client().ZRevRangeByScoreWithScores(ctx, key, opt)
}

func (r *Redis) ZRevRank(ctx context.Context, key, member string) *redis.IntCmd {
	return r.Client().ZRevRank(ctx, key, member)
}

func (r *Redis) ZScore(ctx context.Context, key, member string) *redis.FloatCmd {
	return r.Client().ZScore(ctx, key, member)
}

func (r *Redis) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) *redis.IntCmd {
	return r.Client().ZUnionStore(ctx, dest, store)
}
