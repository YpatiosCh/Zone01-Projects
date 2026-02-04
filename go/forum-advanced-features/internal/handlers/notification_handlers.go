package handlers

import (
	"forum/internal/middleware"
	"forum/internal/models"
	"forum/internal/services"
	"html/template"
	"net/http"
)

// NotificationsHandler handles the notifications page
func NotificationsHandler(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the current user
		user := middleware.GetUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get unread notifications count
		unreadCount, err := service.Notify().GetUnreadNotificationCount(user.ID)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get notifications")
			return
		}

		// Get all notifications for the user
		notifications, err := service.Notify().GetUserNotifications(user.ID)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to get notifications")
			return
		}

		data := models.NotificationsData{
			User:          user,
			Notifications: notifications,
			Unread:        unreadCount,
		}

		if err := tmpl.ExecuteTemplate(w, "notifications.html", data); err != nil {
			// Print the actual error to console
			RenderError(tmpl, w, http.StatusInternalServerError, "Template error: "+err.Error())
			return
		}
	}
}

// DeleteNotificationHandler handles deleting a notification
func DeleteNotificationHandler(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Get authenticated user
		user := middleware.GetUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Get notification ID from URL
		notificationID := r.PathValue("notification_id")

		// Delete the notification
		err := service.Notify().DeleteNotification(notificationID, user.ID)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, "Failed to delete notification")
			return
		}

		// Redirect back to notifications page
		http.Redirect(w, r, "/notifications", http.StatusSeeOther)
	}
}

func DeleteAllNotifications(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Get authenticated user
		user := middleware.GetUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		notifications, err := service.Notify().GetUserNotifications(user.ID)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, err.Error())
		}

		for _, notif := range notifications {
			// Delete the notification
			err := service.Notify().DeleteNotification(notif.ID, user.ID)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to delete notification")
				return
			}
		}

		// Redirect back to notifications page
		http.Redirect(w, r, "/notifications", http.StatusSeeOther)
	}
}

func MarkAllAsRead(service services.Services, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			RenderError(tmpl, w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		// Get authenticated user
		user := middleware.GetUser(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		notifications, err := service.Notify().GetUserNotifications(user.ID)
		if err != nil {
			RenderError(tmpl, w, http.StatusInternalServerError, err.Error())
		}

		for _, notif := range notifications {
			// Mark notification as read
			err := service.Notify().MarkNotificationAsRead(notif.ID, user.ID)
			if err != nil {
				RenderError(tmpl, w, http.StatusInternalServerError, "Failed to delete notification")
				return
			}
		}

		// Redirect back to notifications page
		http.Redirect(w, r, "/notifications", http.StatusSeeOther)
	}
}
