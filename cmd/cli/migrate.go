package cli

import (
	"fmt"
	"log"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/glebarez/sqlite"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// MigrateCmd représente la commande 'migrate'
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Exécute les migrations de la base de données pour créer ou mettre à jour les tables.",
	Long: `Cette commande se connecte à la base de données configurée (SQLite)
			et exécute les migrations automatiques de GORM pour créer les tables 'links' et 'clicks'
			basées sur les modèles Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := cmd2.Cfg

		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: impossible de se connecter à la base SQLite: %v", err)
		}

		// Récupère la connexion SQL sous-jacente
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: impossible d'obtenir la DB SQL sous-jacente: %v", err)
		}

		defer sqlDB.Close()

		// Exécute les migrations automatiques de GORM.
		err = db.AutoMigrate(&models.Link{}, &models.Click{})
		if err != nil {
			log.Fatalf("FATAL: impossible d'exécuter les migrations: %v", err)
		}

		// Pas touche au log
		fmt.Println("Migrations de la base de données exécutées avec succès.")
	},
}

func init() {
	cmd2.RootCmd.AddCommand(MigrateCmd)
}
