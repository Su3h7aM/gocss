module github.com/su3h7am/gocss

go 1.24.4

replace github.com/su3h7am/gocss/pkg/core => ./pkg/core

replace github.com/su3h7am/gocss/pkg/preset => ./pkg/preset

require github.com/fsnotify/fsnotify v1.9.0

require golang.org/x/sys v0.13.0 // indirect
