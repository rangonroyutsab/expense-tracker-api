package utils

import (
	"fmt"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/getsentry/sentry-go"
)

// InitSentry initializes Sentry if a DSN is configured.
func InitSentry() error {
	dsn, _ := beego.AppConfig.String("sentry_dsn")
	if dsn == "" {
		return nil
	}

	return sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
	})
}

// CaptureError sends an error to Sentry if the error is not nil.
func CaptureError(err error) {
	if err != nil {
		sentry.CaptureException(err)
	}
}

// CapturePanicValue captures a recovered panic value.
func CapturePanicValue(value interface{}) {
	if value != nil {
		sentry.CaptureException(fmt.Errorf("panic: %v", value))
		sentry.Flush(2 * time.Second)
	}
}

// FlushSentry waits briefly for pending Sentry events before shutdown.
func FlushSentry() {
	sentry.Flush(2 * time.Second)
}
