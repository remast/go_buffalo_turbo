package actions

import (
	"log"
	"net/http"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gorilla/websocket"
	"github.com/r3labs/sse/v2"
	"github.com/unrolled/secure"

	csrf "github.com/gobuffalo/mw-csrf"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr/v2"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_go_buffalo_turbo_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Setup and use translations:
		app.Use(translations())

		app.Use(SkipTurboMiddleware)

		app.GET("/", TaskIndex)
		app.GET("/task/completed", TaskCompleted)
		app.POST("/task/create", TaskCreate)
		app.POST("/task/check", TaskCheck)
		app.GET("/task/new", TaskNew)

		// Setup server sent events
		server := sse.New()
		server.CreateStream("messages")

		go func() {
			for {
				event := RandomFeed("feed-frame-sse")
				server.Publish("messages", &sse.Event{
					Data: []byte(event),
				})
				app.Logger.Printf("Sending Server Sent Event.", event)
				time.Sleep(5 * time.Second)
			}
		}()

		app.GET("/events", buffalo.WrapHandlerFunc(server.HTTPHandler))
		app.GET("/ws", HandeWs)

		app.GET("/feed", FeedIndex)
		app.GET("/feed_slow", FeedIndexSlow)
		app.GET("/feed_sse", FeedIndexWithSSE)
		app.GET("/feed_ws", FeedIndexWithWebsocket)
		app.GET("/feed-frame", FeedFrame)

		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

func SkipTurboMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		skipTurboCookieValue, _ := c.Cookies().Get("skipTurbo")
		if skipTurboCookieValue == "" {
			skipTurboCookieValue = "false"
			skipTurboCookie := http.Cookie{
				Name:    "skipTurbo",
				Value:   skipTurboCookieValue,
				Path:    "/",
				Expires: time.Now().Add(30 * 24 * time.Hour),
			}
			http.SetCookie(c.Response(), &skipTurboCookie)
		}
		c.Set("skipTurbo", skipTurboCookieValue == "true")

		return next(c)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandeWs(cc buffalo.Context) error {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(cc.Response(), cc.Request(), nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	go func() {
		for {
			event := RandomFeed("feed-frame-ws")
			if err := ws.WriteMessage(1, []byte(event)); err != nil {
				log.Println(err)
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()
	return cc.Render(http.StatusOK, r.String("DONE"))

}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.New("app:locales", "../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
