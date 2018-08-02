# junit-exception-test-snippet
At Junit5, attribute "expected" at @Test is removed. this is snippet code generator.



## Build

go build main.go lexer.go


## How to use

### Example code modification

Output sample modification to stdout

```
./main -example
```

### code modification

```
./main -input ./src/test/com/github/yokotaso/SampleTest.java -output ./src/test/com/github/yokotaso/SampleTest.java
```


