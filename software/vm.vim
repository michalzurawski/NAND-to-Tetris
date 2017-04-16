" Vim syntax file
" Language:	Jack Virtual Machine
" Maintainer:	Michal Zurawski (https://github.com/michalzurawski/NAND-to-Tetris)
" Last Change:	2017 Apr 15

if exists("b:current_syntax")
  finish
endif

syn region jackVMFunctionBlock start="function" end="return" transparent fold

syn keyword jackVMKeywords push pop call
syn keyword jackVMConditionals label goto
syn match jackVMConditionals "if-goto"
syn keyword jackVMTypes static constant local temp this that argument pointer
syn keyword jackVMSymbols and add or sub neg eq gt lt not

hi def link jackVMKeywords        Keyword
hi def link jackVMConditionals    Conditional
hi def link jackVMTypes           Type
hi def link jackVMSymbols         Type

syn match jackVMComment "//.*$"

hi def link jackVMComment         Comment

syn match jackVMNumber "\d\+"

hi def link jackVMNumber          Number

let b:current_syntax = "vm"
