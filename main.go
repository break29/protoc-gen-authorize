package main

import (
	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"

	"github.com/gh1st/protoc-gen-authorize/module"
)

func main() {
	pgs.Init(pgs.DebugEnv("AUTHORIZE_DEBUG")).
		RegisterModule(module.New()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
