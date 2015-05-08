PGNparser -- May, 2015


# Introduction #

This tool uses regular expressions to parse PGN files and extract
information from them.

# Install #

First, clone the repository with:

    hg clone ssh://hg@bitbucket.org/clinares/pgnparser

from within your `$GOPATH/src/` directory. To compile `pgnparser`:

    $ go build

Finally, install the binary with:

    $ go install

And you are ready to execute:

    $ pgnparser --help

anywhere from your filesystem provided that the environment variable
`$PATH` contains the path to your `$GOPATH/bin/` directory


# Usage #

`pgnparser` provides additional information with the commands `--help`
and `--version`


## Example ##

The following command reads the contents of the file `mygames.pgn` and
parse its contents. It shows then some information on the standard
output

    pgnparser --file=~/tmp/mygames.pgn

note that values can be given either using `=` or not


# License #

PGNparser is free software: you can redistribute it and/or modify it under
the terms of the GNU General Public License as published by the Free
Software Foundation, either version 3 of the License, or (at your
option) any later version.

PGNparser is distributed in the hope that it will be useful, but WITHOUT
ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License
for more details.

You should have received a copy of the GNU General Public License
along with PGNparser.  If not, see <http://www.gnu.org/licenses/>.


# Author #

Carlos Linares Lopez <carlos.linares@uc3m.es>

