package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// // Create a new middleware chain containing the middleware specific to
	// // our dynamic application routes. For now, this chain will only contain
	// // the session middleware but we'll add more to it later.
	// dynamicMiddleware := alice.New(app.session.Enable)

	sessionMiddleware := []func(http.Handler) http.Handler{app.session.Enable, noSurf, app.authenticate}
	dynamicMiddleware := append(sessionMiddleware, app.requireAuthentication)

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", app.home)
	// mux.HandleFunc("/snippet", app.showSnippet)
	// mux.HandleFunc("/snippet/create", app.createSnippet)

	router := chi.NewRouter()

	router.With(sessionMiddleware...).Group(func(router chi.Router) {
		router.Get("/", app.home)
		router.Get("/about", app.about)
	})

	// Without Using Alice: mux.Get("/", app.session.Enable(http.HandlerFunc(app.home)))

	// RESTy routes for "snippet" resource
	router.Route("/snippet", func(router chi.Router) {
		router.With(dynamicMiddleware...).Group(func(router chi.Router) {
			router.Get("/create", app.createSnippetForm) // GET /snippet/create
			router.Post("/create", app.createSnippet)    // POST /snippet/create
		})
		router.With(sessionMiddleware...).Get("/{id:[0-9]+}", app.showSnippet) // GET /snippet/:id
	})

	router.Route("/user", func(router chi.Router) {
		router.With(sessionMiddleware...).Group(func(router chi.Router) {
			router.Get("/signup", app.signupUserForm)
			router.Post("/signup", app.signupUser)
			router.Get("/login", app.loginUserForm)
			router.Post("/login", app.loginUser)
		})

		router.With(dynamicMiddleware...).Group(func(router chi.Router) {
			router.Get("/profile", app.userProfile)
			router.Post("/logout", app.logoutUser)
			router.Get("/change-password", app.changePasswordForm)
			router.Post("/change-password", app.changePassword)
		})
	})

	// Add a new GET /ping route.
	router.Get("/ping", http.HandlerFunc(ping))

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Return the 'standard' middleware chain followed by the servemux.
	return standardMiddleware.Then(router)
}
