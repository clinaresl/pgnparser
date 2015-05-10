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

`pgnparser` has a mandatory argument that shall be given with `--file`
and shall consist of a path to a PGN file. If no more directives are
given, `pgnparser` prints out a text table with the headers of all
games found in the PGN file.

It also recognizes an optional argument with `--template`. It consists
of a path to a LaTeX template that is used to automatically generate a
LaTeX file with the transcription of all games found in the PGN
file. `pgnparser` acknowledges placeholders that are conveniently
substituted to produce the desired output. All tags are of the form
`%name`. Every *tag name* found in the header of a game can be used as
a placeholder (e.g., `%Date` is substituted by the value of the tag
`Date`). Additionally, `%moves` is substituted by a list of the moves
of a game. The output generated can be processed with `pdflatex`.

`pgnparser` provides additional information with the commands `--help`
and `--version`


## Example ##

The following command reads the contents of the file `mygames.pgn` and
parse its contents. It shows then some information on the standard
output (note that values can be given either using `=` or not):

    pgnparser --file=mygames.pgn

The output shown is then of the form:

    |  DBGameNo  | Date                | White                     | Black                     | ECO | Time  | Moves | Result |
    +------------+---------------------+---------------------------+---------------------------+-----+-------+-------+--------+
    |  375431273 | 2015.05.04 17:18:00 | samotnik           (1618) | clinares           (1520) | C20 | 180+0 |    60 |    1-0 |
    |  375458181 | 2015.05.05 11:21:00 | clinares           (1514) | walls              (1532) | C00 | 180+0 |    67 |    ½-½ |
    |  375505558 | 2015.05.06 14:10:00 | indianpool         (1582) | clinares           (1524) | A40 | 180+0 |    26 |    1-0 |
    |  375529641 | 2015.05.07 05:54:00 | clinares           (1517) | tsoymen            (1575) | D00 | 180+0 |    10 |    0-1 |
    |  375529748 | 2015.05.07 05:54:00 | tsoymen            (1582) | clinares           (1510) | A40 | 180+0 |    34 |    0-1 |
    +------------+---------------------+---------------------------+---------------------------+-----+-------+-------+--------+

    # Games found: 5

If `--template` is given then a LaTeX file using the specified
template. For example (file `templates/simple.tex` is included in this
distribution):

    pgnparser --file mygames.pgn --template templates/simple.pgn

will generate a file called `mygames.pgn.tex`. The following image
shows one of the pages generated:

![View of the PDF file](images/sample-1.png)


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

