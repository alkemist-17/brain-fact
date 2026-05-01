TESTS Parse brainfuck code and not brainfuck commnads well.
EXPECTED A (ASCII 65 = 13 x 5)

+++++++++++++  Set cell 0 to 13
[              While cell 0 is not 0:
  >+++++       Add 5 to cell 1
  <-           Go back; subtract 1 from cell 0
]              Cell 0 is now 0; cell 1 is 65
>.             Move to cell 1; print it (A)