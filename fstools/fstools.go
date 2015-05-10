/* 
  fstools.go
  Description: Simple tools for handling the filesystem paths
  ----------------------------------------------------------------------------- 

  Started on  <Thu Jun 19 13:36:57 2014 Carlos Linares Lopez>
  Last update <domingo, 10 mayo 2015 12:52:47 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

// fstools provides various simple services for handling paths and files. They
// are grouped in a different namespace since they are expected to be used often
// by other packages.
package fstools

import (
	"errors"		// for signaling errors
	"log"			// logging services
	"os"			// access to env variables
	"path"			// path manipulation
)

// global variables
// ----------------------------------------------------------------------------

// MAXLEN is the size of the largest block read at once when reading the
// contents of text files. By default, 1Kbyte
var MAXLEN int32 = 1024

// functions
// ----------------------------------------------------------------------------

// it returns an absolute path of the path given in dirin. It deals with strings
// starting with the symbol '~' and cleans the result (see path.Clean)
func ProcessDirectory (dirin string) (dirout string) {

	// initially, make the dirout to be equal to the dirin
	dirout = dirin

	// first, in case the input directory starts with the symbol
	// '~'
	if dirin [0] == '~' {

		// substitute '~' with the value of the $HOME variable
		dirout = path.Join (os.Getenv ("HOME"), dirin[1:])
	}

	// finally, clean the given directory specification
	dirout = path.Clean (dirout)

	return dirout
}


// returns true if the given path is a directory which is accessible to the user
// and false otherwise (thus, it is much like os.IsDir but it works from strings
// directly). It also returns a pointer to the os.File and its info in case they
// exist
func IsDir (path string) (isdir bool, filedir *os.File, fileinfo os.FileInfo) {

	var err error

	// open and stat the given location
	if filedir, err = os.Open (path); err!= nil {
		return false, nil, nil
	}
	if fileinfo, err = filedir.Stat (); err != nil {
		return false, filedir, nil
	}

	// return now whether this is a directory or not
	return fileinfo.IsDir (), filedir, fileinfo
}



// returns true if the given string names a regular file (ie., that no mode bits
// are set) and false otherwise (thus, it is much like os.IsRegular but it works
// from strings directly). It also returns the fileinfo in case the file exists
func IsRegular (path string) (isregular bool, fileinfo os.FileInfo) {

	var err error;
	
	// stat the specified path
	if fileinfo, err = os.Lstat (path); err != nil {
		return false, nil
	}
	
	// return now whether this is a regular file or not
	return fileinfo.Mode().IsRegular (), fileinfo
}


// returns a slice of bytes with the contents of the given file. If maxlen takes
// a positive value then data returns no more than max bytes. In case the file
// does not exist or it can not be accessed, a fatal error is raised
func Read (path string, maxlen int32) (contents []byte) {

	var err error

	// open the file in read access
	file, err := os.Open(path); if err != nil {
		log.Fatal(err)
	}

	// make sure the file is closed anyway
	defer file.Close ()
	
	// read the file in chunks of MAXLEN until EOF is reached or maxlen
	// bytes have been read
	var count int
	data := make([]byte, MAXLEN)

	for err == nil {
		count, err = file.Read (data)
		if err == nil {
			contents = append (contents, data[:count]...)
		}
	}
	
	// and return the data
	return contents
}

// write the contents specified in the given file. It returns the number of
// bytes written and nil if everything went fine. Otherwise, it returns any
// number and an error. In case it was not possible to create the file, a fatal
// error is raised
func Write (path string, contents []byte) (nbytes int, err error) {

	// check if the file exists
	if _, err = os.Stat(path); err == nil {
		return 0, errors.New ("the file already exists")
	}

	// now, open the file in read/write mode
	file, err := os.Create(path); if err != nil {
		log.Fatalf ("it was not possible to create the file")
	}

	// make sure the file is closed before leaving
	defer file.Close ()

	// and now write the contents into this file
	nbytes, err = file.Write(contents); if err != nil {
		log.Fatalf ("it was not possible to write to the file")
	}

	// syncing ...
	file.Sync ()

	return nbytes, nil
}


/* Local Variables: */
/* mode:go */
/* fill-column:80 */
/* End: */
