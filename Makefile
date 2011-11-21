include $(GOROOT)/src/Make.inc

.PHONY: all install clean nuke fmt

all:
	gomake -C kmp
	gomake -C serial
	gomake -C parallel

install: all
	gomake -C kmp install
	gomake -C serial install
	gomake -C parallel install

clean:
	gomake -C kmp clean
	gomake -C serial clean
	gomake -C parallel clean

nuke:
	gomake -C kmp nuke
	gomake -C serial nuke
	gomake -C parallel nuke

fmt:
	gomake -C kmp fmt
	gomake -C serial fmt
	gomake -C parallel fmt
