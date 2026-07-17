package api

import (
	//"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	//"github.com/gorilla/csrf"
)

func (api *Api) BindRoutes() {
	api.Router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger, api.Sessions.LoadAndSave)

	// csrfMiddleware := csrf.Protect(
	// 	[]byte(os.Getenv("GOBID_CSRF_KEY")), // Chave de autenticação para proteger contra ataques CSRF.
	// 	csrf.Secure(false), // Define se o cookie CSRF deve ser enviado apenas em conexões HTTPS. Em produção, isso deve ser true para maior segurança.
	// 	csrf.TrustedOrigins([]string{"localhost:3080"}), // Permite requests de origem http://localhost:3080 em desenvolvimento.
	// )
	// api.Router.Use(csrfMiddleware) // Adiciona o middleware CSRF à cadeia de middlewares do roteador, garantindo que todas as rotas estejam protegidas contra ataques CSRF.

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			//r.Get("/csrftoken", api.HandleGetCSRFToken)
			r.Route("/users", func(r chi.Router) {
				r.Post("/signup", api.handleSignupUser)
				r.Post("/login", api.handleLoginUser)
				r.With(api.AuthMiddleware).Post("/logout", api.handleLogoutUser)
			})

			r.Route("/products", func(r chi.Router) {
					r.Post("/", api.handleCreateProduct)
			})
		})
	})
}
// OBS: No Postman é preciso enviar no Headers quando for fazer o POST, PUT ou DELETE, o token CSRF que é retornado na rota /api/v1/csrftoken. O token deve ser enviado no Header com a chave "X-CSRF-Token" e o valor do token retornado. Além de enviar content-type: application/json, Origin: http://localhost:3080, Referer: http://localhost:3080/api/v1/users/login, e incluir o JSON com os dados necessários para a requisição no Body. Isso é necessário para que o servidor possa validar a requisição e garantir que ela não seja maliciosa.