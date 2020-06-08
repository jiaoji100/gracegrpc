module server-demo

go 1.13

require (
	github.com/jiaoji100/gracegrpc/gracegrpc v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.30.0-dev.1
)

replace github.com/jiaoji100/gracegrpc/gracegrpc => ../..
