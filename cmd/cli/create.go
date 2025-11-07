package cli

import (
	"fmt"
	"log"
	"net/url"
	"os"

	// Pour valider le format de l'URL

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// Driver SQLite pour GORM
)

var longURLFlag string

// TODO : Faire une variable longURLFlag qui stockera la valeur du flag --url

// CreateCmd représente la commande 'create'
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO 1: Valider que le flag --url a été fourni.
		if longURLFlag == "" {
			fmt.Println("Erreur: le flag --url est requis.")
			os.Exit(1)
		}

		_, err := url.ParseRequestURI(longURLFlag)
		if err != nil {
			fmt.Printf("Erreur: l'URL fournie n'est pas valide: %v\n", err)
			os.Exit(1)
		}

		// TODO : Charger la configuration chargée globalement via cmd.cfg
		cfg := cmd2.Cfg

		// TODO : Initialiser la connexion à la base de données SQLite.

		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: impossible de se connecter à la base SQLite: %v", err)
		}

		// Récupère la connexion SQL sous-jacente
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: impossible d'obtenir la DB SQL sous-jacente: %v", err)
		}

		// Ferme la connexion à la fin
		defer sqlDB.Close()

		// TODO : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService

		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)

		// TODO : Appeler le LinkService et la fonction CreateLink pour créer le lien court.
		link, err := linkService.CreateLink(longURLFlag)
		if err != nil {
			log.Fatalf("FATAL: Échec de la création du lien court: %v", err)
			os.Exit(1)
		}

		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)
		fmt.Printf("URL courte créée avec succès:\n")
		fmt.Printf("Code: %s\n", link.ShortCode)
		fmt.Printf("URL complète: %s\n", fullShortURL)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	cmd2.RootCmd.AddCommand(CreateCmd)
	CreateCmd.Flags().StringVar(&longURLFlag, "url", "", "URL longue à raccourcir")
	CreateCmd.MarkFlagRequired("url")
}
