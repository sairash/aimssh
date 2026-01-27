package app

import (
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/gen2brain/beeep"
)

// SendNotification sends a desktop notification.
// For local mode, it uses beeep. For SSH mode, it attempts to use
// the client's notification system (currently supports Linux via notify-send).
func SendNotification(sess ssh.Session, title, body string, isSSH bool) {
	if !isSSH {
		beeep.Alert(title, body, "assets/logo.png")
		return
	}

	// Only attempt SSH notification if we have a valid session
	if sess == nil {
		return
	}

	// Check the client's operating system and send notification accordingly
	switch sess.Context().Value("operating_system") {
	case "linux":
		wish.Command(sess, "notify-send", "-a", "Aimssh", title, body).Run()
	}
}
