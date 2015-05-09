/* 
  pgnparser.go
  Description: PGN parser
  ----------------------------------------------------------------------------- 

  Started on  <Sun May  3 23:44:57 2015 Carlos Linares Lopez>
  Last update <sábado, 09 mayo 2015 17:35:26 Carlos Linares Lopez (clinares)>
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

	// verify that the pgn file given exists and is accessible
	isregular, _ := fstools.IsRegular (pgnfile); if !isregular {
		log.Fatalf (" the pgn file '%s' does not exist or is not accessible",
			pgnfile)
	}
}


// Main body
// ----------------------------------------------------------------------------
func main () {

	// verify the values parsed
	verify ()

	// process the contents of the given file
	games := pgntools.GetGamesFromFile (pgnfile, verbose)

	// show the headers of all games
	fmt.Printf ("\n")
	fmt.Println (games.ShowHeaders ())
	fmt.Printf ("\n")
	fmt.Printf (" # Games found: %v\n\n", games.GetNbGames ())

	fmt.Println (games.GameToLaTeX ())
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
