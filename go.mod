module github.com/damienjacinto/goproxy

replace github.com/damienjacinto/internal/utils => ./internal/utils

go 1.18

require (
	github.com/damienjacinto/internal/utils v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
	github.com/rs/zerolog v1.27.0
	github.com/steveyen/gkvlite v0.0.0-20141117050110-5b47ed6d7458
)

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)
