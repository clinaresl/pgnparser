{{/*

	This template provides the basic functionality to just show
	each game of a collection in a different page.

	Every page starts with a nice header showing some
	administrative information about the game including players'
	names, their ELO, the winner, ECO, ... 

*/}}

\documentclass[svgnames]{report}

\usepackage[a4paper, total={6in, 8in}]{geometry}

\usepackage[utf8]{inputenc}
\usepackage[english]{babel}

\usepackage{xcolor}

\usepackage{FiraSans}

\usepackage{skak}
\usepackage{hyperref}

\def\hrulefill{\leavevmode\leaders\hrule height 10pt\hfill\kern\z@}

{{/* ----------------------------- Main Body ----------------------------- */}}

\begin{document}

\sffamily

{{/*
	For all games, just show the header and then the moves
	Finally, show a diagram with the final position of the game
*/}}

{{range .GetGames}} 

{{/* ------------------------------- Header ------------------------------ */}}

\begin{center}
  {\Large \href{%
{{.GetTagValue ("Site")}}}{ {{.GetTagValue ("Event")}} ({{.GetTagValue ("TimeControl")}}) } }
\end{center}

\hrule
\vspace{0.1cm}
\noindent
\raisebox{-5pt}{\WhiteKnightOnWhite} \textcolor{Olive}{%
{{.GetTagValue ("White")}} ({{.GetTagValue ("WhiteElo")}})} \hfill \textcolor{Sienna}{%
{{.GetTagValue ("Date")}}}\\
\raisebox{-5pt}{\BlackKnightOnWhite} \textcolor{Olive}{%
{{.GetTagValue ("Black")}} ({{.GetTagValue ("BlackElo")}})} \hfill \textcolor{IndianRed}{%
{{.GetTagValue ("Opening")}} ({{.GetTagValue ("ECO")}})}
\hrule

\vspace{0.5cm}

{{/* -------------------------------- Moves ------------------------------ */}}

\newgame
{{.GetLaTeXMovesWithComments}}\hfill \textbf{ {{.GetTagValue ("Result")}}}\\

{{/* --------------------------- Final position -------------------------- */}}

\begin{center}
  \showboard
\end{center}
\noindent
\hfill \textcolor{IndianRed}{Termination: {{.GetTagValue ("Termination")}}}

\clearpage

{{end}}

\end{document}
