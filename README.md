# brain-fact
A [brainf**k](http://brainfuck.org/) interpreter written in [Go](https://go.dev/).


Import the library in your project with: 

```

go get github.com/alkemist-17/brain-fact

```


If you want to build the cli and run file scripts then:

```

cd bf
go build
./bf <script-name>.bf

```


Or run the REPL with

```

cd bf
go build
./bf

```

## Project Milestones.
| Component      	| Status 	|
|----------------	|--------	|
| Lexer          	| ✅      	|
| Compiler       	| ✅      	|
| Interpreter    	| ✅      	|
| Code Examples  	| ✅      	|
| Configurations 	| ✅      	|
| CLI tool       	| ✅      	|
| REPL       	    | ✅      	|
| Documentation  	| ✅     	|
| Release       	| ✅      	|
| Tests             | ✅         |
| Config Addr Size  | ✅         |
| Compiled Files    | ✅         |


## Implementation Details.
According to [The epistle to the Implementors](http://brainfuck.org/epistle.html), there are some specific topics to consider before doing an implementation of the language. In our implementation we follow some of them. Here is the list:


- When the value of a cell is 10, it means end-of-line.
- When found an EOF, it does not affect anything.
- The supported character set of operators is { '+', ',', '-', '.', '<', '>', '[', ']' }. Any other character is considered a comment.
- The input must be a valid utf-8 encoded file.
- The default size of the slice is 32767 cells.
- The max bytecode size is configurable to 65535 (default) or 4294967295 bytes.
- The pointer starts at position 0 of the tape.
- Output is not necessarily done with a monospace font. It will depend on the local context.
- I/O operations are done during execution, so all programs are interactive.
- The input operation has a default, yet configurable prompt.
- The compiler matches [] before execution, and it will report an error with unbalanced [].
- The interpreter have some configurable features like: prompt, and slice size.
- The interpreter does not work with negative pointers. (is it too much? 🤔)
- The interpreter will wrap around individual cells.
- May peace be with all implementors.


The interpreter is intended to be used as a package to embed in any Go project. It exports two configurable global parameters (tape size and promtp), and the function `Run(code string)`. The interpreter is thread safe, and does not produce any data races. 


Nobody knows when is needed a concurrent safe [brainf**k](http://brainfuck.org/) interpreter written in [Go](https://go.dev/).


Check out [this](https://brainfuck.net) interpreter.
It may work as an interpreter validator for implementators. 


Contributions are welcome!


Happy brain-fact coding! I❤️🧠




