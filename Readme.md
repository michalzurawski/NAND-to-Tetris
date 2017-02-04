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

## Software
