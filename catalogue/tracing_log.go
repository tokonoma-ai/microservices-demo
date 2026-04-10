package catalogue

import (
	"context"
	"net/http"

	stdopentracing "github.com/opentracing/opentracing-go"
	zipkintracer "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go/model"
)

// TraceLogKV returns trace_id and span_id key-value pairs (B3 hex) for logfmt lines.
func TraceLogKV(ctx context.Context) []interface{} {
	sp := stdopentracing.SpanFromContext(ctx)
	if sp == nil {
		return []interface{}{"trace_id", "", "span_id", ""}
	}
	tid, sid := b3HexFromSpan(sp)
	return []interface{}{"trace_id", tid, "span_id", sid}
}

func b3HexFromSpan(sp stdopentracing.Span) (traceID, spanID string) {
	if zctx, ok := sp.Context().(zipkintracer.SpanContext); ok {
		mc := model.SpanContext(zctx)
		traceID = mc.TraceID.String()
		spanID = mc.ID.String()
		if traceID != "" && !mc.TraceID.Empty() && spanID != "" && mc.ID != 0 {
			return traceID, spanID
		}
	}
	m := make(map[string]string)
	if err := stdopentracing.GlobalTracer().Inject(sp.Context(), stdopentracing.TextMap, stdopentracing.TextMapCarrier(m)); err == nil {
		traceID, spanID = pickB3(m)
		if traceID != "" || spanID != "" {
			return traceID, spanID
		}
	}
	h := make(http.Header)
	if err := stdopentracing.GlobalTracer().Inject(sp.Context(), stdopentracing.HTTPHeaders, stdopentracing.HTTPHeadersCarrier(h)); err == nil {
		traceID = h.Get("X-B3-Traceid")
		if traceID == "" {
			traceID = h.Get("X-B3-TraceId")
		}
		spanID = h.Get("X-B3-Spanid")
		if spanID == "" {
			spanID = h.Get("X-B3-SpanId")
		}
	}
	return traceID, spanID
}

func pickB3(m map[string]string) (tid, sid string) {
	for k, v := range m {
		switch {
		case k == "x-b3-traceid" || k == "X-B3-TraceId":
			tid = v
		case k == "x-b3-spanid" || k == "X-B3-SpanId":
			sid = v
		}
	}
	return tid, sid
}
