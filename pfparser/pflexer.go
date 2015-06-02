/* 
  pflexer.go
  Description: returns different tokens of a propositional formulae to
  be parsed
  ----------------------------------------------------------------------------- 

  Started on  <Sat May 23 13:10:40 2015 Carlos Linares Lopez>
  Last update <martes, 02 junio 2015 17:46:28 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

package pfparser

import (
	"errors"		// for signaling errors
	"fmt"			// Sprintf
	"log"			// logging services
	"regexp"                // pgn files are parsed with a regexp
	"strconv"		// Atoi
)

// global variables
// ----------------------------------------------------------------------------

// the following regexps are used just to recognize different tokens that can
// appear in a propositional formula

// -- EOF: end of formula
var reEOF = regexp.MustCompile (`^\s*$`)

// -- opening parenthesis
var reOpenParen = regexp.MustCompile (`^\s*\(`)

// -- closing parenthesis
var reCloseParen = regexp.MustCompile (`^\s*\)`)

// -- integers
var reInteger = regexp.MustCompile (`^\s*(?P<value>[0-9]+)`)

// -- strings
var reString = regexp.MustCompile (`^\s*(?P<value>'[^']+')`)

// -- variables
var reVariable = regexp.MustCompile (`^\s*%(?P<varname>[a-zA-Z0-9_]+)`)

// -- relational operators
var reRelationalOperator = regexp.MustCompile (`^\s*(?P<operator>(<=|<|=|!=|>=|>))`)

// -- logical operators
var reLogicalOperator = regexp.MustCompile (`^\s*(?P<operator>(and|or|AND|OR))`)


// typedefs
// ----------------------------------------------------------------------------

// the type of a token is represented with an integer that is matched against
// consts (either digits or strings) or operators (both relational and logical)
type tokenType int

// a token is either a constant or an operator. Only in case it is a constant,
// its value is computed (otherwise, it shall be nil). Hence, values should
// satisfy the relational interface
type tokenItem struct {
	tokenType tokenType
	tokenValue RelationalEvaluator
}


// consts
// ----------------------------------------------------------------------------

// The type of a token can be any of the following: integer and string constants
// or relational or logical operators. Additionally, EOF (end of formula) is
// used as a token also to signal termination and parenthesis can be used to
// nest formulas
const (
	constInteger tokenType = 1 << iota	// integer constants
	constString				// string constants
	variable				// variables
	and					// -- logical operators
	or
	leq					// -- relational operators
	lt
	eq
	neq
	gt
	geq
	eof					// end of formula
	openParen				// parenthesis
	closeParen
)

// functions
// ----------------------------------------------------------------------------

// Return the next token in the propositional formula given in pformula and nil
// if any is successfully recognized, otherwise return nil and a syntax
// error. Additionally, if consume is true, the function modifies the string to
// point to the chunk to process in the next invocation
func nextToken (pformula *string, consume bool) (token tokenItem, err error) {

	// just apply regular expressions successively until one matches

	// -- EOF - End of Formula
	// --------------------------------------------------------------------
	if reEOF.MatchString (*pformula) {

		if consume {
			*pformula = ""
		}
		
		return tokenItem{eof, nil}, nil
		
	} else if reOpenParen.MatchString (*pformula) {

		// -- Opening parenthesis
		// ------------------------------------------------------------
		if consume {

			// process the string to locate the position of the
			// parenthesis and move forward
			tag := reOpenParen.FindStringSubmatchIndex (*pformula)
			*pformula = (*pformula)[tag[1]:]
		}
		
		return tokenItem{openParen, nil}, nil
		
	} else if reCloseParen.MatchString (*pformula) {

		// -- Closing parenthesis
		// ------------------------------------------------------------
		if consume {

			// process the string to locate the position of the
			// parenthesis and move forward
			tag := reCloseParen.FindStringSubmatchIndex (*pformula)
			*pformula = (*pformula)[tag[1]:]
		}
		
		return tokenItem{closeParen, nil}, nil
		
	} else if reInteger.MatchString (*pformula) {

		// -- Integer constants
		// ------------------------------------------------------------
		
		// process the string and extract the relevant group
		tag := reInteger.FindStringSubmatchIndex (*pformula)

		// convert this group to an integer value
		value, err := strconv.Atoi ((*pformula)[tag[2]:tag[3]]); if err != nil {
			return tokenItem{eof, nil}, errors.New ("It was not possible to process an integer")
		}
		
		// move forward in the propositional formula if required
		if consume {
			*pformula = (*pformula)[tag[3]:]
		}

		// and return a valid token
		return tokenItem{constInteger, ConstInteger (value)}, nil
				
	} else if reString.MatchString (*pformula) {

		// -- String constants
		// ------------------------------------------------------------
		
		// process the string and extract the relevant group
		tag := reString.FindStringSubmatchIndex (*pformula)

		// convert this group to a string value - note that
		// single quotes are automatically removed
		value := (*pformula)[1+tag[2]:tag[3]-1]
		
		// move forward in the propositional formula if required
		if consume {
			*pformula = (*pformula)[tag[3]:]
		}

		// and return a valid token
		return tokenItem{constString, ConstString (value)}, nil

	} else if reVariable.MatchString (*pformula) {

		// -- Variables
		// ------------------------------------------------------------
		
		// process the string and extract the relevant group
		tag := reVariable.FindStringSubmatchIndex (*pformula)

		// convert this group to a variable which just stores the name
		// of the variable
		value := (*pformula)[tag[2]:tag[3]]
		
		// move forward in the propositional formula if required
		if consume {
			*pformula = (*pformula)[tag[3]:]
		}

		// and return a valid token
		return tokenItem{variable, Variable (value)}, nil
		
	} else if reRelationalOperator.MatchString (*pformula) {

		// -- Relational operators
		// ------------------------------------------------------------
		
		// process the string and extract the relevant group
		tag := reRelationalOperator.FindStringSubmatchIndex (*pformula)

		// derive the type of the relational operator
		var relOp tokenType
		switch (*pformula)[tag[2]:tag[3]] {

		case "<":
			relOp = lt
		case "<=":
			relOp = leq
		case "=":
			relOp = eq
		case "!=":
			relOp = neq
		case ">":
			relOp = gt
		case ">=":
			relOp = geq
		default:
			log.Fatalf ("Unknown relational operator '%s'", (*pformula)[tag[2]:tag[3]])
		}

		// move forward in the propositional formula if required
		if consume {
			*pformula = (*pformula)[tag[1]:]
		}

		// and return a valid token
		return tokenItem {relOp, nil}, nil

	} else if reLogicalOperator.MatchString (*pformula) {

		// -- Logical operators
		// ------------------------------------------------------------
		
		// process the string and extract the relevant group
		tag := reLogicalOperator.FindStringSubmatchIndex (*pformula)

		// derive the type of logical operator
		var logop tokenType
		switch (*pformula)[tag[2]:tag[3]] {

		case "and":
			logop = and
		case "or":
			logop = or
		default:
			log.Fatalf ("Unknown logical operator '%s'", (*pformula)[tag[2]:tag[3]])
		}

		// move forward in the propositional formula
		if consume {
			*pformula = (*pformula)[tag[1]:]
		}

		// and return a valid token
		return tokenItem {logop, nil}, nil
	}

	// at this point, a syntax error happened, so that an
	// arbitrary token is returned in conjunction with an error
	// that points to the position in the string where the error
	// was found
	return tokenItem{and, nil}, errors.New (fmt.Sprintf("Syntax error in '%v'", *pformula))
}



/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
