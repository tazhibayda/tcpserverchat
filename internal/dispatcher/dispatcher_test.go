package dispatcher

import (
	"context"
	"testing"
)

func TestJoinBroadcastLeave(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	d := New(ctx)
	go d.Run()

	c1 := Client{"c1", make(chan string, 1)}
	d.Join(c1)
	d.Broadcast("hello")
	got := <-c1.MsgChan
	if got != "hello" {
		t.Errorf("измещено: %q", got)
	}

	d.Leave(c1)
	cancel()
}
