package sessioncleanup

import (
	"fmt"
	"forum/internal/repository"
	"time"
)

// CleanupExpiredSessions removes expired sessions from the database.
func CleanupExpiredSessions(repo *repository.Manager) error {

	sessions, err := repo.GetAll("session")
	if err != nil {
		return fmt.Errorf("failed to retrieve sessions: %w", err)
	}
	for _, session := range sessions {
		if expiresAt, ok := session["expires_at"].(time.Time); ok {
			if time.Now().After(expiresAt) {
				err = repo.DeleteRecord("session", session["id"].(string))
				if err != nil {
					return fmt.Errorf("failed to delete expired session %s: %w", session["id"], err)
				}
			}
		} else {
			return fmt.Errorf("session %s does not have a valid expiry date", session["id"])
		}
	}
	return nil
}

// CleanupExpiredSessionsCron is a cron job that runs CleanupExpiredSessions every hour.
func CleanupExpiredSessionsCron(repo *repository.Manager) {
	for {
		err := CleanupExpiredSessions(repo)
		if err != nil {
			fmt.Printf("Error cleaning up expired sessions: %v\n", err)
		} else {
			fmt.Println("Expired sessions cleaned up successfully.")
		}
		time.Sleep(1 * time.Hour) // Run every hour
	}
}
