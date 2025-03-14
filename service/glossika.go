package service

import (
	"context"
	"fmt"
	"glossika/service/api"
	"glossika/service/controller/accountCtrl"
	"glossika/service/controller/merchandiseCtrl"
	"glossika/service/internal/config"
	"glossika/service/internal/database"
	"glossika/service/internal/flags"
	"glossika/service/internal/model"
	"go.uber.org/dig"
	"gorm.io/gorm"
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
	if err := container.Provide(accountCtrl.New); err != nil {
		panic(err)
	}

	if err := container.Provide(merchandiseCtrl.New); err != nil {
		panic(err)
	}
}

type migratePack struct {
	dig.In

	MySQLGlossika *gorm.DB `name:"glossika"`
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

	if err := container.Invoke(api.NewAccount); err != nil {
		panic(err)
	}

	if err := container.Invoke(api.NewMerchandise); err != nil {
		panic(err)
	}
}

func (srv *glossika) run(server *http.Server) {
	fmt.Printf("Glossika API starts at %s\n", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
