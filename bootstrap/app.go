package bootstrap

import "gorm.io/gorm"

type Application struct {
	Env *Env
	DB  *gorm.DB
}

func App() *Application {
	env := NewEnv()
	db := NewPostgresDatabase(env)

	return &Application{
		Env: env,
		DB:  db,
	}
}

func (a *Application) CloseDBConnection() {
	ClosePostgresDBConnection(a.DB)
}
