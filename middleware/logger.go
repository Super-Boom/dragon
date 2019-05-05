package middleware

import (
	"dragon/core/dragon/dlogger"
	"io/ioutil"
	"net/http"
	"time"
)

func LogInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// todo parse params will sometimes happen twice, because in controller will call it

		reqBody, _ := ioutil.ReadAll(r.Body)
		dlogger.SugarLogger.Infow("Request Info",
			"Method", r.Method,
			"Path", r.URL.Path,
			"Time", start.Format("2006-01-02 15:04:05"),
			"Params", string(reqBody),
		)
		next.ServeHTTP(w, r)
		dlogger.SugarLogger.Infow("Request Finish Info",
			"Method", r.Method,
			"Path", r.URL.Path,
			"Time", start.Format("2006-01-02 15:04:05"),
			"CostTime", time.Since(start).String(),
		)
		dlogger.Logger.Sync() // flushes buffer, if any
	})
}