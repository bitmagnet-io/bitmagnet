package circuitbreaker

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	Err           = errors.New("circuitbreaker")
	ErrAlreadySet = fmt.Errorf("%w: option already set", Err)
)

type Handler interface {
	http.Handler
	SetOption(option gin.OptionFunc) error
}

func New() Handler {
	return &handler{
		optionSet: make(chan struct{}),
	}
}

type handler struct {
	mtx       sync.Mutex
	optionSet chan struct{}
	option    gin.OptionFunc
	handler   http.Handler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.getHandler().ServeHTTP(w, r)
}

func (h *handler) getHandler() http.Handler {
	<-h.optionSet
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if h.handler == nil {
		h.handler = gin.New(h.option)
	}

	return h.handler
}

func (h *handler) SetOption(option gin.OptionFunc) error {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if h.option != nil {
		return ErrAlreadySet
	}

	h.option = option
	close(h.optionSet)

	return nil
}
