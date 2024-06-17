// -*- coding: utf-8 -*-
// pgntest.go
// -----------------------------------------------------------------------------
//
// Started on <lun 03-06-2024 13:58:11.713581059 (1717415891)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

// Description
package testdata

import (
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"

	"math/rand"
)

// General purpose functions
// ----------------------------------------------------------------------------

// Return a string of length n randomly generated with runes from the given
// string
func RandString(n int, characters string) (output string) {

	// create the random generator with an arbitrary seed
	src := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	// Decode all runes in the given string
	nbrunes := utf8.RuneCountInString(characters)
	runes := make([]rune, nbrunes)
	for idx, rune := range characters {
		runes[idx] = rune
	}

	// create the string
	for len(output) < n {
		output += string(runes[src.Intn(nbrunes)])
	}

	// And return the string built
	return
}

// Remove n runes randomly from the given string and return the resulting string
func RandRemove(n int, input string) (output string) {

	// obviously, n can not be larger than the input string. If that happens,
	// return the empty string
	if n >= len(input) {
		return ""
	}

	// create the random generator with an arbitrary seed
	src := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	// Decode all runes in the given string
	nbrunes := utf8.RuneCountInString(input)
	runes := make([]rune, nbrunes)
	for idx, rune := range input {
		runes[idx] = rune
	}

	// Remove n characters randomly
	for rem := n; rem > 0; rem-- {

		// Randomly choose a location among the runes in the current input
		// string
		loc := src.Intn(nbrunes - (n - rem))

		// And remove it. I know, this is terrible but I want to preserve the order
		nextrunes := make([]rune, len(runes)-1)
		for idx := 0; idx < loc; idx++ {
			nextrunes[idx] = runes[idx]
		}
		for idx := loc + 1; idx < len(runes); idx++ {
			nextrunes[idx-1] = runes[idx]
		}

		// and update the runes
		runes = nextrunes
	}

	// Recreate the string with the resulting runes
	for _, rune := range runes {
		output += string(rune)
	}

	return
}

// Generates a random sequence of the first field of a FEN code which contains
// the piece placement for a single row. The result contains neither wildcards
// nor sequences of empty squares
func RandGenerateOneFENNoEmptySquares() (output string) {

	// chess symbols
	symbols := "prnbqkPRNBQK"

	// create the random generator with an arbitrary seed
	src := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	// generate the contents of up to 8 consecutive cells
	for length := 0; length < 8; length++ {

		// Then randomly choose a symbol (this is not necessarily correct
		// and there could be an arbitrary number of kings, for example)
		output += string(symbols[src.Intn(len(symbols))])
	}

	return
}

// Generates a random sequence of the first field of a FEN code which contains
// the piece placement for a single row. It does not use wildcards but might
// contain sequences of empty squares
func RandGenerateOneFEN() (output string) {

	// chess symbols
	symbols := "prnbqkPRNBQK"

	// create the random generator with an arbitrary seed
	src := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	// generate the contents of up to 8 consecutive cells considering also empty
	// cells. To generate alike fen codes, count the number of empty cells in
	// each row
	length := 0
	var digit bool

	// The first digit is randomly chosen between a symbol and a digit
	if src.Intn(100) >= 50 {

		// insert a symbol
		output += string(symbols[src.Intn(len(symbols))])
		digit = false
		length = 1
	} else {

		// Otherwise, insert a random number of empty squares
		squares := 1 + src.Intn(8)
		output += fmt.Sprintf("%v", squares)

		// and increment the length accordingly also remembering that the
		// last insertion was a digit
		digit = true
		length += squares
	}

	nbempty := 0
	for length < 8 {

		// First, randomly choose to take a number of consecutive empty cells or
		// not. In general, half the squares should be empty. Toss a coin and
		// accept symbols with probability 1-(length-nbempty)/length unless the
		// last symbol was already a digit ---consecutive digits are forbidden.
		// Note that we introduce a perturbation of 10%
		prob := src.Intn(100)
		if digit || prob >= max(0,
			int(100.0*(float64(length)-float64(nbempty))/float64(length))-10) {

			// Then randomly choose a symbol (this is not necessarily correct
			// and there could be an arbitrary number of kings, for example)
			output += string(symbols[src.Intn(len(symbols))])

			// and increment the length accordingly
			length++
		} else {

			// Otherwise, add a random number of consecutive empty squares
			squares := 1 + src.Intn(8-length)
			output += fmt.Sprintf("%v", squares)

			// and increment the length accordingly also remembering that the
			// last insertion was a digit
			digit = true
			length += squares
		}
	}

	return
}

