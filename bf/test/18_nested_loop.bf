TESTS Simple Nested loops
EXPECTED Q (ASCII 91)

+++++++++       Cell 0 = 9 (outer counter)
[               Outer loop (runs 9 times):
  >+++++++++    Cell 1 = 9 (inner counter; reset each time)
  [             Inner loop (runs 9 times):
    >+          Add 1 to cell 2
    <-          Subtract 1 from cell 1
  ]             Inner loop ends when cell 1 = 0
  <-            Subtract 1 from cell 0
]               Outer loop ends when cell 0 = 0
>>.             Should print Q (ASCII 91)