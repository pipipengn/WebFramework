package prometheus

import (
	"WebFramework/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestPrometheus(t *testing.T) {
	builder := NewBuilder(Options{
		Namespace:   "namespace_test",
		Subsystem:   "subsystem_test",
		ConstLabels: nil,
		Help:        "help_test",
	})
	s := web.NewHttpServer(web.WithMiddleware(builder.Build()))

	s.Get("/user", func(c *web.Context) {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)+1))
		c.JSON(200, map[string]string{"name": "ppp"})
	})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":8081", nil)
	}()

	_ = s.Start(":8080")
}
