include $(GOROOT)/src/Make.inc

GOFMT=gofmt -w -tabindent -tabwidth=8

TARG=wp_s
GOFILES=\
	main_s.go\
	wplib_s.go

GC=6g -I../kmp/_obj
LD=6l -L../kmp/_obj

include $(GOROOT)/src/Make.pkg

fmt:
	${GOFMT} ${GOFILES}
