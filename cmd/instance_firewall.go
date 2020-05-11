package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"github.com/spf13/cobra"
)

var instanceSetFirewallCmd = &cobra.Command{
	Use:     "firewall",
	Aliases: []string{"set-firewall", "change-firewall"},
	Short:   "Use different firewall",
	Long: `Change an instance's firewall by part of the instance's ID or name and the full firewall ID.
If you wish to use a custom format, the available fields are:

	* ID
	* Hostname
	* FirewallID

Example: civo instance firewall ID/NAME 12345`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Printf("You must specify %s parameters (you gave %s), the ID/name and the firewall ID\n", utility.Red("2"), utility.Red(strconv.Itoa(len(args))))
			os.Exit(1)
		}

		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Unable to create a Civo API Client %s %s", err)
			os.Exit(1)
		}

		instance, err := client.FindInstance(args[0])
		if err != nil {
			utility.Error("Finding instance %s %s", err)
			os.Exit(1)
		}

		firewall, err := client.FindFirewall(args[1])
		if err != nil {
			utility.Error("Finding firewall %s %s", err)
			os.Exit(1)
		}

		_, err = client.SetInstanceFirewall(instance.ID, args[1])
		if err != nil {
			utility.Error("Setting firewall %s %s", err)
			os.Exit(1)
		}

		if outputFormat == "human" {
			fmt.Printf("Setting the firewall for the instance %s (%s) to %s (%s)\n", utility.Green(instance.Hostname), instance.ID, utility.Green(firewall.Name), firewall.ID)
		} else {
			ow := utility.NewOutputWriter()
			ow.StartLine()
			ow.AppendData("ID", instance.ID)
			ow.AppendData("Hostname", instance.Hostname)
			ow.AppendDataWithLabel("FirewallID", firewall.ID, "Firewall ID")
			if outputFormat == "json" {
				ow.WriteSingleObjectJSON()
			} else {
				ow.WriteCustomOutput(outputFields)
			}
		}
	},
}
