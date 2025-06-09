package cmd

import (
	"bufio"
	"fmt"
	"github.com/l-ILINDAN-l/WorkerPool/internal/pool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run pool workers in interactive mode",
	// TODO: Update help, do doc with variables
	Long: `This command run main cycle application, which use to manage workers and jobs`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Running application pool workers")

		// TODO: initial count workers from config file
		initialWorkers := 10
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
