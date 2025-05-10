package main

import (
	"fmt"
	"os"

	"github.com/meap/logwatch-llm/internal/config"
	"github.com/meap/logwatch-llm/internal/llm"
	"github.com/meap/logwatch-llm/internal/logwatch"
	"github.com/meap/logwatch-llm/internal/presenter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "logwatch-llm",
	Short: "Translate Logwatch output into human-readable HTML reports using LLMs.",
	Long:  `A CLI tool to convert Logwatch output into actionable, human-readable HTML reports using LLMs.`,
	Run: func(cmd *cobra.Command, args []string) {
		var inputData []byte
		var err error

		// Check if stdin is being piped
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			inputData, err = os.ReadFile("/dev/stdin")
			if err != nil {
				logrus.Errorf("Failed to read stdin: %v\n", err)
				os.Exit(1)
			}
		} else {
			logrus.Errorf("Please provide Logwatch output via stdin.")
			os.Exit(1)
		}

		report := logwatch.ParseSectionsFromString(string(inputData))

		logrus.Infof("Parsing Logwatch output for host %s", report.Report.Host)

		if len(report.Sections) == 0 {
			logrus.Error("No sections found in the input data.")
			os.Exit(1)
		}

		result, err := llm.AnalyzeLogwatchSections(viper.GetString("model"), report.Sections)
		if err != nil {
			logrus.Errorf("Error analyzing sections: %v", err)
			os.Exit(1)
		}

		html, err := presenter.ConvertMarkdownToHTML(result.Content)
		if err != nil {
			logrus.Errorf("Error converting markdown to HTML: %v", err)
			os.Exit(1)
		}

		email, err := presenter.ComposePlainEmailWithAttachments(presenter.EmailMessage{
			To:        viper.GetString("email"),
			Subject:   "Logwatch Report for " + report.Report.Host,
			PlainBody: "Say something nice :)",
			Attachments: []presenter.EmailAttachment{
				{
					Name:        "logwatch.log",
					Content:     []byte(inputData),
					ContentType: "text/plain",
				},
				{
					Name:        "logwatch-report.html",
					Content:     []byte(html),
					ContentType: "text/html",
				},
			},
		})

		if err != nil {
			logrus.Errorf("Error composing email: %v", err)
			os.Exit(1)
		}

		errr := presenter.SendRawEmail(email)
		if errr != nil {
			logrus.Errorf("Error sending email: %v", errr)
			os.Exit(1)
		}

		logrus.Infof("Email sent to %s", viper.GetString("email"))

		if viper.GetString("output") != "" {
			outputHTML(html, viper.GetString("output"))
		}
	},
}

func outputHTML(html string, outputPath string) {
	err := os.WriteFile(outputPath, []byte(html), 0644)
	if err != nil {
		logrus.Errorf("Error writing HTML to file: %v", err)
		os.Exit(1)
	}

	logrus.Infof("HTML report saved to %s", outputPath)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.toml in standard locations)")
	rootCmd.PersistentFlags().StringP("model", "m", "", "LLM model to use (e.g., gpt-4o)")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Path to output HTML file")
	rootCmd.PersistentFlags().StringP("email", "e", "", "Email address to send the report to")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose (debug) logging")

	viper.BindPFlag("model", rootCmd.PersistentFlags().Lookup("model"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("email", rootCmd.PersistentFlags().Lookup("email"))

	viper.SetEnvPrefix("LOGWATCH_LLM")
	viper.AutomaticEnv()
}

func initConfig() {
	config.SetupLogging()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath("/etc/logwatch-llm/")
		viper.AddConfigPath("$HOME/.config/logwatch-llm")
		viper.AddConfigPath(".")
	}

	_ = viper.ReadInConfig() // Ignore error if config not found

	// Set up logrus log level based on verbose flag
	if viper.GetBool("verbose") {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}
}

func main() {
	Execute()
}
