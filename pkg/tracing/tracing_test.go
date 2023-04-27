package tracing

import (
	"context"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"net/http/httptest"
	"testing"
)

func TestInit(t *testing.T) {
	_ = Init(Config{})
	_ = Init(Config{Address: "localhost:8000", Name: "-"})
	InitForTest()

	TraceID(context.Background())
	ctx := StartSpan(context.Background(), "startspan")
	EndSpan(ctx)
}

func TestClose(t *testing.T) {
	saved := tp
	tp = nil
	Close(context.Background())
	tp = new(tracesdk.TracerProvider)
	Close(context.Background())
	tp = saved
}

func TestAddEvent(t *testing.T) {
	AddEvent(context.Background(), "")
}

func TestSpanID(t *testing.T) {
	SpanID(context.Background())
}

func TestTraceID(t *testing.T) {
	TraceID(context.Background())
}

func TestCopySpan(t *testing.T) {
	CopySpan(context.Background(), context.TODO())
}

func TestInject(t *testing.T) {
	Inject(context.Background(), httptest.NewRequest("POST", "https://domain.com", nil))
}

func TestHttpReqTextMapCarrier(t *testing.T) {
	r := httptest.NewRequest("POST", "https://domain.com", nil)
	h := httpReqTextMapCarrier{h: &r.Header}
	h.Get("")
	h.Set("", "")
	h.Keys()
}
