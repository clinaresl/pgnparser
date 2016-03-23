{{/*

	This template creates a document with sections where each
	section is a game. Sections are indexed in a table of contents
	at the beginning of the document, each identified with the
	players' names.

	Immediately after the table of contents, a summary table is
	shown. This is expected to be useful only for short documents
	---as long collections could actually make this list to go
	beyond the physical limits of the page.

	Every section contains then the same information portrayed in
	the simple template.

*/}}

\documentclass{article}

\usepackage[utf8]{inputenc}
\usepackage[english]{babel}
\usepackage{mathpazo}
\usepackage{skak}
\usepackage{booktabs}

\def\hrulefill{\leavevmode\leaders\hrule height 10pt\hfill\kern\z@}

{{/* ----------------------------- Main Body ----------------------------- */}}

\begin{document}

{{/* -------------------------- Table of Contents ------------------------ */}}

\cleardoublepage

\tableofcontents

\cleardoublepage

{{/* --------------------------- Summary table --------------------------- */}}

{{with $x := .GetTable "|lr|lr|c|" (.GetSlice "White" "WhiteElo" "Black" "BlackElo" "Result")}}
\vspace*{\fill}
{{printf "%v" $x.ToLaTeX}}
\vspace*{\fill}
{{end}}

\clearpage

{{/*
	For all games, just show the header and then the moves
	Finally, show a diagram with the final position of the game
*/}}

{{range .GetGames}}

{{/* ------------------------------- Section ----------------------------- */}}

\section{ {{.GetTagValue ("White")}} -- {{.GetTagValue ("Black")}} }

{{/* ------------------------------- Header ------------------------------ */}}

\begin{center}
  {\Large {{.GetTagValue ("Event")}} ({{.GetTagValue ("TimeControl")}})}
\end{center}

\hrule
\vspace{0.1cm}
\noindent
\raisebox{-5pt}{\WhiteKnightOnWhite} {{.GetTagValue ("White")}} ({{.GetTagValue ("WhiteElo")}}) \hfill {{.GetTagValue ("Date")}}\\
\raisebox{-5pt}{\BlackKnightOnWhite} {{.GetTagValue ("Black")}} ({{.GetTagValue ("BlackElo")}}) \hfill {{.GetTagValue ("ECO")}}
\hrule

\vspace{0.5cm}

{{/* -------------------------------- Moves ------------------------------ */}}

\newgame
{{.GetLaTeXMovesWithComments}}\hfill \textbf{ {{.GetTagValue ("Result")}}}

{{/* --------------------------- Final position -------------------------- */}}

\begin{center}
  \showboard
\end{center}

\clearpage

{{end}}

\end{document}
