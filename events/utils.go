package events

import (
	"context"
)

var (
	eventNotifyKey = "eventNotify"
)

// EventNotifyWithContext is that bind a event notify channel with a giving context
func EventNotifyWithContext(ctx context.Context, c chan bool) context.Context {
	return context.WithValue(ctx, eventNotifyKey, c)
}

// EventNotifyFromContext is that acquire a event notify channel from giving context
func EventNotifyFromContext(ctx context.Context) (chan bool, bool) {
	c, ok := ctx.Value(eventNotifyKey).(chan bool)
	return c, ok
}
