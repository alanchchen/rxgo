# rxgo

Rx for Go

# Examples

Two observables that emit odd and even numbers.
```bash
$ go run examples/even-odd/main.go
evenSource 0
oddSource 1
evenSource 2
oddSource 3
evenSource 4
oddSource 5
evenSource 6
oddSource 7
evenSource 8
```

Map source to current time
```bash
$ go run examples/map/main.go
2018-12-13T14:26:44
2018-12-13T14:26:45
2018-12-13T14:26:46
2018-12-13T14:26:47
2018-12-13T14:26:48
2018-12-13T14:26:49
2018-12-13T14:26:50
2018-12-13T14:26:51
2018-12-13T14:26:52
```

`sub2` subscribes after `sub1` by 1 second
```bash
$ go run examples/subject/main.go
sub1 0
sub1 1
sub2 1
sub1 2
sub2 2
sub1 3
sub2 3
sub1 4
sub2 4
sub2 5
sub1 5
sub2 6
sub1 6
sub1 7
sub2 7
sub2 8
sub1 8
sub2 9
sub1 9
```
