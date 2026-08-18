// Package discordclient defines the boundary between eveping's batch logic
// and the Discord API, so the rest of the codebase can be tested against a
// fake implementation instead of a real discordgo session.
package discordclient

import "time"

// Event status values, matching Discord's GuildScheduledEvent status field.
const (
	EventStatusScheduled = "SCHEDULED"
	EventStatusActive    = "ACTIVE"
	EventStatusCompleted = "COMPLETED"
	EventStatusCanceled  = "CANCELED"
)

// Guild is the internal domain representation of a Discord guild.
type Guild struct {
	ID string
}

// Event is the internal domain representation of a Discord scheduled event.
type Event struct {
	ID                 string
	GuildID            string
	Name               string
	ScheduledStartTime time.Time
	Status             string
}

// User is the internal domain representation of a Discord user.
type User struct {
	ID string
}

// Client is the boundary interface for all Discord API operations eveping
// needs. Implementations: the production discordgo-backed client, and Fake
// for tests.
type Client interface {
	Guilds() []Guild
	ScheduledEvents(guildID string) ([]Event, error)
	EventUsers(guildID, eventID, after string, limit int) ([]User, error)
	SendDM(userID, message string) error
}
