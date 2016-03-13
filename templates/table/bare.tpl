{{/* this is a bare template which provides almost no information,
other than the name of both players, their elo score and the final
result of each game.

In other words, it provides ONLY bare information about each game */}}

{{.GetTable "|lr|lr|c|" (.GetSlice "White" "WhiteElo" "Black" "BlackElo" "Result") }}

# Games found: {{.Len}}
{{""}}
