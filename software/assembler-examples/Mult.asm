// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)

@i  // i=0
M=0

@R2  // R2=0
M=0

(LOOP)

  @i  // if (i >= R0) GOTO LOOP
  D=M
  @R0
  D=M-D
  @END
  D;JLE

  @R1 // R2+=R1 
  D=M
  @R2
  M=M+D

  @i  // i++
  M=M+1

  @LOOP // GOTO LOOP
  0;JMP

(END)
  @END
  0;JMP

