# pipe

pipe: use in go generate statements to pipe commands
* commands are executed left to right
* stdout from preceding is piped into stdin of next
* commands are separated by ::
* format is: pipe cmd0 arg0 arg1 ... :: cmd1 arg0 arg1 ... :: ...
* use pipe -v (verbose) to print output from each command: pipe -v cmd0

```
go get github.com/exyzzy/metaapi
go install $GOPATH/src/github.com/exyzzy/metaapi
```

Then, for example in todo.go:

```
//go:generate  pipe -v metaapi -pipe=true -sql=todo.sql  :: govueintro -pipe=true -txt=api.txt
//go:generate  pipe -v metaapi -pipe=true -sql=todo.sql  :: govueintro -pipe=true -txt=api_test.txt
```