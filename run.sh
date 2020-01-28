go run goroutine.go  && go tool pprof -top goroutine_prof
go run goroutine.go  && go tool pprof -top goroutine_prof
go run goroutine.go  && go tool pprof -top goroutine_prof

go run serial.go && go tool pprof -top serial_prof
go run serial.go && go tool pprof -top serial_prof
go run serial.go && go tool pprof -top serial_prof