// Generates a random sequence of the first field of a FEN code which contains
// the piece placement for an arbitrary number of rows. It does not use
// wildcards
func RandGenerateFullFEN(rows int) (output string) {

	for irow := 0; irow < rows; irow++ {

		// Create the contents of this row
		output += RandGenerateOneFEN()

		// In case this is not the last row, then add the slash
		if irow < rows-1 {
			output += "/"
		}
	}

	return
}

// Return a correct FEN code using wildcards that correctly represents the given
// FEN code which is assumed to consist of only one row. The input (explicit)
// FEN code might contain consecutive empty cells
func WildcardOneFEN(input string) (output string) {

	// create the random generator with an arbitrary seed
	src := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	// Because the input FEN code might contain empty cells, it is necessary to
	// explicitly count the number of cells in the FEN code
	length := 0
	for i := 0; i < len(input); i++ {
		if value, err := strconv.Atoi(string(input[i])); err == nil {
			length += value
		} else {
			length++
		}
	}

	// The procedure is simple, we just create a slice of bytes where we either
	// ignore a specific location or we do not. To do this we always substitute
	// or copy the first symbol in the input FEN code and update it in every
	// iteration to take into consideration the consecutive empty cells
	for i := 0; i < length; {

		// throw a dice to decide whether to mask this location or not
		if src.Intn(100) >= 50 {

			// Mask this position. In case this is a digit, then determine how
			// many locations to ignore
			if value, err := strconv.Atoi(string(input[0])); err == nil {

				locs := 1 + src.Intn(value)
				output += fmt.Sprintf("*%d", locs)

				// In case we are ignoring all locations, then move forward
				if locs == value {
					input = input[1:]
				} else {

					// Otherwise, consume the number of empty cells randomly
					// chosen
					input = fmt.Sprintf("%d", value-locs) + input[1:]
				}

				// and move forward the number of empty cells "wilcard"ed
				i += locs
			} else {

				// Otherwise, mask a single position
				output += "*"

				// since we are masking only position, if the rest of the output
				// starts with a digit then we are forced to add "1", otherwise,
				// we can freely choose to do so or not. It is impossible to
				// know, at this point, whether the output will continue with a
				// digit, so we use the input for determining what to do
				if len(input) > 1 {
					if _, err := strconv.Atoi(string(input[1])); err == nil {
						output += "1"
					} else {

						// in this case, randomly choose whether to add the trailing
						// 1 or not
						if src.Intn(100) >= 50 {
							output += "1"
						}
					}
				} else {

					// also, in this case, randomly choose
					if src.Intn(100) >= 50 {
						output += "1"
					}
				}

				// and move forward only one location
				input = input[1:]
				i++
			}
		} else {

			// Then we do not mask this position. First, are we starting with a
			// number?
			if value, err := strconv.Atoi(string(input[0])); err == nil {

				// If so, then determine how many consecutive cells to preserve
				locs := 1 + src.Intn(value)
				output += fmt.Sprintf("%d", locs)

				// In case we consumed them all, then move forward
				if locs == value {
					input = input[1:]
				} else {

					// Otherwise, consume the number of empty cells randomly
					// chosen
					input = fmt.Sprintf("%d", value-locs) + input[1:]
				}

				// and move forward the number of empty cells "wildcard"ed
				i += locs
			} else {

				// If this is not a digit, then just simply copy it
				output += string(input[0])

				// and move forward only one position
				input = input[1:]
				i++
			}
		}
	}

	// and return the string computed so far
	return
}

// Generates a random full sequence of the first field of a FEN code which
// contains the piece placement for an arbitrary number of rows. It returns two
// strings, the first one is the sequence without using wildcards, and the
// second one, uses wildcards and is guaranteed to match with the first one
func WildcardFullFEN(rows int) (contents, masked string) {

	for irow := 0; irow < rows; irow++ {

		// Create explicit contents of this row and then mask it
		icontents := RandGenerateOneFEN()
		contents += icontents
		masked += WildcardOneFEN(icontents)

		// In case this is not the last row, then add the slash
		if irow < rows-1 {
			contents += "/"
			masked += "/"
		}
	}

	return
}

// Local Variables:
// mode:go
// fill-column:80
// End:
