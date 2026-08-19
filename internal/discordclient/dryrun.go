package discordclient

import "log"

// DryRunClient wraps a Client so SendDM only logs what would be sent instead
// of contacting Discord. Guilds, ScheduledEvents, and EventUsers are
// delegated straight through to the wrapped Client (via the embedded
// interface), since dry-run mode still needs to discover and log the real
// target events and users.
type DryRunClient struct {
	Client
	logger *log.Logger
}

// NewDryRun wraps client so SendDM logs instead of sending, letting an
// operator run the bot against a real guild to see what a batch would do
// without delivering duplicate reminder DMs.
func NewDryRun(client Client, logger *log.Logger) *DryRunClient {
	return &DryRunClient{Client: client, logger: logger}
}

func (c *DryRunClient) SendDM(userID, message string) error {
	c.logger.Printf("[dry-run] would send DM to user %s: %s", userID, message)
	return nil
}
