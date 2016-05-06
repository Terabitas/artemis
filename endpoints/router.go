package endpoints

import (
	"net/http"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/nildev/artemis/config"
	"github.com/nildev/artemis/domain"
	"github.com/nildev/artemis/version"
	"github.com/nildev/lib/router"
	"github.com/rs/cors"
)

var (
	ctxLog        *log.Entry
	ASGSupervisor *domain.MultiSupervisor
)

func init() {
	ctxLog = log.WithField("version", version.Version).WithField("git-hash", version.GitHash).WithField("build-time", version.BuiltTimestamp)
}

func Router(cfg config.Config) http.Handler {
	r := mux.NewRouter()

	RegisterRoutes(r, cfg)

	return negroni.New(
		negroni.Wrap(r),
	)
}

func BuildRoutes() []router.Routes {
	rt := make([]router.Routes, 0)

	asgRoutes := router.Routes{
		BasePattern: "/api/v1",
		Routes:      make([]router.Route, 9),
	}

	asgRoutes.Routes[0] = router.Route{
		Name: "github.com/nildev/artemis:Healthz",
		Method: []string{
			"GET",
		},
		Pattern:     "/healthz",
		Protected:   false,
		HandlerFunc: HealthzHandler,
		Queries:     []string{},
	}

	asgRoutes.Routes[1] = router.Route{
		Name: "github.com/nildev/artemis:Setup",
		Method: []string{
			"POST",
		},
		Pattern:     "/asgs",
		Protected:   false,
		HandlerFunc: SetupHandler,
		Queries:     []string{},
	}

	asgRoutes.Routes[2] = router.Route{
		Name: "github.com/nildev/artemis:ReadNodes",
		Method: []string{
			"GET",
		},
		Pattern:     "/nodes",
		Protected:   false,
		HandlerFunc: ReadNodesHandler,
		Queries:     []string{},
	}

	asgRoutes.Routes[3] = router.Route{
		Name: "github.com/nildev/artemis:AddNode",
		Method: []string{
			"POST",
		},
		Pattern:     "/nodes",
		Protected:   false,
		HandlerFunc: AddNodeHandler,
		Queries:     []string{},
	}

	asgRoutes.Routes[4] = router.Route{
		Name: "github.com/nildev/artemis:RemoveNode",
		Method: []string{
			"DELETE",
		},
		Pattern:     "/nodes",
		Protected:   false,
		HandlerFunc: RemoveNodeHandler,
		Queries:     []string{},
	}

	asgRoutes.Routes[5] = router.Route{
		Name: "github.com/nildev/artemis:AddMetrics",
		Method: []string{
			"POST",
		},
		Pattern:     "/metrics",
		Protected:   false,
		HandlerFunc: AddMetricsHandler,
		Queries:     []string{},
	}

	asgRoutes.Routes[6] = router.Route{
		Name: "github.com/nildev/artemis:ReadASG",
		Method: []string{
			"GET",
		},
		Pattern:     "/asg",
		Protected:   false,
		HandlerFunc: ReadASGHandler,
		Queries:     []string{},
	}

	asgRoutes.Routes[7] = router.Route{
		Name: "github.com/nildev/artemis:UpdatePolicies",
		Method: []string{
			"POST",
		},
		Pattern:     "/policies",
		Protected:   false,
		HandlerFunc: UpdatePolicyHandler,
		Queries:     []string{},
	}

	asgRoutes.Routes[8] = router.Route{
		Name: "github.com/nildev/artemis:RemoveASG",
		Method: []string{
			"DELETE",
		},
		Pattern:     "/asgs",
		Protected:   false,
		HandlerFunc: RemoveASGHandler,
		Queries:     []string{},
	}

	rt = append(rt, asgRoutes)

	return rt
}

func RegisterRoutes(r *mux.Router, cfg config.Config) {
	routes := BuildRoutes()
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Secret), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
		Debug:         true,
	})

	corsMlw := cors.New(cors.Options{
		AllowedOrigins:     cfg.CORSAllowedOrigins,
		AllowedMethods:     cfg.CORSAllowedMethods,
		AllowedHeaders:     cfg.CORSAllowedHeaders,
		ExposedHeaders:     cfg.CORSExposedHeaders,
		AllowCredentials:   cfg.CORSAllowCredentials,
		MaxAge:             cfg.CORSMaxAge,
		OptionsPassthrough: cfg.CORSOptionsPassThrough,
		Debug:              cfg.CORSDebug,
	})

	for _, rt := range routes {
		apiRouter := mux.NewRouter().StrictSlash(true)
		for _, route := range rt.Routes {
			handlerFunc := negroni.New()

			// Add JWT middleware if route is protected
			if route.Protected {
				handlerFunc.Use(negroni.HandlerFunc(jwtMiddleware.HandlerWithNext))
			}

			handlerFunc.Use(negroni.Wrap(route.HandlerFunc))

			apiRouter.
				NewRoute().
				Methods(route.Method...).
				Name(route.Name).
				Path(fmt.Sprintf("%s%s", rt.BasePattern, route.Pattern)).
				Queries(route.Queries...).
				Handler(handlerFunc)

			ctxLog.WithField("base-path", rt.BasePattern).
				WithField("protected", route.Protected).
				WithField("name", route.Name).
				WithField("path", fmt.Sprintf("%s%s", rt.BasePattern, route.Pattern)).
				WithField("method", route.Method).
				WithField("query", fmt.Sprintf("%+v", route.Queries)).
				Debugf("%+v", route)
		}

		r.PathPrefix(rt.BasePattern).Handler(negroni.New(corsMlw, negroni.Wrap(apiRouter)))
	}
}
