" Vim syntax file
" Language:	Jack
" Maintainer:	Michal Zurawski (https://github.com/michalzurawski/NAND-to-Tetris)
" Last Change:	2017 Apr 14

if exists("b:current_syntax")
  finish
endif

syn keyword jackKeywords class constructor function method field static var this let do return
syn keyword jackConditionals if else while
syn keyword jackTypes int char boolean void true false null

hi def link jackKeywords        Keyword
hi def link jackConditionals    Conditional
hi def link jackTypes           Type

syn keyword jackTodos containted TODO FIXME

syn match jackComment "//.*$" contains=jackTodos
syn region jackComment start="/\*" end="\*/" contains=jackTodos fold

hi def link jackTodos           Todo
hi def link jackComment         Comment

syn region jackBlock start="{" end="}" transparent fold

syn match jackNumber "\d\+"
syn match jackString "\".*\""

hi def link jackNumber          Number
hi def link jackString          String

let b:current_syntax = "jack"
