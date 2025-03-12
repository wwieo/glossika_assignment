package service

import (
	"context"
	"fmt"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"meme_coin_api/service/api"
	"meme_coin_api/service/controller/memeCoinCtrl"
	"meme_coin_api/service/internal/config"
	"meme_coin_api/service/internal/database"
	"meme_coin_api/service/internal/flags"
	"meme_coin_api/service/internal/model"
	"net/http"
)

func Glossika() Service {
	once.Do(func() {
		srv = &glossika{}
	})

	return srv
}

type glossika struct{}

func (srv *glossika) Run() {
	flags.Parse()

	container := dig.New()

	srv.provideConfig(container)

	srv.provideService(container)

	srv.provideController(container)

	if err := container.Invoke(srv.invokeDBMigrate); err != nil {
		panic(err)
	}

	srv.invokeApiRoutes(container)

	if err := container.Invoke(srv.run); err != nil {
		panic(err)
	}
}

func (srv *glossika) provideConfig(container *dig.Container) {
	if err := container.Provide(config.NewGlossika); err != nil {
		panic(err)
	}
}

func (srv *glossika) provideService(container *dig.Container) {
	if err := container.Provide(func() context.Context {
		return context.TODO()
	}); err != nil {
		panic(err)
	}

	if err := container.Provide(database.NewGlossika); err != nil {
		panic(err)
	}

	if err := container.Provide(api.NewServer); err != nil {
		panic(err)
	}

	if err := container.Provide(api.NewGinEngine); err != nil {
		panic(err)
	}

	if err := container.Provide(api.NewRouterRoot); err != nil {
		panic(err)
	}
}

func (srv *glossika) provideController(container *dig.Container) {
	if err := container.Provide(memeCoinCtrl.New); err != nil {
		panic(err)
	}
}

type migratePack struct {
	dig.In

	MySQLGlossika *gorm.DB `name:"meme_coin"`
}

func (srv *glossika) invokeDBMigrate(pack migratePack) {
	if err := model.GlossikaMigrate(pack.MySQLGlossika); err != nil {
		panic(err)
	}
}

func (srv *glossika) invokeApiRoutes(container *dig.Container) {
	if err := container.Invoke(api.NewBasic); err != nil {
		panic(err)
	}

	if err := container.Invoke(api.NewGlossika); err != nil {
		panic(err)
	}
}

func (srv *glossika) run(server *http.Server) {
	fmt.Printf("Meme Coin API starts at %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
