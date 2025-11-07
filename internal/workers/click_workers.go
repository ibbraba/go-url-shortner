package workers

import (
	"log"
	"time"

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository" // Nécessaire pour interagir avec le ClickRepository
)

// StartClickWorkers lance un pool de goroutines "workers" pour traiter les événements de clic.
// Chaque worker lira depuis le même 'clickEventsChan' et utilisera le 'clickRepo' pour la persistance.
func StartClickWorkers(workerCount int, clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
	log.Printf("Starting %d click worker(s)...", workerCount)
	for i := 0; i < workerCount; i++ {
		// Lance chaque worker dans sa propre goroutine.
		// Le channel est passé en lecture seule (<-chan) pour renforcer l'immutabilité du channel à l'intérieur du worker.
		go clickWorker(clickEventsChan, clickRepo)
	}
}

// clickWorker est la fonction exécutée par chaque goroutine worker.
// Elle tourne indéfiniment, lisant les événements de clic dès qu'ils sont disponibles dans le channel.
func clickWorker(clickEventsChan <-chan models.ClickEvent, clickRepo repository.ClickRepository) {
	for event := range clickEventsChan { // Boucle qui lit les événements du channel
		//  Convertir le 'ClickEvent' (reçu du channel) en un modèle 'models.Click'.
		click := &models.Click{
			LinkID:    event.LinkID,
			UserAgent: event.UserAgent,
			IPAddress: event.IpAddress,
			Timestamp: event.Timestamp,
		}

		// Persiste le clic en base de données
		// logique de retry implémentée
		maxRetries := 3
		retryDelay := time.Millisecond * 200
		var err error
		for i := 1; i <= maxRetries; i++ {
			err = clickRepo.CreateClick(click)
			if err == nil {
				log.Printf("Click recorded successfully for LinkID %d", event.LinkID)
				break // ✅ Success
			}

			if i == maxRetries {
				log.Printf("ERROR: Failed to save click after %d attempts: %v", maxRetries, err)
			} else {
				log.Printf("WARN: Failed to save click (attempt %d/%d): %v", i, maxRetries, err)
			}
			time.Sleep(retryDelay)

		}
	}
}
