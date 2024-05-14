{{/*

	This is a bare template which provides almost no information,
	other than the name of both players, their elo score and the
	final result of each game.

	In other words, it provides ONLY bare information about each
	game

*/}}{{.GetTable " | l r | l r | c |" (.GetSlice "White" "WhiteElo" "Black" "BlackElo" "Result") }}
 # Games found: {{.Len}}
{{""}}
