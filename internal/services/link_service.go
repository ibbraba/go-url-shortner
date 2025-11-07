package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm" // Nécessaire pour la gestion spécifique de gorm.ErrRecordNotFound

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository" // Importe le package repository
)

// Définition du jeu de caractères pour la génération des codes courts.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type LinkService struct {
	linkRepo repository.LinkRepository
}

// LinkService est une structure qui g fournit des méthodes pour la logique métier des liens.
// Elle détient linkRepo qui est une référence vers une interface LinkRepository.
// IMPORTANT : Le champ doit être du type de l'interface (non-pointeur).

// NewLinkService crée et retourne une nouvelle instance de LinkService.
func NewLinkService(linkRepo repository.LinkRepository) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
	}
}

// GenerateShortCode est une méthode rattachée à LinkService
// Elle génère un code court aléatoire d'une longueur spécifiée. Elle prend une longueur en paramètre et retourne une string et une erreur
// Il utilise le package 'crypto/rand' pour éviter la prévisibilité.
// Je vous laisse chercher un peu :) C'est faisable en une petite dizaine de ligne
func (s *LinkService) GenerateShortCode(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	shortCode := make([]byte, length)
	for i := range b {
		shortCode[i] = charset[int(b[i])%len(charset)]
	}
	return string(shortCode), nil
}

// CreateLink crée un nouveau lien raccourci.
// Il génère un code court unique, puis persiste le lien dans la base de données.
func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {
	const maxRetries = 5
	var shortCode string

	// Essayez de générer un code, vérifiez s'il existe déjà en base, et retentez si une collision est trouvée.
	// Limitez le nombre de tentatives pour éviter une boucle infinie.

	for i := 0; i < maxRetries; i++ {
		// Génère un code de 6 caractères (GenerateShortCode)

		code, err := s.GenerateShortCode(6)

		if err != nil {
			return nil, fmt.Errorf("failed to generate short code: %w", err)
		}

		// Vérifie si le code généré existe déjà en base de données (GetLinkbyShortCode)
		_, err = s.linkRepo.GetLinkByShortCode(code)
		// On ignore la première valeur

		if err != nil {
			// Si l'erreur est 'record not found' de GORM, cela signifie que le code est unique.
			if errors.Is(err, gorm.ErrRecordNotFound) {

				shortCode = code // Le code est unique, on peut l'utiliser
				break            // Sort de la boucle de retry
			}
			// Si c'est une autre erreur de base de données, retourne l'erreur.
			return nil, fmt.Errorf("database error checking short code uniqueness: %w", err)
		}

		// Si aucune erreur (le code a été trouvé), cela signifie une collision.
		log.Printf("Short code '%s' already exists, retrying generation (%d/%d)...", code, i+1, maxRetries)
		// La boucle continuera pour générer un nouveau code.
	}

	// Si après toutes les tentatives, aucun code unique n'a été trouvé... Errors.New

	if shortCode == "" {
		return nil, errors.New("failed to generate a unique short code after multiple attempts")
	}

	// Crée une nouvelle instance du modèle Link.

	link := &models.Link{
		LongURL:   longURL,
		Shortcode: shortCode,
		CreatedAt: time.Now(),
	}

	// Persiste le nouveau lien dans la base de données via le repository (CreateLink)

	if err := s.linkRepo.CreateLink(link); err != nil {
		return nil, fmt.Errorf("failed to save link: %w", err)
	}

	// Retourne le lien créé

	return link, nil
}

// GetLinkByShortCode récupère un lien via son code court.
// Il délègue l'opération de recherche au repository.
func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	// Retourner le lien trouvé ou une erreur si non trouvé/problème DB.

	link, err := s.linkRepo.GetLinkByShortCode(shortCode)

	return link, err
}

// GetLinkStats récupère les statistiques pour un lien donné (nombre total de clics).
// Il interagit avec le LinkRepository pour obtenir le lien, puis avec le ClickRepository
func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	// Récupérer le lien par son shortCode
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)

	// Compter le nombre de clics pour ce LinkID
	nbr, err := s.linkRepo.CountClicksByLinkID(link.ID)

	// on retourne les 3 valeurs
	return link, nbr, err
}
