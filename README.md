# TISC-100
TISC-100 is a virtual machine for the TIS-100 and an interpreter of its assembly language.

## Disclaimer
The game TIS-100 is owned entirely by Zachtronics Industries. The archetecture and assembly
language used in this project was originally created by them. This silly command line tool
does not replicate the experience of the game. There is no puzzle, no narrartive, and no
GUI.

TIS-100 is an awesome and incredibly creative game. If you haven't played it yet, you should.

## Purpose
The goal of this project is to create a virtual environment that fully implements the TIS-100
specification in order to experiment with the idea of writing real applications under the
limitations of this archetecture.

## Getting Started With the TIS-100
Information on the archetecture and assembly language is found in the TIS-100 manual, which can
be read online. The best way to learn how to program for the TIS-100, and by extension, this
command line tool, is to play the game.

## Creating a Project
A TISC-100 project is characterized by a set of `.tis` files and a single `machine.json`. The
`machine.json` is where the node structure of your program is defined. Projects in TISC-100
are named in the `machine.json` instead of in a special comment.

Nodes are defined as a two-dimensional array of strings. The letter "e" is an execution node
and "s" is a stack node.

Console input and output are defined with the side the input or output plugs into the node
array and the position of it on that side. If it plugs into the top or bottom, the position
value refers to its x position. If it plugs into the left or right, the position value
refers to its y position.

See the example project for a better idea of how to set up a TISC-100 project.
