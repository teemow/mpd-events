package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	globalFlags struct {
		debug   bool
		verbose bool
	}

	mainCmd = &cobra.Command{
		Use:   "mpd-events",
		Short: "React on mpd events",
		Long:  "React on mpd events",
		Run:   mainRun,
	}

	projectVersion string
	projectBuild   string
)

func init() {
	mainCmd.PersistentFlags().BoolVarP(&globalFlags.debug, "debug", "d", false, "Print debug output")
	mainCmd.PersistentFlags().BoolVarP(&globalFlags.verbose, "verbose", "v", false, "Print verbose output")
}

func assert(err error) {
	if err != nil {
		if globalFlags.debug {
			fmt.Printf("%#v\n", err)
			os.Exit(1)
		} else {
			log.Fatal(err)
		}
	}
}

func confirm(question string) error {
	for {
		fmt.Printf("%s ", question)
		bio := bufio.NewReader(os.Stdin)
		line, _, err := bio.ReadLine()
		if err != nil {
			return err
		}

		if string(line) == "yes" {
			return nil
		}
		fmt.Println("Please enter 'yes' to confirm.")
	}
}

func mainRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func main() {
	mainCmd.AddCommand(versionCmd)
	mainCmd.AddCommand(listenCmd)

	mainCmd.Execute()
}
