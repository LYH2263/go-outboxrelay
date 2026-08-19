package outboxrelay_test

import (
	"testing"

	"github.com/LYH2263/go-outboxrelay"
)

func TestBug01_AppendPayloadSliceAlias(t *testing.T) {
	o := outboxrelay.New(outboxrelay.WithTargetURL("http://127.0.0.1:9/hook"))
	defer o.Close()
	payload := []byte(`{"id":1}`)
	id, err := o.Append("order.created", payload, nil)
	if err != nil {
		t.Fatal(err)
	}
	payload[2] = 'Z'
	ev, err := o.Get(id)
	if err != nil {
		t.Fatal(err)
	}
	if string(ev.Payload) != `{"id":1}` {
		t.Fatalf("store payload polluted by caller alias: %q", ev.Payload)
	}
}
