// -*- coding: utf-8 -*-
// metatemplate.go
// -----------------------------------------------------------------------------
//
// Started on <vie 10-05-2024 22:51:51.586460886 (1715374311)>
// Carlos Linares López <carlos.linares@uc3m.es>
//

// Meta text/templates modify a few services from the text/template standard
// library to allow the usage of variables that are properly substituted before
// executing the template. The services provided in this package return ordinary
// text/templates that can then be processed with functions from template
// package.
package metatemplate

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// globals
// ----------------------------------------------------------------------------

// The following regexp looks for variables appearing in the metatemplate in the
// form ${variable} optionally followed by a prompt and a default value. The
// variable is a sequence of alphanumeric characters (both upper and lower case
// are allowed) and the underscore. The prompt and the default value can contain
// any character but ']'
var reTmplExtendedIdentifier = regexp.MustCompile(`\$(\{(?P<idname1>[a-zA-Z0-9_]+)(\[prompt:(?P<prompt>[^\]]+)\])?(\[default:(?P<default>[^\]]+)\])?\})`)

// types
// ----------------------------------------------------------------------------

// Meta-variables might be given either a prompt or a default value and
// certainly a name
type metaVar struct {
	name         string
	prompt       string
	defaultValue string
}

// so that metavars are defined as a dictionary indexed by the variable name
type metaVars map[string]metaVar

// functions
// ----------------------------------------------------------------------------

// process the contents of a match and return the definition of a new
// meta-variable with all fields provided by the user. The given string has to
// match the regular expression that matches meta-variables, i.e., it has to be
// known it is a correct description of a meta-variable
func getMetaVar(metavar string) metaVar {

	// get the different groups of this regular expression
	locs := reTmplExtendedIdentifier.FindAllStringSubmatchIndex(metavar, -1)

	// According to the regular expression, the different groups are found in
	// the following slices:
	//
	// [ 4: 5]: name
	// [ 8: 9]: prompt
	// [12:13]: default

	// the name is guaranteed to exist
	name := metavar[locs[0][4]:locs[0][5]]

	// in case a prompt has been given extract it
	var prompt string
	if locs[0][8] >= 0 {
		prompt = metavar[locs[0][8]:locs[0][9]]
	}

	// in case a default value was given, extract it as well
	var defaultVal string
	if locs[0][12] >= 0 {
		defaultVal = metavar[locs[0][12]:locs[0][13]]
	}

	// and finally return a meta-variable with all information extracted
	return metaVar{
		name:         name,
		prompt:       prompt,
		defaultValue: defaultVal,
	}
}

// The union of two meta-variables with the same name consists of the
// information in the first meta-var and, only if some field is empty the
// information from the second variable is used
func unionMetaVars(var1, var2 metaVar) (union metaVar) {

	// Copy the attributes of the first variable
	union = var1

	// Update the prompt and default value of the union if the respective field
	// in the first variable is empty
	if len(var1.prompt) == 0 {
		union.prompt = var2.prompt
	}
	if len(var1.defaultValue) == 0 {
		union.defaultValue = var2.defaultValue
	}

	// and return the union
	return
}

// The following function returns information about all meta-variables found in
// the given file. Meta-variables can be qualified with either a "prompt" or a
// "default" between square brackets after the name of the variable
func infoMetaVars(file io.Reader) metaVars {

	// initialize the output slice
	result := make(metaVars)

	// Looks in the given reader for all occurrences of meta-variables
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// examine this line
		contents := scanner.Text()
		for _, imetavar := range reTmplExtendedIdentifier.FindAllString(contents, -1) {

			// get this meta-variable and store it
			metavar := getMetaVar(imetavar)
			if value, ok := result[metavar.name]; !ok {
				result[metavar.name] = metavar
			} else {
				result[metavar.name] = unionMetaVars(value, metavar)
			}
		}
	}

	// and return the information of all meta vars
	return result
}

