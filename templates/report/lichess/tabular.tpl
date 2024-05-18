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

\usepackage{booktabs}
\usepackage{marvosym}
\usepackage{FiraSans}

\usepackage{array}
\usepackage{longtable}
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
	Show an index of all games produced in this report along with
    hyperrefs that can be used to jump to any game
*/}}

{{.ShowIndex}}
\newpage

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
{{.GetField ("Opening")}} ({{.GetField ("ECO")}})}
\hrule

\vspace{0.5cm}

{{/* -------------------------------- Moves ------------------------------ */}}

\newchessgame
{{.GetLaTeXMovesWithCommentsTabular "4.2in" "3.0in" 8}}\hfill \textbf{ {{.GetField ("Result")}}}\\
{{.SetLabel}}

{{/* ------------------------------ Postface ----------------------------- */}}
\hfill \textcolor{IndianRed}{Termination: {{.GetField ("Termination")}}}

\newpage

{{end}}

\end{document}
