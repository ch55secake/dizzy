// Package cmd provides command entrypoint
package cmd

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/ch55secake/dizzy/pkg/executor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd will build execution context based on the flags provided by the user and then pass it to the default executor,
// once passed to default executor it will be batched based on size and then passed further down to the requester
var rootCmd = &cobra.Command{
	Use:     "dizzy",
	Short:   "A sub-domain enumeration tool",
	Args:    cobra.ExactArgs(1), // Can extract url from here
	Aliases: []string{"diz", "di"},
	Example: "dizzy http://localhost:8080 -w /path/to/wordlist -l 1000 -X GET -H {'Accept': 'Application/JSON'} -t 10",
	Long: `                ___
           ____/ (_)_______  __  __
          / __  / /_  /_  / / / / /
         / /_/ / / / /_/ /_/ /_/ /
         \__,_/_/ /___/___/\__, /
                          /____/
          An unsung hero.    `,
	Run: func(cmd *cobra.Command, args []string) {

		wordlistFlag, _ := cmd.Flags().GetString("wordlist")
		methodFlag, _ := cmd.Flags().GetString("method")
		timeoutFlag, _ := cmd.Flags().GetInt("timeout")
		headersFlag, _ := cmd.Flags().GetString("headers")
		// TODO: Filter output by length
		lengthFlag, _ := cmd.Flags().GetInt("length")
		debugFlag, _ := cmd.Flags().GetBool("debug")

		var headers map[string]string
		if headersFlag != "" {
			err := json.Unmarshal([]byte(headersFlag), &headers)
			if err != nil {
				log.Fatalf("Error unmarshalling headers: %s", err)
			}
		}

		if debugFlag {
			logrus.SetLevel(logrus.DebugLevel)
		}

		ctx := executor.ExecutionContext{
			Filepath:       wordlistFlag,
			URL:            args[0],
			ResponseLength: lengthFlag,
			Timeout:        time.Duration(timeoutFlag) * time.Second,
			Method:         methodFlag,
			Headers:        headers,
		}

		executor.Execute(ctx)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Can define persistent flags which will be used and stored for the yaml file
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dizzy.yaml)")

	rootCmd.Flags().StringP("wordlist", "w", "", "provide wordlist to use")
	rootCmd.Flags().StringP("method", "X", "", "specify which http request method to use")
	rootCmd.Flags().Int32P("timeout", "t", 0, "specify timeout for each request")
	rootCmd.Flags().StringP("headers", "H", "", "specify headers to add to each request, accepted as json")
	rootCmd.Flags().Int32P("length", "l", 0, "filter output by length of response body")
	rootCmd.Flags().BoolP("debug", "d", false, "enable extra debug logging")
}
