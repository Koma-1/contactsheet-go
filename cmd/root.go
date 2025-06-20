package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var config userConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "contactsheet-go",
	Short: "A CLI tool to generate contact sheets from image files",
	Long: `contactsheet-go is a command-line tool that generates contact sheets from a directory of image files.
It arranges images in a grid layout with customizable options such as tile size, number of rows and columns, margins, padding, interpolation method, and background colors.
You can choose how images fit into each tile: either by fitting them while preserving aspect ratio, or cropping them to fill the tile area.

Example:
	contactsheet-go -i ./images -o ./out -r 5 -c 4 -w 320 -h 320 --tilemode fit

This command creates a 5x4 contact sheet from images in ./images, with each tile resized to fit within 320x320 pixels, and saves the result in ./out.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(config)
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.contactsheet-go.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolP("help", "", false, "Help default flag")
	rootCmd.Flags().StringVarP(&config.indir, "indir", "i", "", "input directory")
	rootCmd.MarkFlagRequired("indir")
	rootCmd.Flags().StringVarP(&config.outdir, "outdir", "o", "", "output directory")
	rootCmd.MarkFlagRequired("outdir")
	rootCmd.Flags().StringVarP(&config.prefix, "prefix", "p", defaultConfig.prefix, "output file prefix")
	rootCmd.Flags().StringVarP(&config.interpolator, "interpolator", "", defaultConfig.interpolator, "Interpolator (a (ApproxBiLinear) | b (BiLinear) | c (CatmullRom) | n (NearestNeighbor))")
	rootCmd.Flags().StringVarP(&config.tileMode, "tilemode", "", defaultConfig.tileMode, "tilemode (fit | crop)")
	rootCmd.Flags().StringVarP(&config.backgroundColor, "background-color", "", defaultConfig.backgroundColor, "backgroundColor")
	rootCmd.Flags().StringVarP(&config.tileBackgroundColor, "tile-background-color", "", defaultConfig.tileBackgroundColor, "tileBackgroundColor")
	rootCmd.Flags().IntVarP(&config.rows, "rows", "r", defaultConfig.rows, "rows")
	rootCmd.Flags().IntVarP(&config.cols, "cols", "c", defaultConfig.cols, "cols")
	rootCmd.Flags().IntVarP(&config.width, "width", "w", defaultConfig.width, "width")
	rootCmd.Flags().IntVarP(&config.height, "height", "h", defaultConfig.height, "height")
	rootCmd.Flags().IntVarP(&config.innerMargin, "inner-margin", "", defaultConfig.innerMargin, "inner margin")
	rootCmd.Flags().IntVarP(&config.outerMargin, "outer-margin", "", defaultConfig.outerMargin, "outer margin")
	rootCmd.Flags().IntVarP(&config.padding, "padding", "", defaultConfig.padding, "padding")
}
