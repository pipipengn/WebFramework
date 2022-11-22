package web

import "sync"

type Middleware func(next HandleFunc) HandleFunc

// Intercept =====================================
type Intercept interface {
	Before(c *Context)
	After(c *Context)
	Surround(c *Context)
}

// Chain =========================================
type Chain []HandleFuncV1

type HandleFuncV1 func(c *Context) (next bool)

type ChainV1 struct {
	handlers []HandleFuncV1
}

func (c ChainV1) Run(ctx *Context) {
	for _, handler := range c.handlers {
		if next := handler(ctx); !next {
			return
		}
	}
}

// Net =========================================
type Net struct {
	handlers []HandleFuncV2
}

func (n Net) Run(ctx *Context) {
	wg := sync.WaitGroup{}
	for _, handler := range n.handlers {
		h := handler
		if h.concurrent {
			wg.Add(1)
			go func() {
				h.Run(ctx)
				wg.Done()
			}()
		} else {
			h.Run(ctx)
		}
	}
	wg.Wait()
}

type HandleFuncV2 struct {
	concurrent bool
	handlers   []*HandleFuncV2
}

func (h HandleFuncV2) Run(ctx *Context) {

}
