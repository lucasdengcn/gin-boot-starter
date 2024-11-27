package logging

import (
	"bytes"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Options for logger
type Options struct {
	//
	Name string

	// Custom logger
	Logger *zerolog.Logger

	// FieldsOrder defines the order of fields in output.
	FieldsOrder []string

	// FieldsExclude defines contextual fields to not display in output.
	FieldsExclude []string
}

var (
	NameFieldName          = "name"
	HostnameFieldName      = "hostname"
	ClientIPFieldName      = "client_ip"
	UserAgentFieldName     = "user_agent"
	TimestampFieldName     = zerolog.TimestampFieldName
	DurationFieldName      = "elapsed"
	MethodFieldName        = "method"
	PathFieldName          = "path"
	PayloadFieldName       = "payload"
	RefererFieldName       = "referer"
	statusCodeFieldName    = "status_code"
	DataLengthFieldName    = "data_length"
	BodyFieldName          = "body"
	TraceIdFieldName       = "trace_id"
	CorrelationIdFieldName = "correlation_id"
	SpanIdFieldName        = "span_id"
)

// Logger is a gin middleware which use zerolog
func Logger() gin.HandlerFunc {
	o := &Options{
		FieldsOrder: ginDefaultFieldsOrder(),
	}
	return LoggerWithOptions(o)
}

// LoggerWithOptions is a gin middleware which use zerolog
func LoggerWithOptions(opt *Options) gin.HandlerFunc {
	// List of fields
	if len(opt.FieldsOrder) == 0 {
		opt.FieldsOrder = ginDefaultFieldsOrder()
	}

	// Logger to use
	if opt.Logger == nil {
		opt.Logger = &log.Logger
	}

	//
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return func(ctx *gin.Context) {
		// get zerolog
		z := opt.Logger

		// return if zerolog is disabled
		if z.GetLevel() == zerolog.Disabled {
			ctx.Next()
			return
		}

		// before executing the next handlers
		begin := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		// parse traceId and log
		// correlation ID, if NO TRACING then use correlation ID
		// inject into global context
		injectTracingIDs(ctx, z)

		// Get payload from request
		var payload []byte
		if !opt.isExcluded(PayloadFieldName) {
			payload, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewReader(payload))
		}

		// Get a copy of the body
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
		ctx.Writer = w

		// executes the pending handlers
		ctx.Next()

		// after executing the handlers
		duration := time.Since(begin)
		statusCode := ctx.Writer.Status()

		//
		var event *zerolog.Event
		var eventError bool
		var eventWarn bool

		// set message level
		if statusCode >= 400 && statusCode < 500 {
			eventWarn = true
			event = z.Warn()
		} else if statusCode >= 500 {
			eventError = true
			event = z.Error()
		} else {
			event = z.Trace()
		}

		// add fields
		for _, f := range opt.FieldsOrder {
			// Name field
			if f == NameFieldName && !opt.isExcluded(f) && len(opt.Name) > 0 {
				event.Str(NameFieldName, opt.Name)
			}
			// Hostname field
			if f == HostnameFieldName && !opt.isExcluded(f) && len(hostname) > 0 {
				event.Str(HostnameFieldName, hostname)
			}
			// ClientIP field
			if f == ClientIPFieldName && !opt.isExcluded(f) {
				event.Str(ClientIPFieldName, ctx.ClientIP())
			}
			// UserAgent field
			if f == UserAgentFieldName && !opt.isExcluded(f) && len(ctx.Request.UserAgent()) > 0 {
				event.Str(UserAgentFieldName, ctx.Request.UserAgent())
			}
			// Method field
			if f == MethodFieldName && !opt.isExcluded(f) {
				event.Str(MethodFieldName, ctx.Request.Method)
			}
			// Path field
			if f == PathFieldName && !opt.isExcluded(f) && len(path) > 0 {
				event.Str(PathFieldName, path)
			}
			// Payload field
			if f == PayloadFieldName && !opt.isExcluded(f) && len(payload) > 0 {
				event.Str(PayloadFieldName, string(payload))
			}
			// Timestamp field
			if f == TimestampFieldName && !opt.isExcluded(f) {
				event.Time(TimestampFieldName, begin)
			}
			// Duration field
			if f == DurationFieldName && !opt.isExcluded(f) {
				var durationFieldName string
				switch zerolog.DurationFieldUnit {
				case time.Nanosecond:
					durationFieldName = DurationFieldName + "_ns"
				case time.Microsecond:
					durationFieldName = DurationFieldName + "_us"
				case time.Millisecond:
					durationFieldName = DurationFieldName + "_ms"
				case time.Second:
					durationFieldName = DurationFieldName
				case time.Minute:
					durationFieldName = DurationFieldName + "_min"
				case time.Hour:
					durationFieldName = DurationFieldName + "_hr"
				default:
					z.Error().Interface("zerolog.DurationFieldUnit", zerolog.DurationFieldUnit).Msg("unknown value for DurationFieldUnit")
					durationFieldName = DurationFieldName
				}
				event.Dur(durationFieldName, duration)
			}
			// Referer field
			if f == RefererFieldName && !opt.isExcluded(f) && len(ctx.Request.Referer()) > 0 {
				event.Str(RefererFieldName, ctx.Request.Referer())
			}
			// statusCode field
			if f == statusCodeFieldName && !opt.isExcluded(f) {
				event.Int(statusCodeFieldName, statusCode)
			}
			// DataLength field
			if f == DataLengthFieldName && !opt.isExcluded(f) && ctx.Writer.Size() > 0 {
				event.Int(DataLengthFieldName, ctx.Writer.Size())
			}
			// Body field
			if f == BodyFieldName && !opt.isExcluded(f) && len(w.body.String()) > 0 {
				event.Str(BodyFieldName, w.body.String())
			}
		}

		// Message
		message := ctx.Errors.String()
		if message == "" {
			message = "Request"
		}

		// post the message
		if eventError {
			event.Msg(message)
		} else if eventWarn {
			event.Msg(message)
		} else {
			event.Msg(message)
		}
	}
}

