package app

// Config holds application-wide configuration
type Config struct {
	// RunAsSSH determines if the app runs as an SSH server or local TUI
	RunAsSSH bool

	// Server configuration (only used when RunAsSSH is true)
	Server ServerConfig
}

// DefaultConfig returns the default application configuration
func DefaultConfig() Config {
	return Config{
		RunAsSSH: false,
		Server:   DefaultServerConfig(),
	}
}
