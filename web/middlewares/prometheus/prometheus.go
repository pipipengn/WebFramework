package prometheus

import (
	"WebFramework/web"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type MiddlewareBuilder struct {
	namespace   string
	subsystem   string
	constLabels map[string]string
	help        string
}

type Options struct {
	Namespace   string
	Subsystem   string
	ConstLabels map[string]string
	Help        string
}

func NewBuilder(opt Options) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		namespace:   opt.Namespace,
		subsystem:   opt.Subsystem,
		constLabels: opt.ConstLabels,
		help:        opt.Help,
	}
}

func (m MiddlewareBuilder) Build() web.Middleware {
	vec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   m.namespace,
		Subsystem:   m.subsystem,
		Help:        m.help,
		ConstLabels: m.constLabels,
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"pattern", "method", "status"})
	prometheus.MustRegister(vec)

	return func(next web.HandleFunc) web.HandleFunc {
		return func(c *web.Context) {
			start := time.Now()
			defer func() {
				go func() {
					duration := time.Now().Sub(start).Milliseconds()
					pattern := c.MatchedRoute
					if pattern == "" {
						pattern = "unknown"
					}
					code := strconv.Itoa(c.RespStatusCode)
					vec.WithLabelValues(pattern, c.Request.Method, code).Observe(float64(duration))
				}()
			}()
			next(c)
		}
	}
}
