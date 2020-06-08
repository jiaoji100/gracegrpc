module server-demo

go 1.13

require (
	github.com/jiaoji100/gracegrpc v0.0.0-20200608065336-77fb6d0465c2
	google.golang.org/grpc v1.30.0-dev.1
)

replace github.com/jiaoji100/gracegrpc => ../..
