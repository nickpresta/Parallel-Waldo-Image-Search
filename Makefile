include $(GOROOT)/src/Make.inc

GOFMT=gofmt -w -tabindent -tabwidth=8

TARG=wp
GOFILES=\
				wplib.go\
				wp.go

include $(GOROOT)/src/Make.cmd

format:
	${GOFMT} ${GOFILES}
