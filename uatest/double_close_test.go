// +build integration

package uatest

import (
	"context"
	"testing"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

// TestAutoReconnection performs an integration test the auto reconnection
// from an OPC/UA server.
func TestConcurrentlyCreatingConnection(t *testing.T) {

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	srv := NewServer("reconnection_server.py")
	defer srv.Close()

	c := opcua.NewClient(srv.Endpoint, srv.Opts...)
	if err := c.Connect(ctx); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	f := func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				c.Connect(ctx)
				c.Call(&ua.CallMethodRequest{
					ObjectID:       ua.NewStringNodeID(2, "simulations"),
					MethodID:       ua.NewStringNodeID(2, "simulate_connection_failure"),
					InputArguments: []*ua.Variant{},
				})
			}
		}
	}

	go f()
	go f()

	<-ctx.Done()
}
