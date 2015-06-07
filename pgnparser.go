/* 
  pgnparser.go
  Description: PGN parser
  ----------------------------------------------------------------------------- 

  Started on  <Sun May  3 23:44:57 2015 Carlos Linares Lopez>
  Last update <domingo, 07 junio 2015 16:33:07 Carlos Linares Lopez (clinares)>
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
	flag.StringVar (&pgnfile, "file", "", "pgn file to parse. This utility adheres to the format of ficsgames.org")

	// Flag to store the file with the LaTeX template
	flag.StringVar (&latexTemplate, "template", "", "file with a LaTeX template to use. If given, a file with the same name used in 'file' and extension '.tex' is automatically generated. This template acknowledges placeholders of the form '%name'. Acknowledged placeholders are PGN tags and, additionally, 'moves' which is substituted by the list of moves of each game")

	// Flag to receive a select query
	flag.StringVar (&query, "select", "", "if an expression is provided here, only games meeting it are accepted. For more information on expressions acknowledged by this directive use '--help-expressions'")
	flag.BoolVar (&helpExpressions, "help-expressions", false, "if given, additional information on expressions acknowledged by this application is provided")
	
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
                        <=, <, =, !=, >, >= 
 where a term can be either a constant or a variable. On one hand, constants can be
 either integer (such as 40) or strings (such as '1-0'). Note that strings have to be
 single quoted. Variables are preceded by the character '%'. Any tag appearing in the
 header of a PGN game can be used as a variable such as '%White' or '%WhiteElo'. 

 Logical expressions consist of relational groups related by any of the logical 
 operators:
                             and, or
 where 'and' has precedence over 'or'. To modify the precedence rules, parenthesis
 can be freely used. 

 Note that the names of variables and the logical operators are case sensitive.
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

	// in case further assistance on expressions is requested, then show it
	// here and exit
	if helpExpressions {
		showExpressions (EXIT_SUCCESS)
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
	games := pgntools.GetGamesFromFile (pgnfile, query, verbose)

	// show the headers of all games
	fmt.Printf ("\n")
	fmt.Println (games.ShowHeaders ())
	fmt.Printf ("\n")
	fmt.Printf (" # Games found: %v\n\n", games.GetNbGames ())

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
