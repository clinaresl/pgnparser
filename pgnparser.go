/* 
  pgnparser.go
  Description: PGN parser
  ----------------------------------------------------------------------------- 

  Started on  <Sun May  3 23:44:57 2015 Carlos Linares Lopez>
  Last update <domingo, 24 mayo 2015 19:22:15 Carlos Linares Lopez (clinares)>
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

	// import the parser of propositional formulae
	"bitbucket.org/clinares/pgnparser/pfparser"
)

// global variables
// ----------------------------------------------------------------------------
var VERSION string = "0.1.0"		// current version
var EXIT_SUCCESS int = 0		// exit with success
var EXIT_FAILURE int = 1		// exit with failure

var pgnfile string       		// base directory
var latexTemplate string		// file with the latex template
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
	games := pgntools.GetGamesFromFile (pgnfile, verbose)

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

	expression0 := "(10>=1)"
	LogicalEvaluator0, err0 := pfparser.Parse (&expression0, 0)
	if err0 != nil {
		log.Fatalf ("%v\n", err0)
	}
	log.Println (LogicalEvaluator0)
	log.Printf (" Evaluation of expression #0: %v\n\n", LogicalEvaluator0.Evaluate ())

	expression1:= "10 >= 100 or 100 <= 1000"
	LogicalEvaluator1, err1 := pfparser.Parse (&expression1, 0)
	if err1 != nil {
		log.Fatalf ("%v\n", err1)
	}
	log.Println (LogicalEvaluator1)
	log.Printf (" Evaluation of expression #1: %v\n\n", LogicalEvaluator1.Evaluate ())

	expression2 := "'roberto' >= 'dario' and 'dario' != 'adriana'"
	LogicalEvaluator2, err2 := pfparser.Parse (&expression2, 0)
	if err2 != nil {
		log.Fatalf ("%v\n", err2)
	}
	log.Println (LogicalEvaluator2)
	log.Printf (" Evaluation of expression #2: %v\n\n", LogicalEvaluator2.Evaluate ())

	expression3 := "'roberto' >= 'dario' and 'dario' != 'adriana' or 'dario'>'monica'"
	LogicalEvaluator3, err3 := pfparser.Parse (&expression3, 0)
	if err3 != nil {
		log.Fatalf ("%v\n", err3)
	}
	log.Println (LogicalEvaluator3)
	log.Printf (" Evaluation of expression #3: %v\n\n", LogicalEvaluator3.Evaluate ())

	expression4 := " 3 <= 4 or 5<2 and 3>=2"
	LogicalEvaluator4, err4 := pfparser.Parse (&expression4, 0)
	if err4 != nil {
		log.Fatalf ("%v\n", err4)
	}
	log.Println (LogicalEvaluator4)
	log.Printf (" Evaluation of expression #4: %v\n\n", LogicalEvaluator4.Evaluate ())

	expression5 := " ( 3 <= 4 or 5<2 )and 3>=2"
	LogicalEvaluator5, err5 := pfparser.Parse (&expression5, 0)
	if err5 != nil {
		log.Fatalf ("%v\n", err5)
	}
	log.Println (LogicalEvaluator5)
	log.Printf (" Evaluation of expression #5: %v\n\n", LogicalEvaluator5.Evaluate ())

	expression6 := " 3 <= 4 or (5<2 and 3>=2)"
	LogicalEvaluator6, err6 := pfparser.Parse (&expression6, 0)
	if err6 != nil {
		log.Fatalf ("%v\n", err6)
	}
	log.Println (LogicalEvaluator6)
	log.Printf (" Evaluation of expression #6: %v\n\n", LogicalEvaluator6.Evaluate ())
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
