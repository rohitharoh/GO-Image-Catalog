package system

import (
	log "github.com/sirupsen/logrus"
	"github.com/zenazn/goji/web"
	"net/http"
)

func (application *Application) ApplyAuth(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Applied authorization filter!")

		refId := r.URL.Query().Get("referenceId")

		if (refId != "") {
			c.Env["AuthFailed"] = false
			c.Env["refId"] = refId
		} else {
			c.Env["AuthFailed"] = true
			c.Env["refId"] = ""
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
