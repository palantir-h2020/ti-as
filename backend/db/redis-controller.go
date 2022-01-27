package db

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"

	"PalantirIpAnonymizationService/logger"
)

var log = logger.New(false)
var ctx = context.Background()

type RedisDb struct {
	Host     string
	Port     int
	Username string
	Password string
	Database int

	Client *redis.Client
}

func New(host string, port int, username string, password string, db string) (*RedisDb, error) {
	dbIdx, err := strconv.Atoi(db)
	if err != nil {
		log.Warning("Cannot convert database index to intger. Using default database (index = 0).")
		dbIdx = 0
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: password,
		DB:       dbIdx,
	})

	r := RedisDb{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: dbIdx,

		Client: redisClient,
	}

	pong, err := r.Client.Ping(ctx).Result()
	if err != nil {
		log.Error("Cannot connect to Redis database in address " + host + ":" + strconv.Itoa(port) + ".")
		log.PrintError(err)

		return nil, err
	}

	log.Info("PING -> " + pong)

	return &r, nil
}

func (r RedisDb) Insert(key string, value string) (bool, error) {
	err := r.Client.Set(ctx, key, value, 0).Err()

	if err != nil {
		log.PrintError(err)

		return false, err
	}

	return true, nil
}

func (r RedisDb) Get(key string) (string, error) {
	value, err := r.Client.Get(ctx, key).Result()

	if err == redis.Nil {
		// Key does not exists
		log.Error("There is no record with key " + key + ".")
		log.PrintError(err)

		return "", err
	} else if err != nil {
		// An error occured
		log.Error("Cannot retrieve record with key " + key + ".")
		log.PrintError(err)

		return "", err
	} else {
		// Return retrieved value
		return value, nil
	}
}
