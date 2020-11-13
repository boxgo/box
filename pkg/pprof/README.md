# Use the pprof tool

Then use the pprof tool to look at the heap profile:
```sh
go tool pprof http://localhost:9999/debug/pprof/heap
```

Or to look at a 30-second CPU profile:
```sh
go tool pprof http://localhost:9999/debug/pprof/profile
```

Or to look at the goroutine blocking profile, after calling runtime.SetBlockProfileRate in your program:
```sh
go tool pprof http://localhost:9999/debug/pprof/block
```

Or to collect a 5-second execution trace:
```sh
wget http://localhost:9999/debug/pprof/trace?seconds=5
```