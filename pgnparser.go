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

	// import a package to manage paths

	// also use several tools for handling games in pgn format
	"github.com/clinaresl/pgnparser/pgntools"
)

// global variables
// ----------------------------------------------------------------------------
const VERSION string = "0.1.0" // current version
const AUTHOR string = "Carlos Linares López"
const EMAIL string = "carlos.linares@uc3m.es"

var EXIT_SUCCESS int = 0 // exit with success
var EXIT_FAILURE int = 1 // exit with failure

// Options
var filename string      // base directory
var play int = 0         // number of moves between boards
var list bool            // whether games should be listed or not
var tableTemplate string // file with the table template
var latexTemplate string // file with the latex template
var filter string        // select query to filter games
var sort string          // sorting descriptor
var histogram string     // histogram descriptor

var helpFilter bool    // is help on filters requested?
var helpSort bool      // is help on sorting requested?
var helpHistogram bool // is help about histograms requested?
var verbose bool       // has verbose output been requested?
var version bool       // has version info been requested?

// functions
// ----------------------------------------------------------------------------

// initializes the command-line parser
func init() {

	// Flag to store the pgn file to parse
	flag.StringVar(&filename, "file", "", "pgn file to parse. While this utility is expected to be generic, it specifically adheres to the format of ficsgames.org as used in lichess.org")

	// Flag to store the number of moves between boards
	flag.IntVar(&play, "play", 0, "if given, each game in the PGN file is played, and the chess board is shown between the number of consecutive plies given. The board is not shown by default")

	// Flag to store the number of moves between boards
	flag.BoolVar(&list, "list", false, "if given, a table with general information about all games found in the PGN file is shown")

	// Flag to store the template to use to generate the ascii table
	flag.StringVar(&tableTemplate, "table", "templates/table/simple.tpl", "file with an ASCII template that can be used to override the output shown by default. For more information on how to create and use these templates see the documentation")

	// Flag to store the file with the LaTeX template
	flag.StringVar(&latexTemplate, "latex", "", "file with a LaTeX template to use. If given, a file with the same name used in 'file' and extension '.tex' is automatically generated in the same directory where the pgn file resides. For more information on how to create and use LaTeX templates see the documentation")

	// Flag to receive a filter query
	flag.StringVar(&filter, "filter", "", "if an expression is provided here, only games satisfying it are accepted. For more information on filters use '--help-filter'")
	flag.BoolVar(&helpFilter, "help-filter", false, "if given, additional information on expressions acknowledged by this application is provided")

	// Flag to receive a sorting descriptor
	flag.StringVar(&sort, "sort", "", "if a string is given here, games are sorted according to the sorting descriptor provided. For more information on sorting descriptors use '--help-sort'")
	flag.BoolVar(&helpSort, "help-sort", false, "if given, additional information on sorting descriptors is provided")

	// Flag to receive a histogram descriptor
	flag.StringVar(&histogram, "histogram", "", "if a string is given here, a histogram with the information requested is generated. For more information on how to specify histograms use '--help-histogram'")
	flag.BoolVar(&helpHistogram, "help-histogram", false, "if given, additional information on how histograms are specified is provided")

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

// shows informmation on expressions as they are recognized by the directive
// --select
func showFilters(signal int) {

	fmt.Println(` 
 Expressions are a powerful mechanism to filter games in a PGN file. They consist of
 logical expressions made of relational groups. 

 A relational group consists of two terms related by any of the relational operators 

                        <=, <, =, !=, >, >=, in, not_in

 where a term can be either a constant or a variable. 'in' and 'not_in' can be
 used only with string constants and they serve to verify whether the left term
 is a substring (or not) of the right term. 

 Constants can be either integer (such as 40) or strings (such as '1-0'). Note
 that strings have to be single quoted. Variables are preceded by the character
 '%'. Any tag appearing in the header of a PGN game can be used as a variable
 such as '%White' or '%WhiteElo'.

 Logical expressions consist of relational groups related by any of the logical 
 operators:

                                  and, or

 where 'and' has precedence over 'or'. To modify the precedence rules, parenthesis
 can be freely used. 

 Note that the names of variables and the logical operators are case sensitive.


 Examples:

 The file 'examples/ficsgamesdb_search_1255777.pgn', contains 2564 different
 games played between January, 1, 2015 and June, 5, 2015. The following filter:

    $ ./pgnparser --file examples/ficsgamesdb_search_1255777.pgn
                  --select "%Date <= '2015.01.31' and %Date >= '2015.01.01'"

 filters those games played during January. In total, 805 games. An alternative
 way to filter the same games is shown below:


    $ ./pgnparser --file examples/ficsgamesdb_search_1255777.pgn
                  --select "'2015.01' in %Date"

 To know how many won games are stored in the file by a specific player (in the
 example *clinares*):

    $ ./pgnparser --file examples/ficsgamesdb_search_1255777.pgn
                  --select "(%White = 'clinares' and %Result = '1-0') or 
                            (%Black = 'clinares' and %Result = '0-1')"

 which returns 1229 games.

 To know the number of games won/lost with ECO code C25 by the same player:

    $ ./pgnparser --file examples/ficsgamesdb_search_1255777.pgn
                  --select "((%White = 'clinares' and %Result = '1-0') or 
                             (%Black = 'clinares' and %Result = '0-1')) and
                            %ECO='C25'"

 returns 160 won games, and:

    $ ./pgnparser --file examples/ficsgamesdb_search_1255777.pgn
                  --select "((%White = 'clinares' and %Result = '0-1') or 
                             (%Black = 'clinares' and %Result = '1-0')) and
                            %ECO='C25'"

 returns 140 games.

`)
	os.Exit(signal)
}

// shows informmation on sorting descriptors acknowledged by the directive
// --sort
func showSortingDescriptors(signal int) {

	fmt.Println(` 
 Games can be sorted according to different criteria either in ascending or
 descending order. The keys to use are given as a string which consists of a
 sequence of variables (and hence, they should be preceded with the character
 '%'). A key is applied in increasing order if it is preceded by '<' and in
 decreasing order if it is given as '>'.

 Examples:

 The file 'examples/mygames.pgn', contains 5 games which can be sorted in
 increasing order of white's name as follows:

    $ ./pgnparser --file examples/mygames.pgn
                  --sort "< %White"

 In case two chess games have the same value of the key, ties can be broken
 using additional keys. For example:

    $ ./pgnparser --file examples/mygames.pgn
                  --sort "< %White > %Result"

 sorts games in increasing order of white's name and in decreasing order of the
 result for all games played by the same player as white.

 It is possible to use an arbitrary number of keys for sorting games.
`)
	os.Exit(signal)
}

// shows informmation on expressions as they are recognized by the directive
// --select
func showHistogram(signal int) {

	fmt.Println(` 
 Histograms are used to produce information about the frequencies of a variable
 or a combination of two variables. This is, histograms are limited to one or
 two variables.

 Two different types of variables are recognized:

 1. Variables prefixed with the character '%' and optionally with a title: 

                             title: variable

    Variables here refer mainly to tags defined in *all* PGN games or variables
    defined automatically by this software

    If a variable is given, histograms are computed as the number of ocurrences
    of each observed value of the specified variable.

 2. Cases defined with the following syntax:

              title: (title: expression ; title: expression ...)

    i.e., as a parenthesized sequence of expressions separated by semicolons
    ---for more information on how to define propositional formulas use
    --help-expressions. Every expression can be optionally preceded by a
    title. Likewise, the whole case can be also given a name which is also
    optional.

    In this case, the histogram is computed as the number of times that each
    expression is verified. In general, only one expression should be true. In
    case that more than one expression is evaluated to true a warning is
    automatically generated.

 If only one variable is provided, the histogram simply consists of the number
 of observations of the given variable/case. If two variables are given, the
 ocurrences of the second variable/case are indexed by the value of the first
 variable. Of course, histograms are processed according to the order of the
 variables.

 If a title is given, then it is used in the report generated. Otherwise, a
 verbatim copy of the variable/case definition is printed.

 Examples:


`)
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

	// in case further assistance on a particular subject is requested, then
	// show it here and exit
	if helpFilter {
		showFilters(EXIT_SUCCESS)
	}
	if helpSort {
		showSortingDescriptors(EXIT_SUCCESS)
	}
	if helpHistogram {
		showHistogram(EXIT_SUCCESS)
	}
}

// Main body
func main() {

	// verify the values parsed
	verify()

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

	// show a table with information of the games been processed. For this,
	// a template is used: tableTemplate contains the location of a default
	// template to use; others can be defined with --table
	if list {
		games.GamesToWriterFromTemplate(os.Stdout, tableTemplate)
	}

	// Play all games unconditionally. This is necessary to verify that the
	// transcription of all games is correct. In case a strictly positive value
	// is given then the board is shown on the standard output
	start = time.Now()
	games.Play(play, os.Stdout)
	fmt.Printf(" Games verified!\n")
	fmt.Printf(" [%v]\n", time.Since(start))
	fmt.Println()

	// In case it has been requested to filter games, do so
	if filter != "" {
		start = time.Now()
		if filtered, err := games.Filter(filter); err != nil {
			log.Fatalln(err)
		} else {
			fmt.Printf(" %v games filtered\n", filtered.Len())
		}
		fmt.Printf(" [%v]\n", time.Since(start))
	}
	fmt.Println()

	// // In case at least one histogram was given, then process it over the
	// // whole collection of pgn games
	// if histogram != "" {
	// 	hist := games.ComputeHistogram(histogram)
	// 	fmt.Printf("%v\n", &hist)
	// }

	// // in case a LaTeX template has been given, then generate a LaTeX file
	// // with the same name than the pgn file (and in the same location) with
	// // extension '.tex' from the contents given in the specified template
	// if latexTemplate != "" {

	// 	games.GamesToFileFromTemplate(pgnfile+".tex", latexTemplate)
	// }
}

/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
