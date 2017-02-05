# NAND to Tetris

## Description
This project, based on the course [Nand to Tetris](http://nand2tetris.org), shows an example design of a fully operational, multipurpose 16-bit computer, constructed using only [NAND](https://en.wikipedia.org/wiki/NAND_gate) logic gates and basic [flip-flops](https://en.wikipedia.org/wiki/Flip-flop_(electronics)). This project includes own assembler, compiler, Virtual Machine and basic Operating System. Design of this computer assumes [memory-mapped I/O](https://en.wikipedia.org/wiki/Memory-mapped_I/O) is given for the input (keyboard) and output (monitor). Implementation and design of each component is written in bottom-up manner, starting from creation of basic logic gates through mutexes, CPU, and ends on Tetris game.

## Table of Contents
1. [Hardware](#hardware)
  1. [Basic logic gates](#basic-logic-gates)
    1. [NAND](#nand)
    2. [NOT](#not)
    3. [AND](#and)
    4. [OR](#or)
    5. [XOR](#xor)
  2. [Multiplexers](#multiplexers)
    1. [Multiplexer](#Multiplexer)
    2. [Demultiplexer](#Demultiplexer)
  3. [16 bits versions of gates](#16-bits-versions-of-gates)
    1. [16 bits NOT](#16-bits-not)
    2. [16 bits AND](#16-bits-and)
    3. [16 bits OR](#16-bits-or)
    4. [OR 8 Way](#or-8-way)
    5. [16 bits Multiplexer](#16-bits-multiplexer)
    6. [16 bits Multiplexer 4 Way](#16-bits-multiplexer-4-way)
    7. [16 bits Multiplexer 8 Way](#16-bits-multiplexer-8-way)
    8. [Demultiplexer 4 Way](#demultiplexer-4-way)
    9. [Demultiplexer 8 Way](#demultiplexer-8-way)
  4. [Arithmetic Logic Union](#arithmetic-logic-union)
    1. [Half Adder](#half-adder)
    2. [Full Adder](#full-adder)
    3. [Adder](#adder)
    4. [Incrementer](#incrementer)
    5. [ALU](#alu)
2. [Software](#software)

## Hardware
Each piece of hardware is constructed either from basic NAND, Flip-Flop or using already designed elements.
All elements have been described using [Hardware Description Language](https://en.wikipedia.org/wiki/Hardware_description_language) and can be find in the hardware directory.
This files can be tested using Hardware Simulator in the tool directory.
Additionally there are also presented below as a drawings - each of them contain symbol of that element (on the right/below) and design.

### Basic logic gates
1. [NAND](#nand)
2. [NOT](#not)
3. [AND](#and)
4. [OR](#or)
5. [XOR](#xor)

#### NAND
NAND gates will be used to construct other logic gates.

| In0 | In1 | Out |
| --- | --- | --- |
| 0   | 0   | 1   |
| 0   | 1   | 1   |
| 1   | 0   | 1   |
| 1   | 1   | 0   |

![NAND](hardware/basic-logic-gates/Nand.png "NAND")

#### NOT

| In  | Out |
| --- | --- |
| 0   | 1   |
| 1   | 0   |

![NOT](hardware/basic-logic-gates/Not.png  "NOT")

#### AND

| In0 | In1 | Out |
| --- | --- | --- |
| 0   | 0   | 0   |
| 0   | 1   | 0   |
| 1   | 0   | 0   |
| 1   | 1   | 1   |

![AND](hardware/basic-logic-gates/And.png  "AND")

#### OR

| In0 | In1 | Out |
| --- | --- | --- |
| 0   | 0   | 0   |
| 0   | 1   | 1   |
| 1   | 0   | 1   |
| 1   | 1   | 1   |

![OR](hardware/basic-logic-gates/Or.png  "OR")

#### XOR

| In0 | In1 | Out |
| --- | --- | --- |
| 0   | 0   | 0   |
| 0   | 1   | 1   |
| 1   | 0   | 1   |
| 1   | 1   | 0   |

![XOR](hardware/basic-logic-gates/Xor.png  "XOR")

### Multiplexers

1. [Multiplexer](#multiplexer)
2. [Demultiplexer](#demultiplexer)

#### Multiplexer

| sel | In0 | In1 | Out |
|-----| --- | --- | --- |
| 0   | 0   | 0   | 0   |
| 0   | 0   | 1   | 0   |
| 0   | 1   | 0   | 1   |
| 0   | 1   | 1   | 1   |
| 1   | 0   | 0   | 0   |
| 1   | 0   | 1   | 1   |
| 1   | 1   | 0   | 0   |
| 1   | 1   | 1   | 1   |

Or in the other words:

| sel | Out |
|-----|-----|
| 0   | In0 |
| 1   | In1 |

![Multiplexer](hardware/multiplexers/Mux.png  "Multiplexer")

#### Demultiplexer

| sel | In | Out0 | Out1 |
|-----|----|------|------|
| 0   | 0  | 0    | 0    |
| 0   | 1  | 1    | 0    |
| 1   | 0  | 0    | 0    |
| 1   | 1  | 0    | 1    |

Or in the other words:

| sel | Out0 | Out1 |
| --- | ---- | ---- |
| 0   | In   | 0    |
| 1   | 0    | In   |

![Demultiplexer](hardware/multiplexers/DMux.png  "Demultiplexer")

### 16 bits versions of gates

16 bits version of base logic gates and mutexes.
Notation `[0..15]` indicates 16 bits number where 0 is the least significant digit.

1. [16 bits NOT](#16-bits-not)
2. [16 bits AND](#16-bits-and)
3. [16 bits OR](#16-bits-or)
4. [OR 8 Way](#or-8-way)
5. [16 bits Multiplexer](#16-bits-multiplexer)
6. [16 bits Multiplexer 4 Way](#16-bits-multiplexer-4-way)
7. [16 bits Multiplexer 8 Way](#16-bits-multiplexer-8-way)
8. [Demultiplexer 4 Way](#demultiplexer-4-way)
9. [Demultiplexer 8 Way](#demultiplexer-8-way)

#### 16 bits NOT

Negates each bit independently.

![NOT16](hardware/gates-16-bits/Not16.png  "NOT16")

#### 16 bits AND

Bitwise `and` on each pair of bits.

![AND16](hardware/gates-16-bits/And16.png  "AND16")

#### 16 bits OR

Bitwise `or` on each pair of bits.

![OR16](hardware/gates-16-bits/Or16.png  "OR16")

#### OR 8 Way

Returns 0 if all 8 bits for both of numbers are set to 0, 1 otherwise.

![OR8Way](hardware/gates-16-bits/Or8Way.png  "OR8Way")

#### 16 bits Multiplexer

| sel | Out |
| --- | --- |
| 0   | In0 |
| 1   | In1 |

![Mux16](hardware/gates-16-bits/Mux16.png  "Mux16")

#### 16 bits Multiplexer 4 Way

| sel | Out |
| --- | --- |
| 00  | In0 |
| 01  | In1 |
| 10  | In2 |
| 11  | In3 |

![Mux4Way16](hardware/gates-16-bits/Mux4Way16.png  "Mux4Way16")

#### 16 bits Multiplexer 8 Way

| sel | Out |
| --- | --- |
| 000 | In0 |
| 001 | In1 |
| 010 | In2 |
| 011 | In3 |
| 100 | In4 |
| 101 | In5 |
| 110 | In6 |
| 111 | In7 |

![Mux8Way16](hardware/gates-16-bits/Mux8Way16.png  "Mux8Way16")

#### Demultiplexer 4 Way

| sel | Out0 | Out1 | Out2 | Out3 |
| --- | ---- | ---- | ---- | ---- |
| 00  | In   | 0    | 0    | 0    |
| 01  | 0    | In   | 0    | 0    |
| 10  | 0    | 0    | In   | 0    |
| 11  | 0    | 0    | 0    | In   |

![DMux4Way](hardware/gates-16-bits/DMux4Way.png  "DMux4Way")

#### Demultiplexer 8 Way

| sel | Out0 | Out1 | Out2 | Out3 | Out4 | Out5 | Out6 | Out7 |
| --- | ---  | ---- | ---- | ---- | ---- | ---- | ---- | ---- |
| 000 | In   | 0    | 0    | 0    | 0    | 0    | 0    | 0    |
| 001 | 0    | In   | 0    | 0    | 0    | 0    | 0    | 0    |
| 010 | 0    | 0    | In   | 0    | 0    | 0    | 0    | 0    |
| 011 | 0    | 0    | 0    | In   | 0    | 0    | 0    | 0    |
| 100 | 0    | 0    | 0    | 0    | In   | 0    | 0    | 0    |
| 101 | 0    | 0    | 0    | 0    | 0    | In   | 0    | 0    |
| 110 | 0    | 0    | 0    | 0    | 0    | 0    | In   | 0    |
| 111 | 0    | 0    | 0    | 0    | 0    | 0    | 0    | In   |

![DMux8Way](hardware/gates-16-bits/DMux8Way.png  "DMux8Way")

### Arithmetic Logic Union

1. [Half Adder](#half-adder)
2. [Full Adder](#full-adder)
3. [Adder](#adder)
4. [Incrementer](#incrementer)
5. [ALU](#alu)

#### Half Adder

Used for creating [Full Adder](#full-adder).

| In0 | In1 | sum | carry |
| --- | --- | --- | ----- |
| 0   | 0   | 0   | 0     |
| 0   | 1   | 1   | 0     |
| 1   | 0   | 1   | 0     |
| 1   | 1   | 0   | 1     |

![HalfAdder](hardware/alu/HalfAdder.png  "HalfAdder")

#### Full Adder

Used for creating [Adder](#adder).

| In0 | In1 | In2 | sum | carry |
| --- | --- | --- | --- | ----- |
| 0   | 0   | 0   | 0   | 0     |
| 0   | 0   | 1   | 1   | 0     |
| 0   | 1   | 0   | 1   | 0     |
| 0   | 1   | 1   | 0   | 1     |
| 1   | 0   | 0   | 1   | 0     |
| 1   | 0   | 1   | 0   | 1     |
| 1   | 1   | 0   | 0   | 1     |
| 1   | 1   | 1   | 1   | 1     |

![FullAdder](hardware/alu/FullAdder.png  "FullAdder")

#### Adder

Adds two 16-bit values.
The most significant carry bit is ignored.

![Add16](hardware/alu/Add16.png  "Add16")

#### Incrementer

16 bits Incrementer. Increases the input value by 1.

![Inc16](hardware/alu/Inc16.png  "Inc16")

#### ALU

16 bits Arithmetic Logic Union.

Computes one of the following functions:

`x+y, x-y, y-x, 0, 1, -1, x, y, -x, -y, !x, !y, x+1, y+1, x-1, y-1, x&y, x|y`

on two 16-bit inputs, according to 6 input bits denoted `zx, nx, zy, ny, f, no`.
In addition, the ALU computes two 1-bit outputs:

if the ALU output == 0, `zr` is set to 1; otherwise `zr` is set to 0;

if the ALU output < 0, `ng` is set to 1; otherwise `ng` is set to 0.

`zx` - zero `x`

`nx` - negate `x`

`zy` - zero `y`

`zy` - negate `y`

`f == 0` - computes `x & y`

`f == 1` - computes `x + y`

`no` - negate output

| zx  | nx  | zy  | ny  | f   | no  | Out  |
| --- | --- | --- | --- | --- | --- | ---- |
| 1   | 0   | 1   | 0   | 0   | 0   | 0    |
| 1   | 1   | 1   | 1   | 1   | 1   | 1    |
| 1   | 1   | 1   | 0   | 1   | 0   | -1   |
| 0   | 0   | 1   | 1   | 0   | 0   | x    |
| 1   | 1   | 0   | 0   | 0   | 0   | y    |
| 0   | 0   | 1   | 1   | 0   | 1   | !x   |
| 1   | 1   | 0   | 0   | 0   | 1   | !y   |
| 0   | 0   | 1   | 1   | 1   | 1   | -x   |
| 1   | 1   | 0   | 0   | 1   | 1   | -y   |
| 0   | 1   | 1   | 1   | 1   | 1   | x+1  |
| 1   | 1   | 0   | 1   | 1   | 1   | y+1  |
| 0   | 0   | 1   | 1   | 1   | 0   | x-1  |
| 1   | 1   | 0   | 0   | 1   | 0   | y-1  |
| 0   | 0   | 0   | 0   | 1   | 0   | x+y  |
| 0   | 1   | 0   | 0   | 1   | 1   | x-y  |
| 0   | 0   | 0   | 1   | 1   | 1   | y-x  |
| 0   | 0   | 0   | 0   | 0   | 0   | x&y  |
| 0   | 1   | 0   | 1   | 0   | 1   | x\|y |

![ALU](hardware/alu/ALU.png  "ALU")

## Software
