package main

import (
	"context"

	"github.com/livekit/protocol/auth"
	"github.com/nbd-wtf/go-nostr"
)

func rejectEvent(ctx context.Context, event *nostr.Event) (bool, string) {
	if event.Kind != 20312 {
		return true, "this relay is not writable"
	}
	tag := event.Tags.GetFirst([]string{"h", ""})
	if tag == nil {
		return true, "missing 'h' tag"
	}

	return false, ""
}

func handleEphemeral(ctx context.Context, event *nostr.Event) {
	if event.Kind != 20312 {
		return
	}
	tag := event.Tags.GetFirst([]string{"h", ""})
	if tag == nil {
		return
	}

	at := auth.NewAccessToken(s.LiveKitAPIKey, s.LiveKitAPISecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     (*tag)[1],
	}
}
