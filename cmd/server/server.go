package server

import (
	"busline/internal/app"
	"context"

	"github.com/srutherhub/web-app/controller"
	s "github.com/srutherhub/web-app/server"
)

func StartServer(ctx context.Context, app *app.App) {
	serverConfig := s.ServerConfig{Port: "8080"}
	server := s.New()

	apiController := apiController(app)
	baseController := baseController()

	server.RegisterController(*apiController)
	server.RegisterController(*baseController)

	server.Start(serverConfig)
}

func apiController(app *app.App) *controller.Controller {
	api := &controller.Controller{}
	api.SetBase("/api")

	api.RegisterRoute(controller.Route{Method: "GET", Path: "/{city}/feed", Handler: ApiMiddleWare(app, GetFeedInfoHandler(app))})
	api.RegisterRoute(controller.Route{Method: "GET", Path: "/{city}/routes", Handler: ApiMiddleWare(app, GetActiveRoutes(app))})
	api.RegisterRoute(controller.Route{Method: "GET", Path: "/{city}/stops", Handler: ApiMiddleWare(app, GetNearbyStops(app))})
	return api
}

func baseController() *controller.Controller {
	base := &controller.Controller{}
	base.SetBase("")

	base.RegisterRoute(controller.Route{Method: "GET", Path: "/healthz", Handler: SetJSONContentType(HealthzHandler())})
	return base
}
