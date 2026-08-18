package discordclient

// Fake is an in-memory Client implementation for tests. Each *Func field is
// optional; a nil field returns a zero value with no error so callers only
// need to configure the behavior their test cares about.
type Fake struct {
	GuildsFunc          func() []Guild
	ScheduledEventsFunc func(guildID string) ([]Event, error)
	EventUsersFunc      func(guildID, eventID, after string, limit int) ([]User, error)
	SendDMFunc          func(userID, message string) error
}

func (f *Fake) Guilds() []Guild {
	if f.GuildsFunc == nil {
		return nil
	}
	return f.GuildsFunc()
}

func (f *Fake) ScheduledEvents(guildID string) ([]Event, error) {
	if f.ScheduledEventsFunc == nil {
		return nil, nil
	}
	return f.ScheduledEventsFunc(guildID)
}

func (f *Fake) EventUsers(guildID, eventID, after string, limit int) ([]User, error) {
	if f.EventUsersFunc == nil {
		return nil, nil
	}
	return f.EventUsersFunc(guildID, eventID, after, limit)
}

func (f *Fake) SendDM(userID, message string) error {
	if f.SendDMFunc == nil {
		return nil
	}
	return f.SendDMFunc(userID, message)
}
