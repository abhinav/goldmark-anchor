module go.abhg.dev/goldmark/anchor/demo

go 1.26.0

toolchain go1.27.1

replace go.abhg.dev/goldmark/anchor => ../

require (
	github.com/yuin/goldmark/v2 v2.1.5
	go.abhg.dev/goldmark/anchor v0.2.0
)

require github.com/yuin/goldmark v1.8.6 // indirect
