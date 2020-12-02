package middleware

import (
	"bytes"
	"dragon/core/dragon/tracker"
	"github.com/go-dragon/util"
	"io/ioutil"
	"net/http"
	"time"
)

func LogInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		spanId, _ := util.NewUUID()
		// 读取
		body, _ := ioutil.ReadAll(r.Body)
		// 把刚刚读出来的再写进去
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		trackMan := &tracker.Tracker{
			SpanId:    spanId,
			Uri:       r.RequestURI,
			Method:    r.Method,
			ReqHeader: r.Header,
			Body:      string(body),
			StartTime: start,
			DateTime:  start.Format("2006-01-02 15:04:05"),
			CostTime:  "",
		}
		trackInfo := trackMan.Marshal()
		r.Header.Set(tracker.TrackKey, trackInfo)

		// before_req hook
		beforeReq(r, w)

		next.ServeHTTP(w, r)

		// after_req hook
		afterReq(r, w)
	})
}
