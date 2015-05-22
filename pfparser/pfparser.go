/* 
  pfparser.go
  Description: Parser of propositional formulae
  ----------------------------------------------------------------------------- 

  Started on  <Wed May 20 23:46:05 2015 Carlos Linares Lopez>
  Last update <sábado, 23 mayo 2015 01:07:46 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

// Implementation of a parser of propositional formulae with a
// remarkable exception: instead of propositions, binary predicates
// over variables and constants are allowed. Although variables and
// constants are allowed, this is not predicate logic since it does
// not acknowledge the usage of quantifiers.
//
// Variables are preceded with the character '%'. The constants that
// are currently supported are integers and strings and, indeed,
// variables can only take values of these types
//
// Example:
//
//    (((%name == "Roberto") OR
//      (%name == "Dario")   OR
//      (%name == "Adriana")) AND
//     (%age > 2))
//
// As shown in the example above, parenthesis are allowed as well to
// modify the precedence rules which are applied by default as
// follows:
//
// 1. OR has less precedence than AND so it is applied before it
// 
// 2. AND has more precedence than OR.
//
// Note that NOT is not implemented since all binary operators can be
// reversed as desired. The binary operations recognized by this
// parser are: <= < = != > >=
//
package pfparser

import (
	"errors"		// for signaling errors
	"log"			// logging services
	"regexp"                // pgn files are parsed with a regexp
)

// global variables
// ----------------------------------------------------------------------------

// the following regexps are used just to recognize different tokens
// that can appear in a propositional formula

// -- variables
var reVariable = regexp.MustCompile (`^\s*%(?P<varname>[a-zA-Z0-9_]+)`)

// -- strings
var reString = regexp.MustCompile (`^\s*"(?P<string>[^"]+)"`)

// -- integers
var reInteger = regexp.MustCompile (`^\s*(?P<integer>[0-9]+)`)

// the following regexp recognize binary operators
var reOperator = regexp.MustCompile (`^\s*(?P<operator>(<=|<|=|!=|>=|>))`)


// typedefs
// ----------------------------------------------------------------------------

// The evaluation of relational expressions is performed over operands
// that support relational operators. These operators can be described
// with either Equal and Less or a combination of both
type RelationalInterface interface {
	Less (right RelationalInterface) TypeBool
	Equal (right RelationalInterface) TypeBool
}

// The evaluation of logical expressions requires the ability to apply
// logical operations over them, specifically AND and OR.
type LogicalInterface interface {
	And (right LogicalInterface) TypeBool
	Or (right LogicalInterface) TypeBool
}

// ConstInteger represents a constant integer value
type ConstInteger int32

// ConstString represents a constant string value
type ConstString string

// Relational operators are represented with integers which are
// matched against the constants: LEQ, LT, EQ, NEQ, GT, GEQ
type RelationalOperator int

// Logical operators are represented with integers which are matched
// against the constants: AND, OR
type LogicalOperator int

// The result of a relational expression is an instance of a boolean
// type which is renamed as a TypeBool
type TypeBool bool

// A relational evaluator is an interface that requires the ability to
// produce items that can be compared with a relational operator,
// i.e., that they produce a RelationalInterface
type RelationalEvaluator interface {
	Evaluate () RelationalInterface
}

// A Logical evaluator is an interface that requires the ability to
// produce items that can be compared with a logical operator, i.e.,
// that they produce a LogicalInterface
type LogicalEvaluator interface {
	Evaluate () LogicalInterface
}

// A relational expression consists of a relational operator that is
// applied over items that can be compared with such operator. 
type RelationalExpression struct {
	root RelationalOperator
	children [2]RelationalEvaluator
}

// A logical expression consists of a logical operator that is applied
// over items that can be compared with such operator
type LogicalExpression struct {
	root LogicalOperator
	children [2]LogicalEvaluator
}

// constants
// ----------------------------------------------------------------------------

// A relational operator consists of any of the following: <= < = != > >=
const (
	LEQ RelationalOperator = 1 << iota	// less or equal than
	LT					// less than
	EQ					// equal
	NEQ					// not equal
	GT					// greater than
	GEQ					// greater or equal than
)

// A logical operator consists of any of the following: AND, OR
const (
	AND LogicalOperator = 1 << iota		// AND
	OR					// OR
)

// Methods
// ----------------------------------------------------------------------------

// Compare this integer with the one specified in right and return
// whether the first is less than the second
func (constant ConstInteger) Less (right RelationalInterface) TypeBool {

	var value ConstInteger
	var ok bool
	
	value, ok = right.(ConstInteger); if !ok {
		log.Fatal (" Type mismatch in ConstInteger.Less")
	}

	return int32 (constant) < int32 (value)
}

// Compare this integer with the one specified in right and return
// whether the first is equal to the second
func (constant ConstInteger) Equal (right RelationalInterface) TypeBool {

	var value ConstInteger
	var ok bool
	
	value, ok = right.(ConstInteger); if !ok {
		log.Fatal (" Type mismatch in ConstInteger.Equal")
	}

	return int32 (constant) == int32 (value)
}

// Compare this string with the one specified in right and return
// whether the first is less than the second
func (constant ConstString) Less (right RelationalInterface) TypeBool {

	var value ConstString
	var ok bool
	
	value, ok = right.(ConstString); if !ok {
		log.Fatal (" Type mismatch in ConstString.Less")
	}

	return string (constant) < string (value)
}

// Compare this string with the one specified in right and return
// whether the first is equal to the second
func (constant ConstString) Equal (right RelationalInterface) TypeBool {

	var value ConstString
	var ok bool
	
	value, ok = right.(ConstString); if !ok {
		log.Fatal (" Type mismatch in ConstString.Equal")
	}

	return string (constant) == string (value)
}

// Perform the logical AND of this instance with the one in right and
// return the result
func (operand TypeBool) And (right LogicalInterface) TypeBool {

	var value TypeBool
	var ok bool

	value, ok = right.(TypeBool); if !ok {
		log.Fatal (" Type mismatch in TypeBool.And")
	}
	
	return TypeBool (bool (operand) && bool (value))
}

// Perform the logical OR of this instance with the one in right and
// return the result
func (operand TypeBool) Or (right LogicalInterface) TypeBool {

	var value TypeBool
	var ok bool

	value, ok = right.(TypeBool); if !ok {
		log.Fatal (" Type mismatch in TypeBool.Or")
	}
	
	return TypeBool (bool (operand) || bool (value))
}

// The following methods implement the evaluation procedure over
// different types

// The evaluation of a constant integer returns the same constant integer
func (constant ConstInteger) Evaluate () RelationalInterface {
	return constant
}

// The evaluation of a string constant returns the same constant string
func (constant ConstString) Evaluate () RelationalInterface {
	return constant
}

// The evaluation of a boolean type (TypeBool) returns the same constant
func (constant TypeBool) Evaluate () LogicalInterface {
	return constant
}

// The evaluation of a relational expression is done in two steps:
// first, both children are evaluated and then the relational operator
// is applied.
func (expression RelationalExpression) Evaluate () LogicalInterface {

	var result TypeBool = false
	
	// first, evaluate both children
	lchild := expression.children [0].Evaluate ()
	rchild := expression.children [1].Evaluate ()

	// and now, depending upon the type of relational operator,
	// apply the right combination of Equal and Less
	switch expression.root {

	case LEQ:
		result = lchild.Less (rchild) || lchild.Equal (rchild)
		
	case LT:
		result = lchild.Less (rchild)

	case EQ:
		result = lchild.Equal (rchild)

	case NEQ:
		result = lchild.Less (rchild) || rchild.Less (lchild)

	case GT:
		result = rchild.Less (lchild)

	case GEQ:
		result = rchild.Less (lchild) || rchild.Equal (lchild)

	default:
		log.Fatal (" Unknown relational operator!")
	}

	// and return the result computed so far
	return result
}

// The evaluation of a logical expression is done in two steps: first,
// both children are evaluated and then the logical operator is
// applied.
func (expression LogicalExpression) Evaluate () TypeBool {

	var result TypeBool = false

	// first, evaluate both children
	lchild, rchild := expression.children [0].Evaluate (), expression.children [1].Evaluate ()

	// and now, depending upon the type of the logical operator,
	// apply the right combination of AND and OR
	switch expression.root {

	case AND:
		result = lchild.And (rchild)

	case OR:
		result = lchild.Or (rchild)

	default:
		log.Fatal (" Unknown logical operator")
	}

	// and return the result computed so far
	return result
}

// Functions
// ----------------------------------------------------------------------------

// This function effectively parses the contents of the string given
// in pformula and returns a valid RelationalExpression and nil if no
// errors were found or an invalid RelationalExpression and an error
// otherwise
func Parse (pformula string) (result RelationalExpression, err error) {

	// --- experiments with Relational Expressions
	
	var a ConstInteger = 100
	var b ConstInteger = 10
	
	expression1 := RelationalExpression{LEQ, [2]RelationalEvaluator{a, b}}
	expression2 := RelationalExpression{LT, [2]RelationalEvaluator{a, b}}
	expression3 := RelationalExpression{EQ, [2]RelationalEvaluator{a, b}}
	expression4 := RelationalExpression{NEQ, [2]RelationalEvaluator{a, b}}
	expression5 := RelationalExpression{GT, [2]RelationalEvaluator{a, b}}
	expression6 := RelationalExpression{GEQ, [2]RelationalEvaluator{a, b}}

	log.Printf (" %v <= %v: %v", a, b, expression1.Evaluate ())
	log.Printf (" %v <  %v: %v", a, b, expression2.Evaluate ())
	log.Printf (" %v == %v: %v", a, b, expression3.Evaluate ())
	log.Printf (" %v != %v: %v", a, b, expression4.Evaluate ())
	log.Printf (" %v >  %v: %v", a, b, expression5.Evaluate ())
	log.Printf (" %v >= %v: %v", a, b, expression6.Evaluate ())
	log.Println ()
	
	var c ConstString = "dario"
	var d ConstString = "roberto"

	expression11 := RelationalExpression{LEQ, [2]RelationalEvaluator{c, d}}
	expression12 := RelationalExpression{LT, [2]RelationalEvaluator{c, d}}
	expression13 := RelationalExpression{EQ, [2]RelationalEvaluator{c, d}}
	expression14 := RelationalExpression{NEQ, [2]RelationalEvaluator{c, d}}
	expression15 := RelationalExpression{GT, [2]RelationalEvaluator{c, d}}
	expression16 := RelationalExpression{GEQ, [2]RelationalEvaluator{c, d}}

	log.Printf (" %v <= %v: %v", c, d, expression11.Evaluate ())
	log.Printf (" %v <  %v: %v", c, d, expression12.Evaluate ())
	log.Printf (" %v == %v: %v", c, d, expression13.Evaluate ())
	log.Printf (" %v != %v: %v", c, d, expression14.Evaluate ())
	log.Printf (" %v >  %v: %v", c, d, expression15.Evaluate ())
	log.Printf (" %v >= %v: %v", c, d, expression16.Evaluate ())
	log.Println ()
	
	// --- experiments with Logical Expressions
	
	var e TypeBool = true
	var f TypeBool = true

	expression21 := LogicalExpression{AND, [2]LogicalEvaluator{e, f}}
	expression22 := LogicalExpression{OR, [2]LogicalEvaluator{e, f}}

	log.Printf (" %v AND %v: %v", e, f, expression21.Evaluate ())
	log.Printf (" %v OR  %v: %v", e, f, expression22.Evaluate ())
	log.Println ()

	// --- experiments with both Relational and Logical Expressions

	expression31 := LogicalExpression{AND, [2]LogicalEvaluator{expression1, expression2}}
	expression32 := LogicalExpression{OR,  [2]LogicalEvaluator{expression1, expression2}}
	expression33 := LogicalExpression{AND, [2]LogicalEvaluator{expression3, expression4}}
	expression34 := LogicalExpression{OR,  [2]LogicalEvaluator{expression3, expression4}}
	expression35 := LogicalExpression{AND, [2]LogicalEvaluator{expression5, expression6}}
	expression36 := LogicalExpression{OR,  [2]LogicalEvaluator{expression5, expression6}}

	log.Printf (" %v <= %v AND %v < %v: %v", a, b, a, b, expression31.Evaluate ())
	log.Printf (" %v <= %v OR  %v < %v: %v", a, b, a, b, expression32.Evaluate ())
	log.Printf (" %v == %v AND %v != %v: %v", a, b, a, b, expression33.Evaluate ())
	log.Printf (" %v == %v OR  %v != %v: %v", a, b, a, b, expression34.Evaluate ())
	log.Printf (" %v >  %v AND %v >= %v: %v", a, b, a, b, expression35.Evaluate ())
	log.Printf (" %v >  %v OR  %v >= %v: %v", a, b, a, b, expression36.Evaluate ())
	
	// --- exit	
	
	return expression1, errors.New ("Not implemented yet!")
}



/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
