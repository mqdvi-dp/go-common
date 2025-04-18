package rdc

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/zone"
)

func (d *Db) Get(ctx context.Context, key string) (string, error) {
	log := logger.DB(logger.Redis, "mget", key)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:Get")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	// log tracer
	trace.Log("key", key)

	resp, err := d.DB.Get(ctx, key).Result()
	if err != nil {
		trace.SetError(err)
		return "", err
	}

	// log result
	trace.Log("result", resp)

	return resp, nil
}

func (d *Db) GetStruct(ctx context.Context, dest interface{}, key string) error {
	log := logger.DB(logger.Redis, "mget", key)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:GetJSON")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	// make sure, the destination is a pointer
	err := validate(dest)
	if err != nil {
		trace.SetError(err)
		return err
	}

	// log tracer
	trace.Log("key", key)

	// get the value
	result, err := d.DB.Get(ctx, key).Result()
	if err != nil {
		trace.SetError(err)
		return err
	}

	// log result
	trace.Log("result", result)

	err = json.Unmarshal([]byte(result), dest)
	if err != nil {
		trace.SetError(err)
		return err
	}

	return nil
}

func (d *Db) Incr(ctx context.Context, key string) (int64, error) {
	log := logger.DB(logger.Redis, "incr", key)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:Incr")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	// log tracer
	trace.Log("key", key)

	result, err := d.DB.Incr(ctx, key).Result()
	if err != nil {
		trace.SetError(err)
		return -1, err
	}

	// log value
	trace.Log("result", result)

	return result, nil
}

func (d *Db) Decr(ctx context.Context, key string) (int64, error) {
	log := logger.DB(logger.Redis, "decr", key)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:Decr")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	// log tracer
	trace.Log("key", key)

	result, err := d.DB.Decr(ctx, key).Result()
	if err != nil {
		trace.SetError(err)
		return -1, err
	}

	// log value
	trace.Log("result", result)

	return result, nil
}

func (d *Db) Set(ctx context.Context, key string, value interface{}, durations ...time.Duration) error {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:Set")

	var expired time.Duration
	if len(durations) > 0 {
		expired = durations[0]
	} else {
		// default duration is until end of day
		now := time.Now().In(zone.TzJakarta())
		nd := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()) // get the next day
		expired = nd.Sub(now)
	}

	log := logger.DB(logger.Redis, "set", key, value, expired)
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	// log tracer
	trace.Log("key", key)
	trace.Log("value", value)
	trace.Log("expired_duration", expired)

	err := d.DB.Set(ctx, key, value, expired).Err()
	if err != nil {
		trace.SetError(err)
		return err
	}

	return nil
}

func (d *Db) Del(ctx context.Context, keys ...string) error {
	log := logger.DB(logger.Redis, "del", keys)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:Del")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	// log tracer
	trace.Log("keys", keys)

	err := d.DB.Del(ctx, keys...).Err()
	if err != nil {
		trace.SetError(err)
		return err
	}

	return nil
}

func (d *Db) TTL(ctx context.Context, key string) (time.Duration, error) {
	log := logger.DB(logger.Redis, "ttl", key)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:TTL")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	trace.Log("key", key)

	result, err := d.DB.TTL(ctx, key).Result()
	if err != nil {
		trace.SetError(err)
		return -1, err
	}

	return result, nil
}

func (d *Db) Expire(ctx context.Context, key string, duration time.Duration) error {
	log := logger.DB(logger.Redis, "expire", key, duration)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:Expire")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	trace.Log("key", key)
	trace.Log("expired_duration", duration)

	err := d.DB.Expire(ctx, key, duration).Err()
	if err != nil {
		trace.SetError(err)
		return err
	}

	return nil
}

func (d *Db) Keys(ctx context.Context, patternKey string) ([]string, error) {
	log := logger.DB(logger.Redis, "keys", patternKey)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:Keys")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	trace.Log("key", patternKey)

	result, err := d.DB.Keys(ctx, patternKey).Result()
	if err != nil {
		trace.SetError(err)
		return nil, err
	}

	return result, nil
}

