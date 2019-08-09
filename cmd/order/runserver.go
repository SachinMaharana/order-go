package main

import (
	"fmt"
	"os"

	"github.com/SachinMaharana/isabella/order"
	"github.com/spf13/cobra"
)

func runServer(cmd *cobra.Command, args []string) {
	s, err := order.NewServer()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// s.InitPromReg()

	if err := s.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}

}
