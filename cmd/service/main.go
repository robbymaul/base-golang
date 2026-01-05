package main

import (
	appiinit "application/init"
	"os"
	"time"
)

func main() {
	// capture start time server
	startTime := time.Now()

	// init cli
	cli := appiinit.NewCLI(os.Args)

	// boot options for preparation
	boot := appiinit.HandleBootFlags(cli)

	// load cfg
	cfg, err := appiinit.Load(boot)
	if err != nil {
		panic(err)
	}

	if err := appiinit.UpgradeDB(cfg); err != nil {
		panic(err)
	}

	appInit := appiinit.InitVariables{
		StartedAt:  startTime,
		AppVersion: "N/A",
		Signature:  "N/A",
		Config:     cfg,
	}

	if cli.Start(appInit); cli.Error() != nil {
		panic(cli.Error())
	}
}
