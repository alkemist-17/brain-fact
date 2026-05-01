+++       Cell 0 = 3 (outer counter)
[         Outer loop (runs 3 times):
  >+++    Cell 1 = 3 (inner counter; reset each time)
  [       Inner loop (runs 3 times):
    >+    Add 1 to cell 2
    <-    Subtract 1 from cell 1
  ]       Inner loop ends when cell 1 = 0
  <-      Subtract 1 from cell 0
]         Outer loop ends when cell 0 = 0