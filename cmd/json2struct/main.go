package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/EdgarOrtegaRamirez/json2struct/internal/analyzer"
	"github.com/EdgarOrtegaRamirez/json2struct/internal/generator"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	structName string
	language   string
	format     string
)

var rootCmd = &cobra.Command{
	Use:   "json2struct",
	Short: "Convert JSON data to typed structs for multiple languages",
	Long: `json2struct takes JSON input (from file, stdin, or URL) and generates
strongly-typed struct definitions for Go, Python, TypeScript, Rust, and Java.

Examples:
  echo '{"name":"John","age":30}' | json2struct --name User --lang go
  json2struct --name User --lang python data.json
  json2struct --name User --lang typescript < data.json`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
				os.Exit(1)
			}
			processInput(data)
		} else {
			data, err := os.ReadFile(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
				os.Exit(1)
			}
			processInput(data)
		}
	},
}

var outputCmd = &cobra.Command{
	Use:   "output <language> <input>",
	Short: "Generate code output for a specific language",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Usage: json2struct output <go|python|typescript|rust|java> <input_file_or_stdin>\n")
			os.Exit(1)
		}
		lang := generator.Language(args[0])
		if lang != generator.GoLang && lang != generator.Python && lang != generator.Typescript && lang != generator.Rust && lang != generator.Java {
			fmt.Fprintf(os.Stderr, "Error: unsupported language %q\n", args[0])
			os.Exit(1)
		}

		var data []byte
		if args[1] == "-" || args[1] == "stdin" {
			data, _ = io.ReadAll(os.Stdin)
		} else {
			data, _ = os.ReadFile(args[1])
		}
		processInputWithLang(data, lang)
	},
}

func processInput(data []byte) {
	lang := generator.Language(language)
	if lang == "" {
		fmt.Fprintln(os.Stderr, "Error: --lang flag is required. Use --help for supported languages.")
		os.Exit(1)
	}
	processInputWithLang(data, lang)
}

func processInputWithLang(data []byte, lang generator.Language) {
	var jsonData map[string]interface{}
	var err error

	if format == "yaml" {
		err = yaml.Unmarshal(data, &jsonData)
	} else {
		err = json.Unmarshal(data, &jsonData)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input: %v\n", err)
		os.Exit(1)
	}

	schema := analyzer.Analyze(jsonData)

	gen := generator.New(structName)
	code, err := gen.Generate(schema, lang)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(code)
}

func init() {
	rootCmd.Flags().StringVarP(&structName, "name", "n", "Data", "Name for the generated struct")
	rootCmd.Flags().StringVarP(&language, "lang", "l", "", "Target language (go, python, typescript, rust, java)")
	rootCmd.Flags().StringVarP(&format, "format", "f", "json", "Input format (json or yaml)")
	rootCmd.MarkFlagRequired("lang")

	rootCmd.AddCommand(outputCmd)
	outputCmd.Flags().StringVarP(&structName, "name", "n", "Data", "Name for the generated struct")
	outputCmd.Flags().StringVarP(&format, "format", "f", "json", "Input format (json or yaml)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}