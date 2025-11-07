package monitor

import (
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository"
)

type UrlMonitor struct {
	linkRepo    repository.LinkRepository
	interval    time.Duration
	knownStates map[uint]bool
	mu          sync.Mutex
}

func NewUrlMonitor(linkRepo repository.LinkRepository, interval time.Duration) *UrlMonitor {
	return &UrlMonitor{
		linkRepo:    linkRepo,
		interval:    interval,
		knownStates: make(map[uint]bool),
	}
}

func (m *UrlMonitor) Start() {
	log.Printf("[MONITOR] Démarrage du moniteur d'URLs avec un intervalle de %v...", m.interval)
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Exécute une première vérification immédiatement au démarrage
	m.checkUrls()

	// Boucle principale du moniteur, déclenchée par le ticker
	for range ticker.C {
		m.checkUrls()
	}
}

// checkUrls effectue une vérification de l'état de toutes les URLs longues enregistrées.
func (m *UrlMonitor) checkUrls() {
	log.Println("[MONITOR] Lancement de la vérification de l'état des URLs...")

	links, err := m.linkRepo.GetAllLinks()
	if err != nil {
		log.Printf("[MONITOR] ERREUR lors de la récupération des liens pour la surveillance : %v", err)
		return
	}

	for _, link := range links {

		currentState := m.isUrlAccessible(link.LongURL)

		// Protéger l'accès à la map 'knownStates' car 'checkUrls' peut être exécuté concurremment
		m.mu.Lock()
		previousState, exists := m.knownStates[link.ID] // Récupère l'état précédent
		m.knownStates[link.ID] = currentState           // Met à jour l'état actuel
		m.mu.Unlock()

		// Si c'est la première vérification pour ce lien, on initialise l'état sans notifier.
		if !exists {
			log.Printf("[MONITOR] État initial pour le lien %s (%s) : %s",
				link.Shortcode, link.LongURL, formatState(currentState))
			continue
		}

		// Si l'état a changé, générer une fausse notification dans les logs.
		if currentState != previousState {
			log.Printf("[NOTIFICATION] Le lien %s (%s) est passé de %s à %s !",
				link.Shortcode, link.LongURL, formatState(previousState), formatState(currentState))
		}

		if !currentState && previousState {
			log.Printf("[NOTIFICATION] L'URL %s (%s) est maintenant INACCESSIBLE.", link.Shortcode, link.LongURL)
		} else if currentState && !previousState {
			log.Printf("[NOTIFICATION] L'URL %s (%s) est maintenant ACCESSIBLE.", link.Shortcode, link.LongURL)
		}
	}

	log.Println("[MONITOR] Vérification de l'état des URLs terminée.")
}

func (m *UrlMonitor) isUrlAccessible(url string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Head(url)
	if err != nil {
		log.Printf("[MONITOR] Erreur d'accès à l'URL '%s': %v", url, err)
		return false
	}
	defer resp.Body.Close()

	// Déterminer l'accessibilité basée sur le code de statut HTTP.
	return resp.StatusCode >= 200 && resp.StatusCode < 400 // Codes 2xx ou 3xx
}

// formatState est une fonction utilitaire pour rendre l'état plus lisible dans les logs.
func formatState(accessible bool) string {
	if accessible {
		return "ACCESSIBLE"
	}
	return "INACCESSIBLE"
}
