package cmd

import (
	"fmt"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"
	"github.com/gorhill/cronexpr"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

var snapshotListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list", "all"},
	Short:   "List snapshot",
	Long: `List all available snapshot.
If you wish to use a custom format, the available fields are:

	* ID
	* InstanceID
	* Hostname
	* Template
	* Region
	* Name
	* Safe
	* SizeGigabytes
	* State
	* Cron
	* RequestedAt
	* CompletedAt

Example: civo snapshot ls -o custom -f "ID: Name (Hostname)"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			fmt.Printf("Unable to create a Civo API Client: %s\n", aurora.Red(err))
			return
		}

		snapshots, err := client.ListSnapshots()
		if err != nil {
			fmt.Printf("Unable to list snapshots: %s\n", aurora.Red(err))
			return
		}

		ow := utility.NewOutputWriter()

		for _, snapshot := range snapshots {
			ow.StartLine()
			ow.AppendData("ID", snapshot.ID)
			ow.AppendData("Name", snapshot.Name)
			ow.AppendDataWithLabel("SizeGigabytes", fmt.Sprintf("%s GB", strconv.Itoa(snapshot.SizeGigabytes)), "Size")
			ow.AppendData("Hostname", snapshot.Hostname)
			ow.AppendData("State", snapshot.State)
			ow.AppendData("Cron", snapshot.Cron)
			if snapshot.Cron != "" {
				ow.AppendData("Schedule", cronexpr.MustParse(snapshot.Cron).Next(time.Now()).Format(time.RFC1123))
			} else {
				ow.AppendData("Schedule", "One-off")
			}
			ow.AppendData("RequestedAt", snapshot.RequestedAt.Format(time.RFC1123))
			ow.AppendData("CompletedAt", snapshot.CompletedAt.Format(time.RFC1123))

			if outputFormat == "json" || outputFormat == "custom" {
				ow.AppendData("InstanceID", snapshot.InstanceID)
				ow.AppendData("Template", snapshot.Template)
				ow.AppendData("Region", snapshot.Region)
				ow.AppendData("Safe", strconv.Itoa(snapshot.Safe))
			}

		}

		switch outputFormat {
		case "json":
			ow.WriteMultipleObjectsJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			ow.WriteTable()
		}
	},
}
