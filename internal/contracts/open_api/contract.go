package open_api

import (
	"bytes"
	"crypto/tls"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	v1 "auth/internal/contracts/open_api/v1"
	"auth/libs/http_server"
	"auth/libs/listeners"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func New(tls *tls.Certificate, baseContract *v1.BaseContract, authContract *v1.AuthContract) listeners.PortListener {
	cont := &contract{
		baseContract: baseContract,
		authContract: authContract,
	}

	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	handler.Use(AccessMiddleware)
	{
		handler.GET(PathPing, cont.baseContract.Ping)

		auth := handler.Group(PathAuth)
		{
			auth.POST(PathAuthSignIn, cont.authContract.SignIn)
			auth.POST(PathAuthSignUp, cont.authContract.SignUp)
		}
	}
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return http_server.New(handler, tls, "open api")
}

type contract struct {
	baseContract *v1.BaseContract
	authContract *v1.AuthContract
}
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func AccessMiddleware(ginC *gin.Context) {
	ginC.Set("trace_id", uuid.New().String())
	log.
		Info().
		Str("trace_id", ginC.Value("trace_id").(string)).
		Str("method", ginC.Request.Method).
		Str("user agent", ginC.Request.UserAgent()).
		Str("path", ginC.Request.URL.Path).Msg("INCOMING REQUEST")

	w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: ginC.Writer}
	ginC.Writer = w

	ginC.Next()

	log.
		Info().
		Str("trace_id", ginC.Value("trace_id").(string)).
		Str("response_body", w.body.String()).
		Int("response_status", w.Status()).
		Msg("SENDING RESPONSE")

	return
}
