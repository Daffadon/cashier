package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tiga-putra-cashier-be/controller"
	"tiga-putra-cashier-be/database"
	"tiga-putra-cashier-be/router"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type Server struct {
	Container   *dig.Container
	ServerReady chan bool
}

func (s *Server) Start() {
	err := s.Container.Invoke(func(
		r *gin.Engine,
		db *gorm.DB,
		pc controller.ProductController,
	) {
		defer database.CloseDB(db)
		if len(os.Args) > 1 {
			Command(db)
		}
		router.AppRouter(r, pc)
		srv := &http.Server{
			Addr:    ":8080",
			Handler: r,
		}
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		if s.ServerReady != nil {
			s.ServerReady <- true
		}
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown: ", err)
		}
		log.Println("Server exiting")
	})
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
}
