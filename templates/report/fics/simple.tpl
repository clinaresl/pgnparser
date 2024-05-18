{{/*

	This template provides the basic functionality to just show
	each game of a collection in a different page.

	Every page starts with a nice header showing some
	administrative information about the game including players'
	names, their ELO, the winner, ECO, ... 

*/}}

\documentclass[oneside,svgnames]{report}

\usepackage[a4paper, total={7.5in, 10in}]{geometry}

\usepackage[utf8]{inputenc}
\usepackage[english]{babel}

\usepackage{xcolor}

\usepackage{marvosym}
\usepackage{FiraSans}

\usepackage{fancyhdr}
\pagestyle{fancy}

\renewcommand{\headrulewidth}{0pt}
\renewcommand{\footrulewidth}{0.4pt}

\lfoot{\texttt{\href{https://lichess.org}{\includegraphics[width=0.2in]{lichess.png}\texttt{lichess.org}}}}\cfoot{\thepage}\rfoot{\href{https://github.com/clinaresl/pgnparser}{\texttt{pgnparser}}}

\usepackage{xskak}

\usepackage{hyperref}
\hypersetup{
    colorlinks=true,
    urlcolor=RoyalBlue,
    pdfpagemode=FullScreen,
}
    
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
{{.GetField ("Site")}}}{\Mundus~}{{.GetField ("Event")}} ({{.GetField ("TimeControl")}})}
\end{center}

\hrule
\vspace{0.1cm}
\noindent
\raisebox{-5pt}{\WhiteKnightOnWhite} \textcolor{Olive}{%
{{.GetField ("White")}} ({{.GetField ("WhiteElo")}})} \hfill \textcolor{Sienna}{%
{{.GetField ("Date")}}}\\
\raisebox{-5pt}{\BlackKnightOnWhite} \textcolor{Olive}{%
{{.GetField ("Black")}} ({{.GetField ("BlackElo")}})} \hfill \textcolor{IndianRed}{%
ECO: {{.GetField ("ECO")}}}
\hrule

\vspace{0.5cm}

{{/* -------------------------------- Moves ------------------------------ */}}

\newchessgame
{{.GetLaTeXMovesWithComments}}\hfill \textbf{ {{.GetField ("Result")}}}\\

{{/* --------------------------- Final position -------------------------- */}}

\begin{center}
  \chessboard[print,showmover=true]
\end{center}
\noindent

\newpage

{{end}}

\end{document}
