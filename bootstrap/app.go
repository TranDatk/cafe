package bootstrap

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	Env   *Env
	DB    *gorm.DB
	Redis *redis.Client
}

func App() *Application {
	env := NewEnv()
	db := NewPostgresDatabase(env)
	rdb := NewRedisClient(env)

	return &Application{
		Env:   env,
		DB:    db,
		Redis: rdb,
	}
}

func (a *Application) CloseDBConnection() {
	ClosePostgresDBConnection(a.DB)
}
