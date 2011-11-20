include $(GOROOT)/src/Make.inc

GOFMT=gofmt -w -tabindent -tabwidth=8

TARG=wp
GOFILES=\
	wplib.go\
	main.go

TARG=wp_s
GOFILES=\
	wplib_s.go\
	main_s.go

include $(GOROOT)/src/Make.cmd

format:
	${GOFMT} ${GOFILES}
