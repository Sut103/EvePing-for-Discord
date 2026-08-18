package discordclient

import (
	"sync"
	"testing"

	"github.com/bwmarrin/discordgo"
)

// TestDiscordgoClient_Guilds_NoDataRace reproduces the concurrency pattern
// between eveping's batch loop (calling Guilds()) and discordgo's own
// gateway read-loop goroutine (calling State.GuildAdd on GUILD_CREATE/etc).
// Run with -race: an unsynchronized read of session.State.Guilds races with
// GuildAdd's locked write to the same slice.
func TestDiscordgoClient_Guilds_NoDataRace(t *testing.T) {
	session := &discordgo.Session{State: discordgo.NewState()}
	client := New(session)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = client.Guilds()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = session.State.GuildAdd(&discordgo.Guild{ID: "guild-1"})
		}
	}()

	wg.Wait()
}
