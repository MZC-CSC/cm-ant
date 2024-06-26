package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cloud-barista/cm-ant/pkg/config"
	"github.com/cloud-barista/cm-ant/pkg/database"
	loadHanlder "github.com/cloud-barista/cm-ant/pkg/load/api/handler"
	"github.com/cloud-barista/cm-ant/pkg/load/services"
	priceHanlder "github.com/cloud-barista/cm-ant/pkg/price/api/handler"
	"github.com/cloud-barista/cm-ant/pkg/render"
	"github.com/cloud-barista/cm-ant/pkg/utils"

	_ "github.com/cloud-barista/cm-ant/api"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	zerolog "github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {

	err := config.InitConfig()

	if err != nil {
		log.Fatal("CM-Ant server config error : ", err)
	}

	err = database.InitDatabase()
	if err != nil {
		log.Fatal("CM-Ant database initialize erro :", err)
	}

	router := InitRouter()
	log.Fatal(router.Start(fmt.Sprintf(":%s", config.AppConfig.Server.Port)))
}

// @title CM-ANT API
// @version 0.1
// @description

func InitRouter() *echo.Echo {
	e := echo.New()

	e.Static("/web/templates", utils.RootPath()+"/web/templates")
	e.Static("/css", utils.RootPath()+"/web/css")
	e.Static("/js", utils.RootPath()+"/web/js")

	e.Use(
		middleware.RequestLoggerWithConfig(
			middleware.RequestLoggerConfig{
				LogError:         true,
				LogRequestID:     true,
				LogRemoteIP:      true,
				LogHost:          true,
				LogMethod:        true,
				LogURI:           true,
				LogUserAgent:     false,
				LogStatus:        true,
				LogLatency:       true,
				LogContentLength: true,
				LogResponseSize:  true,
				LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
					if v.Error == nil {
						zerolog.Info().
							Str("id", v.RequestID).
							Str("client_ip", v.RemoteIP).
							Str("host", v.Host).
							Str("method", v.Method).
							Str("URI", v.URI).
							Int("status", v.Status).
							Int64("latency", v.Latency.Nanoseconds()).
							Str("latency_human", v.Latency.String()).
							Str("bytes_in", v.ContentLength).
							Int64("bytes_out", v.ResponseSize).
							Msg("request")
					} else {
						zerolog.Error().
							Err(v.Error).
							Str("id", v.RequestID).
							Str("client_ip", v.RemoteIP).
							Str("host", v.Host).
							Str("method", v.Method).
							Str("URI", v.URI).
							Int("status", v.Status).
							Int64("latency", v.Latency.Nanoseconds()).
							Str("latency_human", v.Latency.String()).
							Str("bytes_in", v.ContentLength).
							Int64("bytes_out", v.ResponseSize).
							Msg("request error")
					}
					return nil
				},
			},
		),
		middleware.Recover(),
		middleware.RequestID(),
		middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)),
		middleware.CORS(),
	)

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "request timeout",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Println(c.Path())
		},
		Timeout: 300 * time.Second,
	}))

	// config template
	tmpl := render.NewTemplate()

	e.Renderer = tmpl

	antRouter := e.Group("/ant")

	{

		antRouter.GET("/swagger/*", echoSwagger.WrapHandler)
		antRouter.GET("", func(c echo.Context) error {
			return c.Render(http.StatusOK, "home.page.tmpl", nil)
		})

		antRouter.GET("/results", func(c echo.Context) error {
			result, err := services.GetAllLoadExecutionConfig()

			if err != nil {
				log.Printf("error while get load test execution config; %+v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, map[string]any{
					"message": "something went wrong.try again.",
				})

			}

			return c.Render(http.StatusOK, "results.page.tmpl", result)
		})

		apiRouter := antRouter.Group("/api")

		versionRouter := apiRouter.Group("/v1")

		versionRouter.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"message": "CM-Ant API server is running",
			})
		})

		connectionRouter := versionRouter.Group("/env")

		{
			connectionRouter.GET("", loadHanlder.GetAllLoadEnvironments())
		}

		loadRouter := versionRouter.Group("/load")

		{
			// load tester
			loadRouter.POST("/tester", loadHanlder.InstallLoadTesterHandler())
			loadRouter.DELETE("/tester/:envId", loadHanlder.UninstallLoadTesterHandler())

			// load test execution
			loadRouter.POST("/start", loadHanlder.RunLoadTestHandler())
			loadRouter.POST("/stop", loadHanlder.StopLoadTestHandler())

			// load test result
			loadRouter.GET("/result", loadHanlder.GetLoadTestResultHandler())
			loadRouter.GET("/result/metrics", loadHanlder.GetLoadTestMetricsHandler())

			// load test history
			loadRouter.GET("/config", loadHanlder.GetAllLoadConfigHandler())
			loadRouter.GET("/config/:loadTestKey", loadHanlder.GetLoadConfigHandler())

			// load test state
			loadRouter.GET("/state", loadHanlder.GetAllLoadExecutionStateHandler())
			loadRouter.GET("/state/:loadTestKey", loadHanlder.GetLoadExecutionStateHandler())

			// load test metrics agent
			loadRouter.POST("/agent", loadHanlder.InstallAgent())
			loadRouter.GET("/agent", loadHanlder.GetAllAgentInstallInfo())
			loadRouter.DELETE("/agent/:agentInstallInfoId", loadHanlder.UninstallAgent())
		}

		priceRouter := versionRouter.Group("/price")

		{
			priceRouter.POST("", priceHanlder.GetPriceInfoHandler())
		}

	}
	return e
}
