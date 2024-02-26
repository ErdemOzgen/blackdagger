package cmd

import (
	"context"
	"github.com/ErdemOzgen/blackdagger/internal/agent"
	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/engine"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/client"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func execDAG(ctx context.Context, e engine.Engine, cmd *cobra.Command, args []string, dry bool) {
	params, err := cmd.Flags().GetString("params")
	checkError(err)

	loadedDAG, err := loadDAG(args[0], removeQuotes(params))
	checkError(err)

	err = start(ctx, e, loadedDAG, dry)
	if err != nil {
		log.Fatalf("Failed to start DAG: %v", err)
	}
}

func start(ctx context.Context, e engine.Engine, d *dag.DAG, dry bool) error {
	// TODO: remove this
	ds := client.NewDataStoreFactory(config.Get())

	a := agent.New(&agent.Config{DAG: d, Dry: dry}, e, ds)
	listenSignals(ctx, a)
	return a.Run(ctx)
}

type signalListener interface {
	Signal(os.Signal)
}

var (
	signalChan chan os.Signal
)

func listenSignals(ctx context.Context, a signalListener) {
	go func() {
		signalChan = make(chan os.Signal, 100)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			a.Signal(os.Interrupt)
		case sig := <-signalChan:
			a.Signal(sig)
		}
	}()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func removeQuotes(s string) string {
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
