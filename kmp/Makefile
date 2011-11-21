include $(GOROOT)/src/Make.inc

GOFMT=gofmt -w -tabindent -tabwidth=8

TARG=kmp
GOFILES=\
	kmp.go\

include $(GOROOT)/src/Make.pkg

fmt:
	${GOFMT} ${GOFILES}