func (d *Db) GetKeysAndDelete(ctx context.Context, patternKey string) (err error) {
	log := logger.DB(logger.Redis, "keys and del", patternKey)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:GetKeysAndDelete")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	trace.Log("key", patternKey)
	result, err := d.DB.Keys(ctx, patternKey).Result()
	if err != nil {
		trace.SetError(err)
		return err
	}

	err = d.DB.Del(ctx, result...).Err()
	if err != nil {
		trace.SetError(err)
		return err
	}

	return nil
}

func (d *Db) DoSMembers(ctx context.Context, key string) (members []string, err error) {
	log := logger.DB(logger.Redis, "do smembers", key)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:DoSMembers")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	// log tracer
	trace.Log("key", key)

	members, err = d.DB.SMembers(context.Background(), key).Result()
	if err != nil {
		trace.SetError(err)
		return
	}

	// log result
	trace.Log("result", members)

	return
}

func (d *Db) DoSIsMember(ctx context.Context, key, member string) (exist bool) {
	log := logger.DB(logger.Redis, "do sismember", key, member)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:DoSIsMember")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	// log tracer
	trace.Log("key", key)
	trace.Log("member", member)

	exist, err := d.DB.SIsMember(context.Background(), key, member).Result()
	if err != nil {
		trace.SetError(err)
		return
	}

	// log result
	trace.Log("result", exist)

	return
}

func (d *Db) DoSadd(ctx context.Context, key string, value []string, durations ...time.Duration) (err error) {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:DoSadd")

	var expired time.Duration
	if len(durations) > 0 {
		expired = durations[0]
	} else {
		// default duration is until end of day
		now := time.Now().In(zone.TzJakarta())
		nd := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()) // get the next day
		expired = nd.Sub(now)
	}

	log := logger.DB(logger.Redis, "DoSadd", key, value, expired)
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	// log tracer
	trace.Log("key", key)
	trace.Log("value", value)
	trace.Log("expired_duration", expired)

	for i := 0; i < len(value); i++ {
		err = d.DB.SAdd(ctx, key, value[i]).Err()
		if err != nil {
			trace.SetError(err)
			return
		}
		err = d.DB.Expire(ctx, key, expired).Err()
		if err != nil {
			trace.SetError(err)
			return
		}
	}

	return
}

func (d *Db) HSet(ctx context.Context, key string, field string, value interface{}, durations ...time.Duration) error {
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:HSet")

	var expired time.Duration
	if len(durations) > 0 {
		expired = durations[0]
	} else {
		// default duration is until end of day
		now := time.Now().In(zone.TzJakarta())
		nd := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location()) // get the next day
		expired = nd.Sub(now)
	}

	log := logger.DB(logger.Redis, "hset", key, value)
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	// log tracer
	trace.Log("key", key)
	trace.Log("field", field)
	trace.Log("value", value)
	trace.Log("expired_duration", expired)

	err := d.DB.HSet(ctx, key, field, value).Err()
	if err != nil {
		trace.SetError(err)
		return err
	}

	err = d.DB.Expire(ctx, key, expired).Err()
	if err != nil {
		trace.SetError(err)
		return err
	}

	return nil
}

func (d *Db) HGet(ctx context.Context, key string, field string) (string, error) {
	log := logger.DB(logger.Redis, "HGet", key)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:HGet")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	// log tracer
	trace.Log("key", key)
	trace.Log("field", field)

	resp, err := d.DB.HGet(ctx, key, field).Result()
	if err != nil {
		trace.SetError(err)
		return "", err
	}

	// log result
	trace.Log("result", resp)

	return resp, nil
}

func (d *Db) Close() error {
	return d.DB.Close()
}

func (d *Db) Ping(ctx context.Context) error {
	return d.DB.Ping(ctx).Err()
}

func (d *Db) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	log := logger.DB(logger.Redis, "HGet", key)
	trace, ctx := tracer.StartTraceWithContext(ctx, "Rdc:HGetAll")
	defer func() {
		log.Store(ctx)
		trace.Finish()
	}()

	// log tracer
	trace.Log("key", key)

	resp, err := d.DB.HGetAll(ctx, key).Result()
	if err != nil {
		trace.SetError(err)
		return nil, err
	}

	// log result
	trace.Log("result", resp)

	return resp, nil
}
