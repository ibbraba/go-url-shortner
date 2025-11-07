package cmd

import (
	"log"
	"os"

	"github.com/axellelanca/urlshortener/internal/config"
	"github.com/spf13/cobra"
)

var Cfg *config.Config

var RootCmd = &cobra.Command{
	Use:   "url-shortener",
	Short: "Un service de raccourcissement d'URLs avec API REST et CLI",
	Long:  "'url-shortener' est une application complète pour gérer des URLs courtes. Elle inclut un serveur API pour le raccourcissement et la redirection, ainsi qu'une interface en ligne de commande pour l'administration. Utilisez 'url-shortener [command] --help' pour plus d'informations sur une commande.",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("Erreur lors de l'exécution de la commande principale : %v", err)
	}

	os.Exit(1)
}

func init() {

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	var err error
	Cfg, err = config.LoadConfig()
	if err != nil {

		log.Printf("Attention: Problème lors du chargement de la configuration: %v. Utilisation des valeurs par défaut.", err)
	}
	log.Println("Configuration chargée avec succès.")
}
