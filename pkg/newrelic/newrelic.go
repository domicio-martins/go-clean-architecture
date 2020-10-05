package newrelic

import (
	"context"
	"fmt"
	"time"

	"github.com/PicPay/picpay-dev-ms-template-manager/pkg/log"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type NewRelic struct {
	App         *newrelic.Application
	appName     string
	transaction *newrelic.Transaction
}

func Setup(appEnv, newRelicLicenseKey string) (appNewRelic NewRelic) {
	if newRelicLicenseKey == "" {
		log.Warn("Missing New Relic key", nil)
		return appNewRelic
	}

	appNewRelic.appName = fmt.Sprintf("template-manager:%s", appEnv)
	var err error
	appNewRelic.App, err = newrelic.NewApplication(
		newrelic.ConfigAppName(appNewRelic.appName),
		newrelic.ConfigLicense(newRelicLicenseKey),
	)
	if err != nil {
		log.Error("Error to init new relic application", err, nil)
		return appNewRelic
	}

	if err = appNewRelic.App.WaitForConnection(5 * time.Second); nil != err {
		log.Error("Error to connection to new relic", err, nil)
		return appNewRelic
	}
	log.Info("New Relic successfully initialized", &log.LogContext{
		"appName": appNewRelic.appName,
	})
	return appNewRelic
}

func (n *NewRelic) StartTransaction() context.Context {
	n.transaction = n.App.StartTransaction(n.appName)
	return newrelic.NewContext(context.Background(), n.transaction)
}

func (n *NewRelic) FinishTransaction() {
	n.transaction.End()
}
