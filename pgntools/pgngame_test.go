/*
  pgngame.go
  Description: Simple tools to handle a single game in PGN format
  -----------------------------------------------------------------------------

  Started on  <Sat May  9 16:59:21 2015 Carlos Linares Lopez>
  Last update <sábado, 07 mayo 2016 16:44:27 Carlos Linares Lopez (clinares)>
  -----------------------------------------------------------------------------

  $Id::                                                                      $
  $Date::                                                                    $
  $Revision::                                                                $
  -----------------------------------------------------------------------------

  Made by Carlos Linares Lopez
  Login   <clinares@atlas>
*/

package pgntools

import (
	"testing"

	"github.com/clinaresl/pgnparser/pgntools/testdata"
	"golang.org/x/exp/rand"
)

func Test_consumeUndefined(t *testing.T) {
	type args struct {
		n    int
		code string
	}
	tests := []struct {
		name    string
		args    args
		advance int
		digits  int
		wantErr bool
	}{

		// Consuming ordinary characters
		// --------------------------------------------------------------------
		{name: "byte",
			args:    args{n: 1, code: "p"},
			advance: 1,
			digits:  0,
			wantErr: false},

		{name: "byte",
			args:    args{n: 1, code: "pp"},
			advance: 1,
			digits:  0,
			wantErr: false},

		{name: "byte",
			args:    args{n: 1, code: "ppp"},
			advance: 1,
			digits:  0,
			wantErr: false},

		{name: "byte",
			args:    args{n: 2, code: "p"},
			advance: 1,
			digits:  0,
			wantErr: true},

		{name: "byte",
			args:    args{n: 2, code: "pp"},
			advance: 2,
			digits:  0,
			wantErr: false},

		{name: "byte",
			args:    args{n: 2, code: "ppp"},
			advance: 2,
			digits:  0,
			wantErr: false},

		{name: "byte",
			args:    args{n: 3, code: "p"},
			advance: 1,
			digits:  0,
			wantErr: true},

		{name: "byte",
			args:    args{n: 3, code: "pp"},
			advance: 2,
			digits:  0,
			wantErr: true},

		{name: "byte",
			args:    args{n: 3, code: "ppp"},
			advance: 3,
			digits:  0,
			wantErr: false},

		// consuming empty squares
		// --------------------------------------------------------------------
		{name: "digit",
			args:    args{n: 1, code: "1"},
			advance: 1,
			digits:  0,
			wantErr: false},

		{name: "digit",
			args:    args{n: 1, code: "2"},
			advance: 1,
			digits:  1,
			wantErr: false},

		{name: "digit",
			args:    args{n: 1, code: "3"},
			advance: 1,
			digits:  2,
			wantErr: false},

		{name: "digit",
			args:    args{n: 2, code: "1"},
			advance: 1,
			digits:  0,
			wantErr: true},

		{name: "digit",
			args:    args{n: 2, code: "2"},
			advance: 1,
			digits:  0,
			wantErr: false},

		{name: "digit",
			args:    args{n: 2, code: "3"},
			advance: 1,
			digits:  1,
			wantErr: false},

		{name: "digit",
			args:    args{n: 3, code: "1"},
			advance: 1,
			digits:  0,
			wantErr: true},

		{name: "digit",
			args:    args{n: 3, code: "2"},
			advance: 1,
			digits:  0,
			wantErr: true},

		{name: "digit",
			args:    args{n: 3, code: "3"},
			advance: 1,
			digits:  0,
			wantErr: false},

		{name: "digit",
			args:    args{n: 2, code: "1p"},
			advance: 2,
			digits:  0,
			wantErr: false},

		{name: "digit",
			args:    args{n: 2, code: "2p"},
			advance: 1,
			digits:  0,
			wantErr: false},

		{name: "digit",
			args:    args{n: 2, code: "3p"},
			advance: 1,
			digits:  1,
			wantErr: false},

		{name: "digit#06",
			args:    args{n: 3, code: "1p"},
			advance: 2,
			digits:  0,
			wantErr: true},

		{name: "digit#07",
			args:    args{n: 3, code: "2p"},
			advance: 2,
			digits:  0,
			wantErr: false},

		{name: "digit",
			args:    args{n: 3, code: "3p"},
			advance: 1,
			digits:  0,
			wantErr: false},

		{name: "digit#06",
			args:    args{n: 4, code: "1p"},
			advance: 2,
			digits:  0,
			wantErr: true},

		{name: "digit#07",
			args:    args{n: 4, code: "2p"},
			advance: 2,
			digits:  0,
			wantErr: true},

		{name: "digit",
			args:    args{n: 4, code: "3p"},
			advance: 2,
			digits:  0,
			wantErr: false},

		// Consuming up to the end of the row
		// --------------------------------------------------------------------
		{name: "slash",
			args:    args{n: 1, code: "ppp/"},
			advance: 1,
			digits:  0,
			wantErr: false},

		{name: "slash",
			args:    args{n: 2, code: "ppp/"},
			advance: 2,
			digits:  0,
			wantErr: false},

		{name: "slash",
			args:    args{n: 3, code: "ppp/"},
			advance: 3,
			digits:  0,
			wantErr: false},

		{name: "slash",
			args:    args{n: 4, code: "ppp/"},
			advance: 3,
			digits:  0,
			wantErr: true},

		{name: "slash",
			args:    args{n: 1, code: "1pp/"},
			advance: 1,
			digits:  0,
			wantErr: false},

		{name: "slash",
			args:    args{n: 2, code: "1pp/"},
			advance: 2,
			digits:  0,
			wantErr: false},

		{name: "slash",
			args:    args{n: 3, code: "1pp/"},
			advance: 3,
			digits:  0,
			wantErr: false},

		{name: "slash",
			args:    args{n: 4, code: "1pp/"},
			advance: 3,
			digits:  0,
			wantErr: true},

		{name: "slash",
			args:    args{n: 1, code: "2p/"},
			advance: 1,
			digits:  1,
			wantErr: false},

		{name: "slash",
			args:    args{n: 2, code: "2p/"},
			advance: 1,
			digits:  0,
			wantErr: false},

		{name: "slash",
			args:    args{n: 3, code: "2p/"},
			advance: 2,
			digits:  0,
			wantErr: false},

		{name: "slash",
			args:    args{n: 4, code: "2p/"},
			advance: 2,
			digits:  0,
			wantErr: true},

		{name: "slash",
			args:    args{n: 1, code: "3/"},
			advance: 1,
			digits:  2,
			wantErr: false},

		{name: "slash",
			args:    args{n: 2, code: "3/"},
			advance: 1,
			digits:  1,
			wantErr: false},

		{name: "slash",
			args:    args{n: 3, code: "3/"},
			advance: 1,
			digits:  0,
			wantErr: false},

		{name: "slash",
			args:    args{n: 4, code: "3/"},
			advance: 1,
			digits:  0,
			wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := consumeUndefined(tt.args.n, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("consumeUndefined() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.advance {
				t.Errorf("consumeUndefined() got advance = %v, want %v", got, tt.advance)
			}
			if got1 != tt.digits {
				t.Errorf("consumeUndefined() got digits = %v, want %v", got1, tt.digits)
			}
		})
	}
}

