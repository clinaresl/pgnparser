{{/* this is a simple template which provides additional (purely
administrative) information in addition to the name of both players,
their elo score and the final result.

In other words, it provides ONLY basic information about each game
*/}}

{{.GetTable "|cc|lr|lr|c|c|c|c|" (.GetSlice "Date" "Time" "White" "WhiteElo" "Black" "BlackElo" "ECO" "TimeControl" "Moves" "Result") }}

# Games found: {{.Len}}
{{""}}
