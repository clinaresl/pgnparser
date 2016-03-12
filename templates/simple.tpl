\documentclass{article}

\usepackage[utf8]{inputenc}
\usepackage[english]{babel}
\usepackage{mathpazo}
\usepackage{skak}

\def\hrulefill{\leavevmode\leaders\hrule height 10pt\hfill\kern\z@}

\begin{document}

{{range .GetGames}}

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

\newgame

{{.GetLaTeXMovesWithComments}}\hfill \textbf{ {{.GetTagValue ("Result")}}}

\begin{center}
  \showboard
\end{center}

\clearpage

{{end}}

\end{document}
