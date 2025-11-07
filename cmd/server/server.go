package server

import (
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/api/handlers"
	"github.com/axellelanca/urlshortener/internal/monitor"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"

	"github.com/spf13/cobra"
	// Driver SQLite pour GORM
)

// RunServerCmd représente la commande 'run-server' de Cobra.
// C'est le point d'entrée pour lancer le serveur de l'application.
var RunServerCmd = &cobra.Command{
	Use:   "run-server",
	Short: "Lance le serveur API de raccourcissement d'URLs et les processus de fond.",
	Long: `Cette commande initialise la base de données, configure les APIs,
démarre les workers asynchrones pour les clics et le moniteur d'URLs,
puis lance le serveur HTTP.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO : créer une variable qui stock la configuration chargée globalement via cmd.cfg
		// Ne pas oublier la gestion d'erreur et faire un fatalF
		cfg := cmd2.Cfg

		// TODO : Initialiser la connexion à la bBDD
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}

		//  Initialiser les repositories.
		linkRepo := repository.GormLinkRepository(sqlDB)
		clickRepo := repository.GormClickRepository(sqlDB)

		// Laissez le log
		log.Println("Repositories initialisés.")

		//  Initialiser les services métiers.
		linkService := services.NewLinkService(linkRepo)
		clickService := services.NewClickService(clickRepo)

		// Laissez le log
		log.Println("Services métiers initialisés.")

		// TODO : Initialiser le channel ClickEventsChannel (api/handlers) des événements de clic et lancer les workers (StartClickWorkers).
		// Le channel est bufferisé avec la taille configurée.
		// Passez le channel et le clickRepo aux workers.

		clickEventsChannel := make(chan handlers.ClickEvent, cfg.Workers.ClickWorkerBufferSize)

		// TODO : Remplacer les XXX par les bonnes variables
		log.Printf("Channel d'événements de clic initialisé avec un buffer de %d. %d worker(s) de clics démarré(s).",
			cfg.Workers.ClickWorkerBufferSize, cfg.Workers.ClickWorkerCount)

		// TODO : Initialiser et lancer le moniteur d'URLs.
		// Utilisez l'intervalle configuré
		monitorInterval := time.Duration(cfg.Workers.UrlMonitorInterval) * time.Minute
		urlMonitor := monitor.NewUrlMonitor(linkRepo, monitorInterval)

		// TODO Lancez le moniteur dans sa propre goroutine.
		go urlMonitor.Start()
		log.Printf("Moniteur d'URLs démarré avec un intervalle de %v.", monitorInterval)

		// TODO : Configurer le routeur Gin et les handlers API.
		// Passez les services nécessaires aux fonctions de configuration des routes.
		router := handlers.SetupRoutes(linkService, clickService, clickEventsChannel)

		// Pas toucher au log
		log.Println("Routes API configurées.")

		// Créer le serveur HTTP Gin
		serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
		srv := &http.Server{
			Addr:    serverAddr,
			Handler: router,
		}

		// TODO : Démarrer le serveur Gin dans une goroutine anonyme pour ne pas bloquer.
		// Pensez à logger des ptites informations...
		go func() {
			log.Printf("Démarrage du serveur HTTP sur %s...", serverAddr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("FATAL: Échec du démarrage du serveur HTTP: %v", err)
			}
		}()

		// Gére l'arrêt propre du serveur (graceful shutdown).
		// TODO Créez un channel pour les signaux OS (SIGINT, SIGTERM), bufferisé à 1.
		quit :=
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Attendre Ctrl+C ou signal d'arrêt

		// Bloquer jusqu'à ce qu'un signal d'arrêt soit reçu.
		<-quit
		log.Println("Signal d'arrêt reçu. Arrêt du serveur...")

		// Arrêt propre du serveur HTTP avec un timeout.
		log.Println("Arrêt en cours... Donnez un peu de temps aux workers pour finir.")
		time.Sleep(5 * time.Second)

		log.Println("Serveur arrêté proprement.")
	},
}

func init() {
	// TODO : ajouter la commande
	cmd2.RootCmd.AddCommand(RunServerCmd)
}
