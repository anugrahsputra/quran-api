package config

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
)

var Logger = logging.MustGetLogger("github.com/anugrahsputra/github.com/anugrahsputra/go-rest")

func ConfigureLogger() {
	fmt.Println("Configuring logger...")
	format := logging.MustStringFormatter(
		`%{color}[%{time:2006-01-02 15:04:05}] ▶ %{level}%{color:reset} %{message} ...[%{shortfile}@%{shortfunc}()]`,
	)

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)

	logging.SetBackend(backendFormatter)
}
