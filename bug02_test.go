package outboxrelay_test

import (
	"testing"

	"github.com/LYH2263/go-outboxrelay"
)

func TestBug02_ListPendingPayloadSliceAlias(t *testing.T) {
	o := outboxrelay.New(outboxrelay.WithTargetURL("http://127.0.0.1:9/hook"))
	defer o.Close()
	id, err := o.Append("order.created", []byte(`{"id":1}`), map[string]string{"k": "v"})
	if err != nil {
		t.Fatal(err)
	}
	list, err := o.ListPending(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("len=%d", len(list))
	}
	list[0].Payload[2] = 'Z'
	list[0].Headers["k"] = "hacked"
	ev, err := o.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if string(ev.Payload) != `{"id":1}` {
		t.Fatalf("payload alias into store: %q", ev.Payload)
	}
	if ev.Headers["k"] != "v" {
		t.Fatalf("headers alias into store: %v", ev.Headers)
	}
}
