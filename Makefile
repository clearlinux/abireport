PROJECT_ROOT := src/
VERSION = 1.0.3

.DEFAULT_GOAL := all

# The resulting binaries map to the subproject names
BINARIES = \
	abireport

include Makefile.gobuild

_PKGS = \
	libabi \
	explode \
	abireport \
	abireport/cmd

# We want to add compliance for all built binaries
_CHECK_COMPLIANCE = $(addsuffix .compliant,$(_PKGS))

# Build all binaries as static binary
BINS = $(addsuffix .statbin,$(BINARIES))

# Ensure our own code is compliant..
compliant: $(_CHECK_COMPLIANCE)
install: $(BINS)
	test -d $(DESTDIR)/usr/bin || install -D -d -m 00755 $(DESTDIR)/usr/bin; \
	install -m 00755 bin/* $(DESTDIR)/usr/bin/.; \
	test -d $(DESTDIR)/usr/share/man/man1 || install -D -d -m 00755 $(DESTDIR)/usr/share/man/man1; \
	install -m 00644 man/*.1 $(DESTDIR)/usr/share/man/man1/.; \

# Credit to swupd developers: https://github.com/clearlinux/swupd-client
MANPAGES = \
	man/abireport.1

gen_docs:
	for MANPAGE in $(MANPAGES); do \
		ronn --roff < $${MANPAGE}.md > $${MANPAGE}; \
		ronn --html < $${MANPAGE}.md > $${MANPAGE}.html; \
	done

ensure_modules:
	@ ( \
		git submodule init; \
		git submodule update; \
	);

# See: https://github.com/meitar/git-archive-all.sh/blob/master/git-archive-all.sh
release: ensure_modules
	git-archive-all.sh --format tar.gz --prefix abireport-$(VERSION)/ --verbose -t HEAD abireport-$(VERSION).tar.gz

all: $(BINS)
