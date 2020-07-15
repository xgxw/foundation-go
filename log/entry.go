package log

import (
	"strings"

	"github.com/sirupsen/logrus"
)

/*
	Entry 与 Logger 的区别:
	当 Logger 记录某些信息后(如WithField), 就会变成一个 Entry, 表示在该Field下的日志记录.
*/

func NewEntry(logger *Logger) *Entry {
	return &Entry{logrus.NewEntry(logger.Logger)}
}

type RequestProtocol string

const (
	ProtocolHTTP RequestProtocol = "http"
	ProtocolGRPC RequestProtocol = "grpc"
	ProtocolTCP  RequestProtocol = "tcp"
)

// Entry is wrapper of logrus.Entry which provides some methods to log structured data.
type Entry struct {
	*logrus.Entry
}

// WithScope indicate the logging scope. Scope can be a function name, a job name,
// a request name, or even a project name.
func (entry *Entry) WithScope(scope string) *Entry {
	return entry.WithField("scope", scope)
}

// WithContext attach some context data to the logger. Context data is used to
// trace your requests or jobs.
func (entry *Entry) WithContext(fields ContextFields) *Entry {
	return entry.WithField("context", fields)
}

// WithServer attach some information about requests on the server. It's something
// like Nginx or Apache access log, but it is strictly structured.
func (entry *Entry) WithServer(protocol RequestProtocol, fields ServerFields) *Entry {
	fields.Kind = "server"
	return entry.WithField(string(protocol), fields)
}

// WithClient is much like WithServer except it is used to log outgoing requests.
func (entry *Entry) WithClient(protocol RequestProtocol, fields ClientFields) *Entry {
	fields.Kind = "client"
	return entry.WithField(string(protocol), fields)
}

// WithRequest is used to log request payload.
func (entry *Entry) WithRequest(protocol RequestProtocol, fields PayloadFields) *Entry {
	name := strings.Join([]string{string(protocol), "request"}, ".")
	return entry.WithField(name, fields)
}

// WithResponse is used to log response payload.
func (entry *Entry) WithResponse(protocol RequestProtocol, fields PayloadFields) *Entry {
	name := strings.Join([]string{string(protocol), "response"}, ".")
	return entry.WithField(name, fields)
}

// WithMetrics is used to log metrics data.
func (entry *Entry) WithMetrics(fields MetricsFields) *Entry {
	return entry.WithField("metrics", fields)
}

func (entry *Entry) WithProcessMetricsFields(fields ProcessMetricsFields) *Entry {
	return entry.WithField("monitor", fields)
}

// WithField is wrapper of logrus.Entry.WithField. You SHOULD NOT use it in your code.
func (entry *Entry) WithField(name string, value interface{}) *Entry {
	return &Entry{entry.Entry.WithField(name, value)}
}

// WithFields is wrapper of logrus.Entry.WithFields. You SHOULD NOT use it in your code.
func (entry *Entry) WithFields(fields Fields) *Entry {
	return &Entry{entry.Entry.WithFields(fields)}
}

// WithError is wrapper of logrus.Entry.WithError
func (entry *Entry) WithError(err error) *Entry {
	return &Entry{entry.Entry.WithError(err)}
}

// ExtractContext returns the context fields in the current entry, in case someone
// want to attach additional data to it.
func (entry *Entry) ExtractContext() ContextFields {
	if ctx, ok := entry.Data["context"]; ok {
		return ctx.(ContextFields)
	}
	return ContextFields{}
}

// ExtractServer returns the server fields in the current entry, in case someone
// want to attach additional data to it.
func (entry *Entry) ExtractServer(protocol RequestProtocol) ServerFields {
	if server, ok := entry.Data[string(protocol)]; ok {
		return server.(ServerFields)
	}
	return ServerFields{}
}

// ExtractClient returns the client fields in the current entry, in case someone
// want to attach additional data to it.
func (entry *Entry) ExtractClient(protocol RequestProtocol) ClientFields {
	if client, ok := entry.Data[string(protocol)]; ok {
		return client.(ClientFields)
	}
	return ClientFields{}
}
