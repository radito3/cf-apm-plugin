package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
)

type BlueGreenUploader struct{}

func (c *BlueGreenUploader) Run(cliConnection plugin.CliConnection, args []string) {
	if args[0] != "bg-upload" && len(args) != 3 {
		return
	}

	httpClient, err := createHttpClient(cliConnection)

	if err != nil {
		fmt.Println(err)
		return
	}

	opsClient := OperationsClient{*httpClient}

	if args[2] == "--continue" {
		opsClient.continueAppUpload(args[1])
	} else {
		opsClient.uploadApp(args[1], args[2])
	}

	monitor := OperationMonitor{args[1], httpClient}

	monitor.monitorOperation()
}

func (c *BlueGreenUploader) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "bgUploader",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "bg-upload",
				HelpText: "Upload or update an application without downtime",
				UsageDetails: plugin.Usage{
					Usage: "cf bg-upload APP_NAME [FILE_PATH|--continue]",
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(BlueGreenUploader))
}
