package discordclient

import "github.com/bwmarrin/discordgo"

func convertStatus(status discordgo.GuildScheduledEventStatus) string {
	switch status {
	case discordgo.GuildScheduledEventStatusScheduled:
		return EventStatusScheduled
	case discordgo.GuildScheduledEventStatusActive:
		return EventStatusActive
	case discordgo.GuildScheduledEventStatusCompleted:
		return EventStatusCompleted
	case discordgo.GuildScheduledEventStatusCanceled:
		return EventStatusCanceled
	default:
		return ""
	}
}

func convertEvent(e *discordgo.GuildScheduledEvent) Event {
	return Event{
		ID:                 e.ID,
		GuildID:            e.GuildID,
		Name:               e.Name,
		ScheduledStartTime: e.ScheduledStartTime,
		Status:             convertStatus(e.Status),
	}
}