func Test_cardinalityUndefined(t *testing.T) {
	type args struct {
		expr string
	}
	tests := []struct {
		name        string
		args        args
		advance     int
		cardinality int
	}{

		// No undefined positions
		// --------------------------------------------------------------------
		{name: "Undefined 0",
			args:        args{expr: "p"},
			advance:     0,
			cardinality: 0},

		{name: "Undefined 0",
			args:        args{expr: "1"},
			advance:     0,
			cardinality: 0},

		{name: "Undefined 0",
			args:        args{expr: "2"},
			advance:     0,
			cardinality: 0},

		{name: "Undefined 0",
			args:        args{expr: "/"},
			advance:     0,
			cardinality: 0},

		// One undefined positions
		// --------------------------------------------------------------------
		{name: "Undefined 0",
			args:        args{expr: "*"},
			advance:     1,
			cardinality: 1},

		{name: "Undefined 0",
			args:        args{expr: "*1"},
			advance:     2,
			cardinality: 1},

		{name: "Undefined 0",
			args:        args{expr: "*p"},
			advance:     1,
			cardinality: 1},

		{name: "Undefined 0",
			args:        args{expr: "*11"},
			advance:     2,
			cardinality: 1},

		{name: "Undefined 0",
			args:        args{expr: "*/"},
			advance:     1,
			cardinality: 1},

		{name: "Undefined 0",
			args:        args{expr: "**"},
			advance:     1,
			cardinality: 1},

		// Two undefined positions
		// --------------------------------------------------------------------
		{name: "Undefined 0",
			args:        args{expr: "*2"},
			advance:     2,
			cardinality: 2},

		{name: "Undefined 0",
			args:        args{expr: "*2p"},
			advance:     2,
			cardinality: 2},

		{name: "Undefined 0",
			args:        args{expr: "*21"},
			advance:     2,
			cardinality: 2},

		{name: "Undefined 0",
			args:        args{expr: "*2/"},
			advance:     2,
			cardinality: 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := cardinalityUndefined(tt.args.expr)
			if got != tt.advance {
				t.Errorf("cardinalityUndefined() got advance = %v, want %v", got, tt.advance)
			}
			if got1 != tt.cardinality {
				t.Errorf("cardinalityUndefined() got cardinality = %v, want %v", got1, tt.cardinality)
			}
		})
	}
}

