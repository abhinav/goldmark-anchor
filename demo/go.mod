module go.abhg.dev/goldmark/anchor/demo

go 1.25

toolchain go1.25.6

replace go.abhg.dev/goldmark/anchor => ../

require (
	github.com/yuin/goldmark/v2 v2.0.0
	go.abhg.dev/goldmark/anchor v0.2.0
)

require github.com/yuin/goldmark v1.7.16 // indirect
