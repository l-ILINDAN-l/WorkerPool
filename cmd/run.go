package cmd

import (
	"bufio"
	"fmt"
	"github.com/l-ILINDAN-l/WorkerPool/internal/pool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

// runCmd represents the run command.
// It starts the worker pool in an interactive command-line mode.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run pool workers in interactive mode",
	Long: `This command starts the main application cycle, which allows you to manage workers and jobs.
Available commands in interactive mode:
  add - add a new worker
  remove - delete one worker.
  exit - shut down the application.
Any other text will be considered a processing job.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Running application pool workers")
		initialWorkers := viper.GetInt("workers.initial")
		p := pool.NewPool(initialWorkers)
		p.Start()

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Введите команду или задачу > ")
			if !scanner.Scan() {
				break
			}
			input := scanner.Text()
			command := strings.TrimSpace(input)

			switch command {
			case "exit":
				logrus.Info("Exiting application pool workers")
				p.Shutdown()
				time.Sleep(500 * time.Millisecond)
				return
			case "add":
				logrus.Info("Adding new worker")
				p.AddWorker()
			case "remove":
				logrus.Info("Removing worker")
				p.RemoveWorker()
			case "":
				continue
			default:
				logrus.WithFields(logrus.Fields{"job": command}).Infoln("Submit new job to pool")
				p.SubmitJob(command)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Scanner from stdin error: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