func Test_consumeDigits(t *testing.T) {
	type args struct {
		n    int
		expr string
	}
	tests := []struct {
		name      string
		args      args
		success   bool
		advance   int
		undefined int
		wantErr   bool
	}{

		// One digit
		// --------------------------------------------------------------------
		{name: "One digit",
			args:      args{n: 1, expr: "1"},
			success:   true,
			advance:   1,
			undefined: 0,
			wantErr:   false},

		{name: "One digit",
			args:      args{n: 1, expr: "*"},
			success:   true,
			advance:   1,
			undefined: 0,
			wantErr:   false},

		{name: "One digit",
			args:      args{n: 1, expr: "*1"},
			success:   true,
			advance:   2,
			undefined: 0,
			wantErr:   false},

		{name: "One digit",
			args:      args{n: 1, expr: "*2"},
			success:   true,
			advance:   2,
			undefined: 1,
			wantErr:   false},

		{name: "One digit",
			args:      args{n: 1, expr: "*3"},
			success:   true,
			advance:   2,
			undefined: 2,
			wantErr:   false},

		{name: "One digit",
			args:      args{n: 1, expr: "p"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   false},

		{name: "One digit",
			args:      args{n: 1, expr: "p1"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   false},

		{name: "One digit",
			args:      args{n: 1, expr: "p*"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   false},

		{name: "One digit",
			args:      args{n: 1, expr: "p*1"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   false},

		{name: "One digit",
			args:      args{n: 1, expr: "p/"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   false},

		{name: "One digit",
			args:      args{n: 1, expr: "/"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   true},

		// Two digits
		// --------------------------------------------------------------------
		{name: "Two digits",
			args:      args{n: 2, expr: "1"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   true},

		{name: "Two digits",
			args:      args{n: 2, expr: "2"},
			success:   true,
			advance:   1,
			undefined: 0,
			wantErr:   false},

		{name: "Two digits",
			args:      args{n: 2, expr: "3"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   true},

		{name: "Two digits",
			args:      args{n: 2, expr: "*"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   true},

		{name: "Two digits",
			args:      args{n: 2, expr: "*1"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   true},

		{name: "Two digits",
			args:      args{n: 2, expr: "*2"},
			success:   true,
			advance:   2,
			undefined: 0,
			wantErr:   false},

		{name: "Two digits",
			args:      args{n: 2, expr: "*3"},
			success:   true,
			advance:   2,
			undefined: 1,
			wantErr:   false},

		{name: "Two digits",
			args:      args{n: 2, expr: "p"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   false},

		{name: "Two digits",
			args:      args{n: 2, expr: "pp"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   false},

		{name: "Two digits",
			args:      args{n: 2, expr: "p1"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   false},

		{name: "Two digits",
			args:      args{n: 2, expr: "p*"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   false},

		{name: "Two digits",
			args:      args{n: 2, expr: "p*1"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   false},

		{name: "Two digits",
			args:      args{n: 2, expr: "p/"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   false},

		{name: "Two digits",
			args:      args{n: 2, expr: "/"},
			success:   false,
			advance:   0,
			undefined: 0,
			wantErr:   true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, err := consumeDigits(tt.args.n, tt.args.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("consumeDigits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.success {
				t.Errorf("consumeDigits() got success = %v, want %v", got, tt.success)
			}
			if got1 != tt.advance {
				t.Errorf("consumeDigits() got advance = %v, want %v", got1, tt.advance)
			}
			if got2 != tt.undefined {
				t.Errorf("consumeDigits() got undefined = %v, want %v", got2, tt.undefined)
			}
		})
	}
}

func Test_matchFENPiecePlacement(t *testing.T) {
	type args struct {
		expr      string
		code      string
		digits    int
		undefined int
	}

	// Definition of ad-hoc test cases
	// ------------------------------------------------------------------------
	tests := []struct {
		name string
		args args
		want bool
	}{

		{name: "SimplePositive",
			args: args{expr: "/",
				code:      "/",
				digits:    0,
				undefined: 0},
			want: true},
	}

	// Execution of ad-hoc cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchFENPiecePlacement(tt.args.expr, tt.args.code, tt.args.digits, tt.args.undefined); got != tt.want {
				t.Errorf("matchFENPiecePlacement() = %v, want %v", got, tt.want)
			}
		})
	}

	// Definition of random cases
	// ------------------------------------------------------------------------

	// Without wildcards
	//
	// Random generation of FEN codes with a different number of rows
	for rows := 1; rows <= 8; rows++ {

		for i := 0; i < 1000; i++ {

			// Randomly generate the piece placement for this number of rows
			fen := testdata.RandGenerateFullFEN(rows)

			// Create a random case that actually matches
			positivecase := struct {
				name string
				args args
				want bool
			}{
				name: "RandFullRowFENEqualNoWildcards",
				args: args{
					expr:      fen,
					code:      fen,
					digits:    0,
					undefined: 0,
				},
				want: true,
			}

			// and execute it
			t.Run(positivecase.name, func(t *testing.T) {
				if got := matchFENPiecePlacement(positivecase.args.expr,
					positivecase.args.code,
					positivecase.args.digits,
					positivecase.args.undefined); got != positivecase.want {
					t.Errorf("matchFENPiecePlacement() = %v, want %v", got, positivecase.want)
				}
			})

			// And now, modify some characters and verify they do not match
			removed := testdata.RandRemove(1+rand.Intn(len(fen)), fen)

			// Create a random case that actually matches
			negativecase := struct {
				name string
				args args
				want bool
			}{
				name: "RandFullRowFENDifferentNoWildcards",
				args: args{
					expr:      fen,
					code:      removed,
					digits:    0,
					undefined: 0,
				},
				want: false,
			}

			// and execute it
			t.Run(negativecase.name, func(t *testing.T) {
				if got := matchFENPiecePlacement(negativecase.args.expr,
					negativecase.args.code,
					negativecase.args.digits,
					negativecase.args.undefined); got != negativecase.want {
					t.Errorf("matchFENPiecePlacement() = %v, want %v", got, negativecase.want)
				}
			})
		}
	}

	// With wildcards
	//
	// Random generation of FEN codes with a different number of rows
	for rows := 1; rows <= 8; rows++ {

		for i := 0; i < 1000; i++ {

			// Randomly generate the piece placement for this number of rows
			fen, wld := testdata.WildcardFullFEN(rows)

			// Create a random case that actually matches
			positivecase := struct {
				name string
				args args
				want bool
			}{
				name: "RandFullRowFENEqualWildcards",
				args: args{
					expr:      wld,
					code:      fen,
					digits:    0,
					undefined: 0,
				},
				want: true,
			}

			// and execute it
			t.Run(positivecase.name, func(t *testing.T) {
				if got := matchFENPiecePlacement(positivecase.args.expr,
					positivecase.args.code,
					positivecase.args.digits,
					positivecase.args.undefined); got != positivecase.want {
					t.Logf("\t> expr: %v\n", positivecase.args.expr)
					t.Logf("\t> code: %v\n", positivecase.args.code)
					t.Errorf("matchFENPiecePlacement() = %v, want %v", got, positivecase.want)
				}
			})

			// And now, modify some characters and verify they do not match
			removed := testdata.RandRemove(1+rand.Intn(len(fen)), fen)

			// Create a random case that actually matches
			negativecase := struct {
				name string
				args args
				want bool
			}{
				name: "RandFullRowFENDifferentNoWildcards",
				args: args{
					expr:      fen,
					code:      removed,
					digits:    0,
					undefined: 0,
				},
				want: false,
			}

			// and execute it
			t.Run(negativecase.name, func(t *testing.T) {
				if got := matchFENPiecePlacement(negativecase.args.expr,
					negativecase.args.code,
					negativecase.args.digits,
					negativecase.args.undefined); got != negativecase.want {
					t.Errorf("matchFENPiecePlacement() = %v, want %v", got, negativecase.want)
				}
			})
		}
	}
}

// Local Variables:
// mode:go
// fill-column:80
// End:
