package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/easy-oj/common/logs"
	"github.com/easy-oj/common/settings"
	"github.com/easy-oj/judger/initial"
)

func main() {
	var confPath, logsPath string
	flag.StringVar(&confPath, "conf", "/etc/eoj/settings.yaml", "Path of config file, default /etc/eoj/settings.yaml")
	flag.StringVar(&logsPath, "logs", "", "Path of logs file, default stdout")
	flag.Parse()
	settings.InitSettings(confPath)
	logs.InitLogs(logsPath)
	if bs, err := json.Marshal(settings.Settings); err == nil {
		logs.Info("[Main] Loaded settings: %s", string(bs))
	}
	initial.Initialize()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	_, _ = fmt.Fprintf(os.Stderr, "Received signal %v, exit...\n", <-ch)
}