func injectTracingIDs(ctx *gin.Context, z *zerolog.Logger) {
	traceparent := ctx.Request.Header.Get("traceparent")
	traceID := ""
	spanID := ""
	correlationID := ""
	if traceparent != "" {
		parts := strings.Split(traceparent, "-")
		if len(parts) == 4 {
			traceID = parts[1]
			spanID = parts[2]
			ctx.Set(TraceIdFieldName, traceID)
			ctx.Set(SpanIdFieldName, spanID)
		}
	}

	if traceID == "" || spanID == "" {
		correlationID = xid.New().String()
		ctx.Set(CorrelationIdFieldName, correlationID)
	}

	z.UpdateContext(func(c zerolog.Context) zerolog.Context {
		zc := c
		if correlationID != "" {
			zc = zc.Str(CorrelationIdFieldName, correlationID)
		}
		if traceID != "" {
			zc = zc.Str(TraceIdFieldName, traceID)
		}
		if spanID != "" {
			zc = zc.Str(SpanIdFieldName, spanID)
		}
		return zc
	})
	if correlationID != "" {
		ctx.Header("X-Correlation-ID", correlationID)
	}
}

// gormDefaultFieldsOrder defines the default order of fields
func ginDefaultFieldsOrder() []string {
	return []string{
		NameFieldName,
		HostnameFieldName,
		ClientIPFieldName,
		CorrelationIdFieldName,
		TraceIdFieldName,
		UserAgentFieldName,
		MethodFieldName,
		PathFieldName,
		PayloadFieldName,
		TimestampFieldName,
		DurationFieldName,
		RefererFieldName,
		statusCodeFieldName,
		DataLengthFieldName,
		BodyFieldName,
	}
}

// isExcluded check if a field is excluded from the output
func (o *Options) isExcluded(field string) bool {
	if o.FieldsExclude == nil {
		return false
	}
	for _, f := range o.FieldsExclude {
		if f == field {
			return true
		}
	}

	return false
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r responseBodyWriter) WriteString(s string) (n int, err error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}
