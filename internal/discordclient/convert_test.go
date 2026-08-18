package discordclient

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

func TestConvertStatus(t *testing.T) {
	cases := []struct {
		name string
		in   discordgo.GuildScheduledEventStatus
		want string
	}{
		{"scheduled", discordgo.GuildScheduledEventStatusScheduled, EventStatusScheduled},
		{"active", discordgo.GuildScheduledEventStatusActive, EventStatusActive},
		{"completed", discordgo.GuildScheduledEventStatusCompleted, EventStatusCompleted},
		{"canceled", discordgo.GuildScheduledEventStatusCanceled, EventStatusCanceled},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := convertStatus(tc.in); got != tc.want {
				t.Fatalf("convertStatus(%v) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestConvertEvent_MapsFields(t *testing.T) {
	start := time.Date(2026, 8, 19, 20, 0, 0, 0, time.UTC)
	in := &discordgo.GuildScheduledEvent{
		ID:                 "event-1",
		GuildID:            "guild-1",
		Name:               "Fleet Op",
		ScheduledStartTime: start,
		Status:             discordgo.GuildScheduledEventStatusScheduled,
	}

	got := convertEvent(in)

	want := Event{
		ID:                 "event-1",
		GuildID:            "guild-1",
		Name:               "Fleet Op",
		ScheduledStartTime: start,
		Status:             EventStatusScheduled,
	}
	if got != want {
		t.Fatalf("convertEvent() = %+v, want %+v", got, want)
	}
}
