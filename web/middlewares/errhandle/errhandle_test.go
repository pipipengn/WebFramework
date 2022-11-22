package errhandle

import (
	"WebFramework/web"
	"net/http"
	"testing"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewBuilder().
		AddCode(http.StatusNotFound, []byte(`
<html>
	<body>
		<h1>哈哈哈，走失了</h1>
	</body>
</html>
`)).
		AddCode(http.StatusBadRequest, []byte(`
<html>
	<body>
		<h1>请求不对</h1>
	</body>
</html>
`))

	s := web.NewHttpServer(web.WithMiddleware(builder.Build()))
	_ = s.Start(":8080")
}
