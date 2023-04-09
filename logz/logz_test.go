package logz

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"bitbucket.org/truedmp/logger"
	mockLogger "bitbucket.org/truedmp/logger/mock"
	"github.com/golang/mock/gomock"
)

const (
	jsonFullCheckMapStr = `{"empty":"","false":false,"float":777.7,"int":777,"int string":"777","null":null,"string":"string","true":true}`
)

var (
	jsonFullCheckMapObj = map[string]interface{}{"empty": "", "false": false, "float": 777.7, "int": 777, "int string": "777", "null": nil, "string": "string", "true": true}
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name string
		want Log
	}{
		{
			name: "New",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogger(); got == nil || got.GetData() == nil {
				t.Errorf("NewLogger() is nil, or data is nil = %v", got)
			}
		})
	}
}

func TestDefaultLog_getInfo(t *testing.T) {
	l := NewLogger()
	tests := []struct {
		name  string
		l     *DefaultLog
		want  *regexp.Regexp
		want2 *regexp.Regexp
		want3 *regexp.Regexp
	}{
		{
			name:  "Validate info",
			l:     l.(*DefaultLog),
			want:  regexp.MustCompile(`function:testing.tRunner`),
			want2: regexp.MustCompile(`file:.*?/testing\.go`),
			want3: regexp.MustCompile(`line:\d+`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmt.Sprintf("%+v", tt.l.getInfo()); !tt.want.Match([]byte(got)) || !tt.want2.Match([]byte(got)) || !tt.want3.Match([]byte(got)) {
				t.Errorf("DefaultLog.getInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeHeader() http.Header {
	header := http.Header{}
	header.Set(HeaderXCorrelationId, "777")
	header.Set(HeaderAuthorization, "My authorization")
	header.Set(HeaderCookie, "My secret cookie")
	header.Set(HeaderSetCookie, "My secret set-cookie")
	return header
}

func TestDefaultLog_Init(t *testing.T) {
	appWithoutReqBody := logger.Application{}
	appWithoutReqBody.XCorrelationID("777")
	appWithoutReqBody.FullPath("http://domain.org/path1/path2?q1=v1")
	appWithoutReqBody.Method("method1")
	appWithoutReqBody.ReqHeader("method1 /path1/path2?q1=v1 HTTP/0.0\r\nHost: domain.org\r\nAuthorization: 3DC8780\r\nCookie: 30F19DFB\r\nSet-Cookie: 5F175E6B\r\nX-Correlation-Id: 777\r\n\r\n")
	type args struct {
		r           *http.Request
		withReqBody bool
	}
	appWithReqBody := appWithoutReqBody
	appWithReqBody.ReqBody("My Body")
	header := makeHeader()
	headerWithoutXCorrelationId := makeHeader()
	headerWithoutXCorrelationId.Del(HeaderXCorrelationId)
	url, _ := url.Parse("http://domain.org/path1/path2?q1=v1")
	tests := []struct {
		name      string
		l         *DefaultLog
		args      args
		want      logger.Application
		wantMatch *regexp.Regexp
	}{
		{
			name: "withReqBody false",
			l:    NewLogger().(*DefaultLog),
			args: args{
				r: &http.Request{
					Method: "method1",
					URL:    url, //Don't forget to set or httputil will crash
					Header: header,
					Body:   io.NopCloser(strings.NewReader("My Body")),
					GetBody: func() (io.ReadCloser, error) {
						return io.NopCloser(strings.NewReader("My Body")), nil
					},
				},
				withReqBody: false,
			},
			want: appWithoutReqBody,
		},
		{
			name: "withReqBody true",
			l:    NewLogger().(*DefaultLog),
			args: args{
				r: &http.Request{
					Method: "method1",
					URL:    url, //Don't forget to set or httputil will crash
					Header: header,
					Body:   io.NopCloser(strings.NewReader("My Body")),
					GetBody: func() (io.ReadCloser, error) {
						return io.NopCloser(strings.NewReader("My Body")), nil
					},
				},
				withReqBody: true,
			},
			want: appWithReqBody,
		},
		{
			name: "without XCorrelationId",
			l:    NewLogger().(*DefaultLog),
			args: args{
				r: &http.Request{
					Method: "method1",
					URL:    url, //Don't forget to set or httputil will crash
					Header: headerWithoutXCorrelationId,
					Body:   io.NopCloser(strings.NewReader("My Body1")),
					GetBody: func() (io.ReadCloser, error) {
						return io.NopCloser(strings.NewReader("My Body1")), nil
					},
				},
				withReqBody: true,
			},
			wantMatch: regexp.MustCompile(`xCorrelationID:[0-f]{8}-[0-f]{4}-[0-f]{4}-[0-f]{4}-[0-f]{12}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantMatch != nil {
				if got := tt.l.Init(tt.args.r, tt.args.withReqBody).GetData(); !tt.wantMatch.Match([]byte(fmt.Sprintf("%+v", *got))) {
					t.Errorf("DefaultLog.Init() = \n%+v\n%+v", *got, tt.wantMatch)
				}
			} else {
				if got := tt.l.Init(tt.args.r, tt.args.withReqBody).GetData(); !reflect.DeepEqual(fmt.Sprintf("%+v", *got), fmt.Sprintf("%+v", tt.want)) {
					t.Errorf("DefaultLog.Init() = \n%+v\n%+v", *got, tt.want)
					t.Errorf("Binary:\n%v\n%v\n", []byte(fmt.Sprintf("%+v", *got)), []byte(fmt.Sprintf("%+v", tt.want)))
				}
			}
		})
	}
}

func TestDefaultLog_Copy(t *testing.T) {
	log1 := NewLogger().(*DefaultLog)
	log1.GetData().AccountID("accountId")
	log1.getInfo()
	tests := []struct {
		name string
		l    *DefaultLog
		want Log
	}{
		{
			name: "str message",
			l:    log1,
			// want: log2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Copy(); !reflect.DeepEqual(got, tt.l) {
				t.Errorf("DefaultLog.Copy() = %v, want %v", got, tt.want)
			} else if got == tt.l {
				t.Error("DefaultLog.Copy() is the same object")
			} else {
				l2 := got.(*DefaultLog)
				l2.BegTime = time.Now().Add(2 * time.Second)
				if l2.BegTime == tt.l.BegTime {
					t.Error("DefaultLog.Copy() must not share BegTime")
				}
				if l2.GetData() == tt.l.GetData() {
					t.Error("DefaultLog.Copy() must not share data object")
				}
				l2.GetData().AdditionalField("somethingnew", "yes")
				if !reflect.DeepEqual(l2.GetData(), tt.l.GetData()) {
					t.Error("DefaultLog.Copy() should share additionalField in data object because truedmp/logger has no getter")
				}
			}
		})
	}
}

func TestDefaultLog_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	AppLoggerMock := mockLogger.NewMockLogger(mockCtrl)
	AppLogger = AppLoggerMock

	type args struct {
		msg any
	}
	tests := []struct {
		name string
		l    *DefaultLog
		args args
		want string
	}{
		{
			name: "str message",
			l:    NewLogger().(*DefaultLog),
			args: args{
				msg: "log message",
			},
			want: "log message",
		},
		{
			name: "error object",
			l:    NewLogger().(*DefaultLog),
			args: args{
				msg: errors.New("err message"),
			},
			want: "err message",
		},
		{
			name: "error object",
			l:    NewLogger().(*DefaultLog),
			args: args{
				msg: jsonFullCheckMapObj,
			},
			want: jsonFullCheckMapStr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//tt.l.getInfo() will valid when gomock test because it's pointer. Our getInfo is tested so we'll ignore and let this pass.
			AppLoggerMock.EXPECT().Error(tt.want, tt.l.getInfo())
			tt.l.Error(tt.args.msg)
		})
	}
}

func TestDefaultLog_Info(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	AppLoggerMock := mockLogger.NewMockLogger(mockCtrl)
	AppLogger = AppLoggerMock

	type args struct {
		msg any
	}
	tests := []struct {
		name string
		l    *DefaultLog
		args args
		want string
	}{
		{
			name: "str message",
			l:    NewLogger().(*DefaultLog),
			args: args{
				msg: "log message",
			},
			want: "log message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AppLoggerMock.EXPECT().Info(tt.want, tt.l.getInfo())
			tt.l.Info(tt.args.msg)
		})
	}
}

func TestDefaultLog_Warning(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	AppLoggerMock := mockLogger.NewMockLogger(mockCtrl)
	AppLogger = AppLoggerMock

	type args struct {
		msg any
	}
	tests := []struct {
		name string
		l    *DefaultLog
		args args
		want string
	}{
		{
			name: "str message",
			l:    NewLogger().(*DefaultLog),
			args: args{
				msg: "log message",
			},
			want: "log message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AppLoggerMock.EXPECT().Warning(tt.want, tt.l.getInfo())
			tt.l.Warning(tt.args.msg)
		})
	}
}

func TestDefaultLog_GetData(t *testing.T) {
	tests := []struct {
		name string
		l    *DefaultLog
		want *logger.Application
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.GetData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultLog.GetData() = %v, want %v", got, tt.want)
			}
		})
	}
}
