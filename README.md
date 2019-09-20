- commands to use
```bash
go test -bench .
go test -bench . -benchmem
go test -v
go test -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 main_test.go
go tool pprof mem.out
go tool pprof cpu.out
```
