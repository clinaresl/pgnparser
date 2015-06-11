/* 
  pgnparser.go
  Description: PGN parser
  ----------------------------------------------------------------------------- 

  Started on  <Sun May  3 23:44:57 2015 Carlos Linares Lopez>
  Last update <jueves, 11 junio 2015 11:08:26 Carlos Linares Lopez (clinares)>
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
	"flag"				// arg parsing
	"fmt"				// printing msgs
	"log"				// logging services
	"os"				// operating system services

	// import a user package to manage paths
	"bitbucket.org/clinares/pgnparser/fstools"

	// also use several tools for handling games in pgn format
	"bitbucket.org/clinares/pgnparser/pgntools"
)

// global variables
// ----------------------------------------------------------------------------
var VERSION string = "0.1.0"		// current version
var EXIT_SUCCESS int = 0		// exit with success
var EXIT_FAILURE int = 1		// exit with failure

var pgnfile string       		// base directory
var latexTemplate string		// file with the latex template
var query string			// select query to filter games
var helpExpressions bool		// is help on expressions requested?
var sort string                         // sorting descriptor
var helpSort bool                       // is help on sorting requested?
var verbose bool			// has verbose output been requested?
var version bool			// has version info been requested?

// functions
// ----------------------------------------------------------------------------

// init
// 
// initializes the command-line parser
// ----------------------------------------------------------------------------
func init () {

	// Flag to store the pgn file to parse
	flag.StringVar (&pgnfile, "file", "", "pgn file to parse. While this utility is expected to be generic, it specifically adheres to the format of ficsgames.org")

	// Flag to store the file with the LaTeX template
	flag.StringVar (&latexTemplate, "template", "", "file with a LaTeX template to use. If given, a file with the same name used in 'file' and extension '.tex' is automatically generated in the same directory where the pgn file resides. For more information on how to create and use LaTeX templates see the documentation")

	// Flag to receive a select query
	flag.StringVar (&query, "select", "", "if an expression is provided here, only games meeting it are accepted. For more information on expressions acknowledged by this directive use '--help-expressions'")
	flag.BoolVar (&helpExpressions, "help-expressions", false, "if given, additional information on expressions acknowledged by this application is provided")
	
	// Flag to receive a sorting descriptor
	flag.StringVar (&sort, "sort", "", "if a string is given here, games are sorted according to the sorting descriptor provided. For more information on sorting descriptors use '--help-sort'")
	flag.BoolVar (&helpSort, "help-sort", false, "if given, additional information on sorting descriptors is provided")
	
	// other optional parameters are verbose and version
	flag.BoolVar (&verbose, "verbose", false, "provides verbose output")
	flag.BoolVar (&version, "version", false, "shows version info and exists")
}

// showVersion
//
// shows version info and exists with the specified signal
// ----------------------------------------------------------------------------
func showVersion (signal int) {

	fmt.Printf ("\n %v", os.Args [0])
	fmt.Printf ("\n Version: %v\n\n", VERSION)
	os.Exit (signal)
}

// showEpressions
//
// shows informmation on expressions as they are recognized by the directive
// --select
// ----------------------------------------------------------------------------
func showExpressions (signal int) {

	fmt.Println (` 
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
 games played between January, 1, 2015 and June, 5, 2015. The following query:

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
	os.Exit (signal)
}


// showSortingDescriptors
//
// shows informmation on sorting descriptors acknowledged by the directive
// --sort
// ----------------------------------------------------------------------------
func showSortingDescriptors (signal int) {

	fmt.Println (` 
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
	os.Exit (signal)
}

// verify
// 
// parse the flags and verifies that proper values were given. If not, a fatal
// error is logged
// ----------------------------------------------------------------------------
func verify () {

	// first, parse the flags ---in case help was given, it is automatically
	// handled by the flags package
	flag.Parse ()

	// if version information was requested show it now and exit
	if version {
		showVersion (EXIT_SUCCESS)
	}

	// in case further assistance on either expressions or sorting is
	// requested, then show it here and exit
	if helpExpressions {
		showExpressions (EXIT_SUCCESS)
	}
	if helpSort {
		showSortingDescriptors (EXIT_SUCCESS)
	}

	// verify that the pgn file given exists and is accessible
	isregular, _ := fstools.IsRegular (pgnfile); if !isregular {
		log.Fatalf ("the pgn file '%s' does not exist or is not accessible",
			pgnfile)
	}
}


// Main body
// ----------------------------------------------------------------------------
func main () {

	// verify the values parsed
	verify ()

	// process the contents of the given file
	games := pgntools.GetGamesFromFile (pgnfile, query, sort, verbose)

	// show the headers of all games
	fmt.Printf ("\n")
	fmt.Println (games.ShowHeaders ())
	fmt.Printf ("\n")
	fmt.Printf (" # Games found: %v\n\n", games.Len ())

	// in case a LaTeX template has been given, then generate a LaTeX file
	// with the same name than the pgn file (and in the same location) with
	// extension '.tex'
	if latexTemplate != "" {

		// compute the contents to write to the file
		contents := games.GamesToLaTeXFromFile (latexTemplate)

		// now, write this contents to the specified file
		_, err := fstools.Write (pgnfile + ".tex", []byte (contents))
		if err != nil {
			log.Fatalf ("An error was issued when writing data to the LaTeX file")
		}
	}
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
