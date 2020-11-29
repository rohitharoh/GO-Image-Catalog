package system

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"

	"github.com/zenazn/goji/web"
	"net/http"
	"reflect"
)

type Controller struct {
}

type Application struct {
}

func (application *Application) Route(controller interface{}, route string, isPublic bool, roles []string) interface{} {

	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		if !isPublic && c.Env["AuthFailed"].(bool) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			response := make(map[string]interface{})
			response["message"] = UnauthorisedErr.Error()
			errResponse, _ := json.Marshal(response)
			w.Write(errResponse)
		} else {

			var l *log.Entry
			if !isPublic {


				if !Contains(roles, "admin") {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					response := make(map[string]interface{})
					response["message"] = UnauthorisedErr.Error()
					errResponse, _ := json.Marshal(response)
					w.Write(errResponse)
					return
				}


			}

			methodValue := reflect.ValueOf(controller).MethodByName(route)
			methodInterface := methodValue.Interface()

			method := methodInterface.(func(c web.C, w http.ResponseWriter, r *http.Request, l *log.Entry) ([]byte, error))
			result, err := method(c, w, r, l)

			if c.Env["Content-Type"] != nil {
				w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
			} else {
				w.Header().Set("Content-Type", "application/json")
			}

			if err != nil {
				response := make(map[string]interface{})
				if IsFunctionalError(err) {
					response["message"] = err.Error()
					w.WriteHeader(http.StatusPreconditionFailed)
				} else {
					response["message"] = InternalServerError.Error()
					w.WriteHeader(http.StatusInternalServerError)
				}

				errResponse, _ := json.Marshal(response)
				w.Write(errResponse)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write(result)
			}

		}
	}
	return fn
}
