package routes

import (
	"github.com/tb/image-catalog/common-packages/system"
	"github.com/tb/image-catalog/imageapp/controllers/api"
	"github.com/zenazn/goji"
)

func PrepareRoutes(application *system.Application) {
	goji.Post("/services/image-catalog/files/list", application.Route(&api.Controller{}, "ListFilesInS3", false, []string{"admin", "user"}))
	goji.Post("/services/image-catalog/files/upload", application.Route(&api.Controller{}, "UploadMultipleFiles", false, []string{"admin", "user"}))
}
