/*
  pgnparser.go
  Description: PGN parser
  -----------------------------------------------------------------------------

  Started on  <Sun May  3 23:44:57 2015 Carlos Linares Lopez>
  Last update <martes, 29 marzo 2016 21:06:50 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

package main

// imports
// ----------------------------------------------------------------------------
import (
	"flag" // arg parsing
	"fmt"  // printing msgs
	"log"  // logging services
	"os"   // operating system services
	"time"

	// also use several tools for handling games in pgn format
	"github.com/clinaresl/pgnparser/pgntools"
)

// global variables
// ----------------------------------------------------------------------------
const VERSION string = "0.1.0" // current version
const AUTHOR string = "Carlos Linares López"
const EMAIL string = "carlos.linares@uc3m.es"

const TABLE_TEMPLATE = "templates/table/simple.tpl"

var EXIT_SUCCESS int = 0 // exit with success
var EXIT_FAILURE int = 1 // exit with failure

// Options
var filename string      // base directory
var list bool            // whether games should be listed or not
var play int = 0         // number of moves between boards
var filter string        // select query to filter games
var histogram string     // histogram descriptor
var sort string          // sorting descriptor
var output string        // name of the file that stores results
var tableTemplate string // file with the table template
var latexTemplate string // file with the latex template

var verbose bool // has verbose output been requested?
var version bool // has version info been requested?

// functions
// ----------------------------------------------------------------------------

// initializes the command-line parser
func init() {

	// Flag to store the pgn file to parse
	flag.StringVar(&filename, "file", "", "pgn file to parse. While this utility is expected to be generic, it specifically adheres to the format of ficsgames.org as used in lichess.org")

	// Flag to store the number of moves between boards
	flag.BoolVar(&list, "list", false, "if given, a table with general information about all games found in the PGN file is shown")

	// Flag to store the number of moves between boards
	flag.IntVar(&play, "play", 0, "if given, each game in the PGN file is played, and the chess board is shown between the number of consecutive plies given. The board is not shown by default")

	// Flag to request filtering games by some criteria
	flag.StringVar(&filter, "filter", "", "generates a new pgn file with those games satisfying the given filtering criteria. For information about the filtering criteria see the documentation.")

	// Flag to request sorting games by some criteria
	flag.StringVar(&sort, "sort", "", "generates a new pgn file with games sorted according to the given criteria. For information about the sorting criteria see the documentation.")

	// Flag to request generating histograms
	flag.StringVar(&histogram, "histogram", "", "generates a table with a summary about the given variables. For information about the histogram variables see the documentation.")

	// Flag to store the output filename
	flag.StringVar(&output, "output", "output.pgn", "name of the file where the result of any manipulations is stored. It is used only in case any of the directives --filter or --sort is given. By default, 'output.pgn'")

	// Flag to store the template to use to generate the ascii table
	flag.StringVar(&tableTemplate, "table", "", "file with an ASCII template that can be used to override the output shown by default. For more information on how to create and use these templates see the documentation")

	// Flag to store the file with the LaTeX template
	flag.StringVar(&latexTemplate, "latex", "", "file with a LaTeX template to use. If given, a file with the same name used in 'file' and extension '.tex' is automatically generated in the same directory where the pgn file resides. For more information on how to create and use LaTeX templates see the documentation")

	// other optional parameters are verbose and version
	flag.BoolVar(&verbose, "verbose", false, "provides verbose output")
	flag.BoolVar(&version, "version", false, "shows version info and exists")
}

// shows version info and exists with the specified signal
func showVersion(signal int) {

	fmt.Printf("\n %v", os.Args[0])
	fmt.Printf("\n Version: %v\n", VERSION)
	fmt.Printf("\n %v", AUTHOR)
	fmt.Printf("\n %v\n\n", EMAIL)
	os.Exit(signal)
}

// parse the flags and verifies that proper values were given. If not, a fatal
// error is logged
func verify() {

	// first, parse the flags ---in case help was given, it is automatically
	// handled by the flags package
	flag.Parse()

	// if version information was requested show it now and exit
	if version {
		showVersion(EXIT_SUCCESS)
	}

	// verify that a pgn file to examine was given. Note that this argument is
	// mandatory
	if len(filename) == 0 {
		log.Fatalf(" Error: a PGN file must be given with --file")
	}
}

// Main body
func main() {

	// verify the values parsed
	verify()

	// PgnFile
	// ------------------------------------------------------------------------
	// Create a new PgnFile
	start := time.Now()
	pgnfile, err := pgntools.NewPgnFile(filename)
	if err != nil {
		log.Fatalf(" Error: %v\n", err)
	}

	// Show information of the PgnFile provided by the user
	fmt.Println()
	fmt.Println(pgnfile)
	fmt.Printf(" [%v]\n", time.Since(start))
	fmt.Println()

	// Obtain all games in this file as a collection of PgnGames
	start = time.Now()
	games, err := pgnfile.Games()
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Printf(" %v games found\n", games.Len())
	}
	fmt.Printf(" [%v]\n", time.Since(start))
	fmt.Println()

	// List games
	// ------------------------------------------------------------------------
	// show a table with information of the games been processed. For this,
	// a template is used: tableTemplate contains the location of a default
	// template to use; others can be defined with --table
	if list || len(tableTemplate) > 0 {

		// In case a list was requested but no template is given to produce it,
		// use the default one
		if len(tableTemplate) == 0 {
			tableTemplate = TABLE_TEMPLATE
		}
		games.GamesToWriterFromTemplate(os.Stdout, tableTemplate)
	}

	// Play/verify games
	// ------------------------------------------------------------------------
	// Play all games unconditionally. This is necessary to verify that the
	// transcription of all games is correct. In case a strictly positive value
	// is given then the board is shown on the standard output
	start = time.Now()
	if err := games.Play(play, os.Stdout); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf(" Games verified!\n")
	fmt.Printf(" [%v]\n", time.Since(start))
	fmt.Println()

	// Filter games
	// ------------------------------------------------------------------------
	// In case it has been requested to filter games, do so
	if filter != "" {
		start = time.Now()
		if filtered, err := games.Filter(filter); err != nil {
			log.Fatalln(err)
		} else {
			fmt.Printf(" %v games filtered\n", filtered.Len())

			// and make the filtered collection the current one
			games = filtered
		}
		fmt.Printf(" [%v]\n", time.Since(start))
		fmt.Println()
	}

	// Sort games
	// ------------------------------------------------------------------------
	if sort != "" {
		start = time.Now()
		if sorted, err := games.Sort(sort); err != nil {
			log.Fatalln(err)
		} else {
			fmt.Printf(" %v games sorted\n", sorted.Len())

			// and make the sorted collection the current one
			games = sorted
		}
		fmt.Printf(" [%v]\n", time.Since(start))
		fmt.Println()
	}

	// In case either sorting and/or filter has been requested, write the result
	// in the output file
	if sort != "" || filter != "" {

		// Check first whether there are some games to write
		if games.Len() == 0 {
			fmt.Println(" No games to store!")
		} else {

			// In case there are effectively some games to store, do! Create the
			// file in write mode and then write the contents of the entire
			// collection
			stream, err := os.Create(output)
			defer stream.Close()
			if err != nil {
				log.Fatalln(err)
			} else {
				games.GetPGN(stream)
			}
		}
	}

	// Histogram
	// ------------------------------------------------------------------------
	if histogram != "" {
		start = time.Now()
		if pgnhistogram, err := games.GetHistogram(histogram); err != nil {
			log.Fatalln(err)
		} else {
			fmt.Println(*pgnhistogram)
		}
		fmt.Printf(" [%v]\n", time.Since(start))
		fmt.Println()
	}

	// LaTeX
	// ------------------------------------------------------------------------

	// in case a LaTeX template has been given, then generate a LaTeX file
	// with the same name than the pgn file (and in the same location) with
	// extension '.tex' from the contents given in the specified template
	if latexTemplate != "" {

		// Create a LaTeX file to write the output
		if latexStream, err := os.Create(output + ".tex"); err != nil {
			log.Fatalln(err)
		} else {
			games.GamesToWriterFromTemplate(latexStream, latexTemplate)
		}

	}
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
