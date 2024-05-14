{{/*

	This is a simple template which provides additional (purely
	administrative) information in addition to the name of both
	players, their elo score and the final result.

	Games from lichess do not record the current time when the
	game was played, whereas fics does. Hence, whether the tag
	"Time" is defined or not is taken into account to set up the
	table to show accordingly
	
*/}}{{if ne (.GetTagValue "Time") ""}}{{.GetTable "| c c | l r | l r | c | c | c | c |" (.GetSlice "Date" "Time" "White" "WhiteElo" "Black" "BlackElo" "ECO" "TimeControl" "Moves" "Result") }}
{{else}}{{.GetTable "| c | l r | l r | c | c | c | c |" (.GetSlice "Date" "White" "WhiteElo" "Black" "BlackElo" "ECO" "TimeControl" "Moves" "Result") }}
{{end}} # Games found: {{.Len}}
{{""}}
