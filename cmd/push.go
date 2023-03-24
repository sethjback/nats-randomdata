package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sethjback/nats-randomdata/pusher"
	"github.com/spf13/cobra"
)

var (
	userCreds string
	service   string
	interval  int
)

var pushCmd = &cobra.Command{
	Use:   "push [subject name]",
	Short: "start pushing data to nats",
	Args:  cobra.ExactArgs(1),
	RunE:  Push,
}

func init() {
	pushCmd.Flags().StringVarP(&userCreds, "creds", "c", "", "nsc path to user creds (nsd://<operator>/<account>/<user>")
	pushCmd.Flags().IntVarP(&interval, "interval", "i", 1, "how often to push new message in seconds")
	pushCmd.Flags().StringVarP(&service, "service", "s", "", "nats URL")
	pushCmd.MarkFlagRequired("creds")
	GetRootCommand().AddCommand(pushCmd)
}

func Push(cmd *cobra.Command, args []string) error {
	p, err := pusher.New(args[0], userCreds, service, interval)
	if err != nil {
		return err
	}

	err = p.Start()
	if err != nil {
		return err
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	p.Stop()
	fmt.Println("finished")
	return nil
}
