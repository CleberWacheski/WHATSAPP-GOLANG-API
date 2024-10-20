package main

import (
	"log"
	"net/http"
	"sync"
	"whatsapp/application/client"
	"whatsapp/application/controllers"
	"whatsapp/application/queue"
	"whatsapp/application/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Erro ao carregar o arquivo .env:", err)
	}
	utils.EnvironmentInitialize()
	client.WhatsappAPI.Initialize()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}))
		r.Route("/whatsapp", func(r chi.Router) {
			r.Use(middleware.BasicAuth(utils.ENV.BASIC_AUTH_REALM, map[string]string{
				utils.ENV.BASIC_AUTH_USERNAME: utils.ENV.BASIC_AUTH_PASSWORD,
			}))
			r.Post("/session", controllers.CreateSession)
			r.Post("/message", controllers.SendMessage)
			r.Post("/queue/message", controllers.SendQueueMessage)
			r.Post("/document", controllers.SendDocument)
			r.Post("/verify", controllers.VerifySession)
			r.Post("/disconnect", controllers.DisconnectedSession)
		})
		log.Println("Servidor HTTP rodando na porta 3000")
		if err := http.ListenAndServe(":3000", r); err != nil {
			log.Fatalf("Erro ao iniciar o servidor HTTP: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		log.Println("Gerenciador de filas iniciado")
		queue.NewAsyncQueueManagerInitialize()
	}()
	wg.Wait()
}
