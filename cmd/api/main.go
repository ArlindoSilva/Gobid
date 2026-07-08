package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/ArlindoSilva/gobid/internal/api"
	"github.com/ArlindoSilva/gobid/internal/services"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// É necessário registrar o tipo uuid.UUID para que ele possa ser armazenado na sessão. O pacote "encoding/gob" é usado para serializar e desserializar dados em Go, e o registro do tipo permite que ele seja corretamente codificado e decodificado quando armazenado na sessão. Sem ele vai dar erro no login, pois o uuid.UUID não é um tipo básico e precisa ser registrado para que o pacote gob saiba como lidar com ele.
	gob.Register(uuid.UUID{}) 

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", 
		os.Getenv("GOBID_DATABASE_USER"),
		os.Getenv("GOBID_DATABASE_PASSWORD"),
		os.Getenv("GOBID_DATABASE_HOST"),
		os.Getenv("GOBID_DATABASE_PORT"),
		os.Getenv("GOBID_DATABASE_NAME"),
	))
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(err)
	}
	// Initialize the session manager with a PostgreSQL store
	s := scs.New()
	s.Store = pgxstore.New(pool)
	s.Lifetime = 24 * time.Hour
	s.Cookie.HttpOnly = true
	s.Cookie.SameSite = http.SameSiteLaxMode //restringe o envio de cookies entre requests de diferentes origens, ajudando a proteger contra ataques CSRF (Cross-Site Request Forgery).

	api := api.Api{
		Router: chi.NewMux(),
		UserService: services.NewUserService(pool),
		Sessions: s,
	}

	api.BindRoutes()

	fmt.Println("Server is running on port :3080")
	if err := http.ListenAndServe("localhost:3080", api.Router); err != nil {
		panic(err)
	}
}