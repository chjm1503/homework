// +build wireinject

package main

import (
	"github.com/chxhjm/app/parse/service/internal/biz"
	"github.com/google/wire"
)

func initApp() {
	panic(wire.Build())
}
