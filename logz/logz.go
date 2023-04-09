package logz

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"regexp"
	"runtime"
	"strings"
	"time"

	"bitbucket.org/truedmp/logger"
	"github.com/google/uuid"
	"github.com/zev-zakaryan/go-util/conv"
	"github.com/zev-zakaryan/go-util/stringz"
)

const (
	HeaderXCorrelationId = "X-Correlation-ID"
	HeaderAuthorization  = "Authorization"
	HeaderCookie         = "Cookie"
	HeaderSetCookie      = "set-cookie"
)

var (
	//AppLogger of truedmp/logger is singleton so we use singleton here. As var so we can inject mock.
	AppLogger           logger.Logger
	PatternHiddenHeader = `(?im)^((?:Authorization|Cookie|set-cookie):\s*)([^\r\n]*)`
)

type Log interface {
	Init(r *http.Request, withReqBody bool) Log
	Copy() Log
	Error(msg interface{})
	Info(msg interface{})
	Warning(msg interface{})
	GetData() *logger.Application
}

type DefaultLog struct {
	BegTime        time.Time
	XCorrelationId string
	Data           *logger.Application
}

func NewLogger() Log {
	return &DefaultLog{
		Data: logger.GetAppLog(),
	}
}

func (l *DefaultLog) getInfo() *logger.Application {
	function := ""
	pc, file, line, ok := runtime.Caller(2)
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			name := strings.Split(fn.Name(), "/")
			function = name[len(name)-1]
		}
	}

	l.Data.
		Function(function).
		AdditionalField("file", file).
		AdditionalField("line", fmt.Sprintf("%v", line)).
		RespTime(time.Since(l.BegTime).Milliseconds()).
		XCorrelationID(l.XCorrelationId)
	return l.Data
}

func (l *DefaultLog) Init(r *http.Request, withReqBody bool) Log {
	l.BegTime = time.Now().UTC()
	if l.XCorrelationId == "" && r != nil {
		l.XCorrelationId = strings.TrimSpace(r.Header.Get(HeaderXCorrelationId))
	}
	if l.XCorrelationId == "" {
		l.XCorrelationId = uuid.New().String()
	}
	l.Data.XCorrelationID(l.XCorrelationId) //Set when possible
	if r != nil {
		headers := ""
		if dump, err := httputil.DumpRequest(r, false); err == nil {
			reg := regexp.MustCompile(PatternHiddenHeader) //Exclude \r in hash, do not use .*$
			headers = stringz.ReplaceAllStringSubmatchFunc(reg, string(dump), func(ms []string) string {
				return ms[1] + stringz.ToCrc32(ms[2])
			})
		}
		//FullPath spec: Scheme+Host+r.URL.RequestURI(). But Scheme+Host maybe empty and omit from request object, we send as is.
		l.Data.FullPath(r.URL.String()).
			Method(r.Method).
			ReqHeader(headers)

		if withReqBody {
			bs, _ := io.ReadAll(r.Body)
			rc := io.NopCloser(bytes.NewBuffer(bs))
			r.Body = io.NopCloser(bytes.NewBuffer(bs))
			bs, _ = io.ReadAll(rc)
			l.Data.ReqBody(string(bs))
		}
	}
	return l
}

// Copy to copy log object with data. AdditionalField is private map so it's the same object (can't clone without reflect/unsafe).
func (l *DefaultLog) Copy() Log {
	out := *l
	newData := *l.Data
	out.Data = &newData
	return &out
}
func (l *DefaultLog) Error(msg interface{}) {
	AppLogger.Error(conv.ToForce[string](msg), l.getInfo())
}
func (l *DefaultLog) Info(msg interface{}) {
	AppLogger.Info(conv.ToForce[string](msg), l.getInfo())
}
func (l *DefaultLog) Warning(msg interface{}) {
	AppLogger.Warning(conv.ToForce[string](msg), l.getInfo())
}
func (l *DefaultLog) GetData() *logger.Application {
	return l.Data
}
