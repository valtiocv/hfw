package hfw

import (
	"errors"
)

type sessRedisStore struct {
	redisIns   *Redis
	prefix     string
	expiration int32
}

var _ sessionStoreInterface = &sessRedisStore{}

var sessRedisStoreIns *sessRedisStore

func NewSessRedisStore() (*sessRedisStore, error) {
	if sessRedisStoreIns == nil {
		redisConfig := Config.Redis
		sessConfig := Config.Session
		if sessConfig.SessID != "" && redisConfig.Server == "" {
			return nil, errors.New("session config error")
		}
		sessRedisStoreIns = &sessRedisStore{
			redisIns:   NewRedis(redisConfig),
			prefix:     "sess_",
			expiration: redisConfig.Expiration,
		}
	}

	return sessRedisStoreIns, nil
}

func (s *sessRedisStore) SetExpiration(expiration int32) {
	s.expiration = expiration
}

func (s *sessRedisStore) IsExist(sessid, key string) (value bool, err error) {
	// key = fmt.Sprintf("%x", md5.Sum([]byte(key)))
	// Debug("IsExist cache key:", sessid, key)

	return s.redisIns.Hexists(s.prefix+sessid, key)
}

func (s *sessRedisStore) Put(sessid, key string, value interface{}) (err error) {
	// key = fmt.Sprintf("%x", md5.Sum([]byte(key)))
	// Debug("Put cache key:", sessid, key, value)

	return s.redisIns.Hset(s.prefix+sessid, key, value)
}

func (s *sessRedisStore) Get(sessid, key string) (value interface{}, err error) {
	// key = fmt.Sprintf("%x", md5.Sum([]byte(key)))
	// Debug("Get cache key:", sessid, key)

	return s.redisIns.Hget(s.prefix+sessid, key)
}

func (s *sessRedisStore) Del(sessid, key string) (err error) {
	// key = fmt.Sprintf("%x", md5.Sum([]byte(key)))
	// Debug("Del cache key:", sessid, key)

	return s.redisIns.Hdel(s.prefix+sessid, key)
}

func (s *sessRedisStore) Destroy(sessid string) (err error) {
	// key = fmt.Sprintf("%x", md5.Sum([]byte(key)))
	// Debug("Del cache key:", sessid)

	_, err = s.redisIns.Del(s.prefix + sessid)

	return
}

func (s *sessRedisStore) Rename(sessid, newid string) (err error) {
	// key = fmt.Sprintf("%x", md5.Sum([]byte(key)))
	// Debug("Rename cache key:", sessid, "to key:", newid)

	err = s.redisIns.Rename(s.prefix+sessid, s.prefix+newid)
	if err != nil {
		return
	}

	if s.expiration > 0 {
		s.redisIns.Expire(s.prefix+newid, s.expiration)
	}

	return
}
