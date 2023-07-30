### GLOX
Golang implementation of the Lox programming language, as described in https://craftinginterpreters.com/

This code attempts to mirror the Java code in the Crafting Interpreters book, which means that it is probably
non-idiomatic to some degree. 


### Grammer (so far)

Grammar syntax ( https://craftinginterpreters.com/representing-code.html#enhancing-our-notation ) :

```

* : zero or more repetitions
? : zero or one
| : or 

```


```
program        → declaration* EOF ;

declaration    → varDecl
               | statement ;

varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;

statement      → exprStmt
               | ifStmt
               | printStmt
               | whileStmt
               | block ;

whileStmt      → "while" "(" expression ")" statement ;

ifStmt         → "if" "(" expression ")" statement
               ( "else" statement )? ;
               
block          → "{" declaration* "}" ;
               
exprStmt       → expression ";" ;
printStmt      → "print" expression ";" ;

expression     → assignment ;
assignment     → IDENTIFIER "=" assignment
               | logic_or ;
logic_or       → logic_and ( "or" logic_and )* ;
logic_and      → equality ( "and" equality )* ;
               
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;

primary        → "true" | "false" | "nil"
               | NUMBER | STRING
               | "(" expression ")"
               | IDENTIFIER ;
```

### How we parse the grammar (see chapter 6.2)

```
Grammar notation               Code representation
Terminal                       Code to match and consume a token
Nonterminal                    Call to that rule’s function
|                              if or switch statement
* or +                         while or for loop
?                              if statement
```

### Current chapter
9 - Control flow

### Helpful links
Precedence and associativity: https://craftinginterpreters.com/parsing-expressions.html#ambiguity-and-the-parsing-game