// The following function performs all the necessary operations to get the value
// of the given meta-variable and nil if no error was detected.
//
// If a default value is given, then it is used, unless a prompt has been given
// also. In this case the user is prompted with a default value which is then
// used in case RET is pressed, i.e., accepting the default value. If no default
// value has been given the user is prompted and the result is assigned to the
// variable. If neither a prompt nor a default value have been given an error is
// returned
func getValue(metavar metaVar) (string, error) {

	// In case a prompt was given, ask the user
	if len(metavar.prompt) > 0 {

		// The prompt to show the user must include the default value in case
		// any has been given in addition to the prompt
		userPrompt := metavar.prompt
		if len(metavar.defaultValue) > 0 {
			userPrompt += fmt.Sprintf(" (%v)", metavar.defaultValue)
		}

		// and ask the user
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf(" %v: ", userPrompt)
		scanner.Scan()
		if scanner.Err() != nil {
			return "", fmt.Errorf(" Error while reading the user input for prompt '%v'\n", userPrompt)
		}
		result := scanner.Text()

		// in case the user immediately pressed RET verify whether the empty
		// string has to be used or whether (s)he was accepting the default
		// value
		if len(result) == 0 {
			if len(metavar.defaultValue) >= 0 {
				result = metavar.defaultValue
			}
		}

		// and return the result
		return result, nil
	}

	// At this point, no prompt was given, so that just check whether a default
	// value was given before returning an error
	if len(metavar.defaultValue) >= 0 {
		return metavar.defaultValue, nil
	}

	// So, if neither a prompt nor a default value was given, then return any
	// value and an error
	return "", errors.New("No value")
}

// getValues returns a map of strings to strings with the substitions to perform
// in the template, and nil if no error occurred.
//
// This function accepts a map of strings to strings. In case the name of a
// meta-variable is found in this map, its value is given preference. Otherwise,
// its default value and/or its prompt are used
//
// If it was not possible to deduce the value of any meta-variable an error is
// returned
func getValues(values map[string]string, metavars metaVars) (substitutions map[string]string, err error) {

	substitutions = make(map[string]string)

	// process all variables
	for k, v := range metavars {

		// in case this name is also found in the dictionary of values, use it
		if value, ok := values[k]; ok {
			substitutions[k] = value
		} else {

			// in case it does not exist then try to deduce it from the prompt
			// and/or the default value in case any were given
			if value, err = getValue(v); err != nil {

				// In case it was not possible stop the process and return an error.
				return nil, fmt.Errorf(" No value found for variable '%v'\n", k)
			} else {

				// Otherwise, use the value deduced
				substitutions[k] = value
			}
		}
	}

	// and return the substitions computed so far
	return
}

// The following function provides a replacement of the function
// regexp.ParseFiles () providing a service for substituting meta-variables with
// the values in the given dictionary
func ParseFiles(values map[string]string, filenames ...string) (*template.Template, error) {

	// create a slice to store the processed files
	tmpfiles := make([]string, 0)

	// create temporary files with a copy of each input file with all
	// substitutions being performed
	for _, ifile := range filenames {

		// Open this file in read mode
		if istream, ierr := os.OpenFile(ifile, os.O_RDONLY, 0644); ierr != nil {
			return nil, fmt.Errorf(" Error opening file '%v': %v\n", ifile, ierr)
		} else {

			// First of all, parse the template and get information of all
			// meta-variables
			metavars := infoMetaVars(istream)

			// Now, compute all substitutions of all values found in the template
			substitutions, err := getValues(values, metavars)
			if err != nil {
				return nil, err
			}

			// And now process the entire file to write the result of performing
			// all substitutions in a temporary file
			if ostream, oerr := os.CreateTemp("", filepath.Base(ifile)); oerr != nil {
				return nil, fmt.Errorf(" It was not possible to create a temp file for '%v'\n", filepath.Base(ifile))
			} else {

				// create a writer to store the result of the substitions into
				// the temporary file
				writer := bufio.NewWriter(ostream)

				// Get the contents of this file, line per line
				if _, err = istream.Seek(0, 0); err != nil {
					return nil, fmt.Errorf(" Error when seeking the start of file '%v'\n", ifile)
				}
				scanner := bufio.NewScanner(istream)
				for scanner.Scan() {

					// for all meta-template variables appearing in this file
					contents := scanner.Text()

					// and look for the ocurrence of any meta-variable in this line
					line := contents
					for _, loc := range reTmplExtendedIdentifier.FindAllStringSubmatchIndex(contents, -1) {

						// Get the name of this occurrence and perform the
						// corresponding substitution
						line = strings.Replace(line, contents[loc[0]:loc[1]], substitutions[contents[loc[4]:loc[5]]], -1)
					}

					// and write it in the temp file
					if _, err := writer.WriteString(line + "\n"); err != nil {
						return nil, fmt.Errorf(" Error writing into the temp file '%v'\n", ostream.Name())
					}
				}

				// flush this writer and add this temporary file to the files to
				// process
				writer.Flush()
				tmpfiles = append(tmpfiles, ostream.Name())
			}
		}
	}

	// Finally, pass the processed files to the function in regexp and return
	// its values
	return template.ParseFiles(tmpfiles...)
}

// Local Variables:
// mode:go
// fill-column:80
// End:
