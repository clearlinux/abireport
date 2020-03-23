VERSION = 1.0.8

.DEFAULT_GOAL := all

clean:
	rm -f abireport abireport-*.tar.xz man/abireport.1 man/abireport.1.html
	rm -fr vendor

install: all
	test -d $(DESTDIR)/usr/bin || install -D -d -m 00755 $(DESTDIR)/usr/bin; \
	install -m 00755 abireport $(DESTDIR)/usr/bin/.; \
	test -d $(DESTDIR)/usr/share/man/man1 || install -D -d -m 00755 $(DESTDIR)/usr/share/man/man1; \
	install -m 00644 man/*.1 $(DESTDIR)/usr/share/man/man1/.; \

gen_docs: man/abireport.1.md
	pandoc -s -f markdown -t man man/abireport.1.md --output man/abireport.1
	pandoc -s -f markdown -t html man/abireport.1.md --output man/abireport.1.html

vendor:
	@go mod vendor

all: vendor
	(cd abi-report && go build --buildmode=pie -mod=vendor -o ../abireport)

dist: vendor gen_docs
	@rm -f abireport-$(VERSION).tar.xz
	@git tag -l | grep -q v$(VERSION) || (echo "tag v$(VERSION) not found"; exit 1)
	$(eval TMP := $(shell mktemp -d))
	@cp -r . $(TMP)/abireport-$(VERSION)
	@( \
		cd $(TMP)/abireport-$(VERSION); \
		git reset --hard v$(VERSION) &>/dev/null; \
		git clean -xf -e man/abireport.1 -e man/abireport.1.html; \
		rm -fr .git .gitignore; \
	);
	@tar -C $(TMP) -cf abireport-$(VERSION).tar abireport-$(VERSION)
	@xz abireport-$(VERSION).tar
	@rm -fr $(TMP)
