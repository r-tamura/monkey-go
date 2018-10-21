# monkey-go
「Writing an Interpreter in Go」写経

# Test
```
$ cd ${GOPATH}
$ mkdir monkey
$ cd monkey
$ git clone git@github.com:r-tamura/monkey-go.git .
```

# Run benchmark program

```
# build
$ go test -o fibonacci ./benchmark

# using 'eval'
$ fibonacci -engine=eval

# using 'vm'
$ fibonacci -engine=vm
```