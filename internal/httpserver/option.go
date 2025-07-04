package httpserver

import "github.com/gin-gonic/gin"

type Option interface {
	Key() string
	Apply(engine *gin.Engine)
}

func NewOption(key string, fn func(engine *gin.Engine)) Option {
	return option{key, fn}
}

type option struct {
	key string
	fn  func(engine *gin.Engine)
}

func (o option) Key() string {
	return o.key
}

func (o option) Apply(engine *gin.Engine) {
	o.fn(engine)
}
