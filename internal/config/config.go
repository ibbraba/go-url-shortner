package config

import (
	"log" // Pour logger les informations ou erreurs de chargement de config

	"github.com/spf13/viper" // La bibliothèque pour la gestion de configuration
)

// Structure principale qui mappe l'intégralité de la configuration de l'application.

type Config struct {
	Server struct {
		Port    int    `mapstructure:"port"`
		BaseURL string `mapstructure:"base_url"`
	} `mapstructure:"server"`

	Database struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"database"`

	Analytics struct {
		BufferSize  int `mapstructure:"buffer_size"`
		WorkerCount int `mapstructure:"worker_count"`
	} `mapstructure:"analytics"`

	Monitor struct {
		IntervalMinutes int `mapstructure:"interval_minutes"`
	} `mapstructure:"monitor"`
}

// LoadConfig charge la configuration de l'application en utilisant Viper.
// Elle recherche un fichier 'config.yaml' dans le dossier 'configs/'.
// Elle définit également des valeurs par défaut si le fichier de config est absent ou incomplet.
func LoadConfig() (*Config, error) {

	viper.AddConfigPath("configs/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Valeurs par défaut si le fichier de config est absent ou incomplet
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.base_url", "http://localhost:8080")
	viper.SetDefault("database.name", "urlshortener.db")
	viper.SetDefault("analytics.buffer_size", 100)
	viper.SetDefault("analytics.worker_count", 5)
	viper.SetDefault("monitor.interval_minutes", 5)

	// Lis le fichier de configuration.
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Avertissement: Impossible de lire le fichier de configuration, utilisation des valeurs par défaut. Détail: %v", err)
	}

	//  Démappe  la configuration lue (ou les valeurs par défaut) dans la structure Config.
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("Erreur: Impossible de démapper la configuration. Détail: %v", err)
		return nil, err
	}

	// Log  pour vérifier la config chargée
	log.Printf("Configuration loaded: Server Port=%d, DB Name=%s, Analytics Buffer=%d, Monitor Interval=%dmin",
		cfg.Server.Port, cfg.Database.Name, cfg.Analytics.BufferSize, cfg.Monitor.IntervalMinutes)

	return &cfg, nil // Retourne la configuration chargée
}
