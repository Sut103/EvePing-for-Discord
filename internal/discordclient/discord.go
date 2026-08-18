package discordclient

import "github.com/bwmarrin/discordgo"

// discordgoClient is the production Client implementation, backed by a real
// discordgo session.
type discordgoClient struct {
	session *discordgo.Session
}

// New wraps an established discordgo session as a Client.
func New(session *discordgo.Session) Client {
	return &discordgoClient{session: session}
}

func (c *discordgoClient) Guilds() []Guild {
	guilds := make([]Guild, 0, len(c.session.State.Guilds))
	for _, g := range c.session.State.Guilds {
		guilds = append(guilds, Guild{ID: g.ID})
	}
	return guilds
}

func (c *discordgoClient) ScheduledEvents(guildID string) ([]Event, error) {
	events, err := c.session.GuildScheduledEvents(guildID, false)
	if err != nil {
		return nil, err
	}
	result := make([]Event, 0, len(events))
	for _, e := range events {
		result = append(result, convertEvent(e))
	}
	return result, nil
}

func (c *discordgoClient) EventUsers(guildID, eventID, after string, limit int) ([]User, error) {
	users, err := c.session.GuildScheduledEventUsers(guildID, eventID, limit, false, "", after)
	if err != nil {
		return nil, err
	}
	result := make([]User, 0, len(users))
	for _, u := range users {
		result = append(result, User{ID: u.User.ID})
	}
	return result, nil
}

func (c *discordgoClient) SendDM(userID, message string) error {
	channel, err := c.session.UserChannelCreate(userID)
	if err != nil {
		return err
	}
	_, err = c.session.ChannelMessageSend(channel.ID, message)
	return err
}
