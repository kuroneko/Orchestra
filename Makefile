#
# version of Orchestra
#
VERSION=0.3.0

#
# packaging revision.
#
REVISION=1

# remove at your own peril.
#
# This tells goinstall to work against the local directory as the
# build/source path, and not use the system directories.
#
GOPATH=$(PWD)/build-tree:$(PWD)
GOINSTALL_FLAGS=-dashboard=false -clean=true -u=false -make=false

export GOPATH

all: build

build:	build-tree
	goinstall $(GOINSTALL_FLAGS) conductor
	goinstall $(GOINSTALL_FLAGS) player
	goinstall $(GOINSTALL_FLAGS) submitjob
	goinstall $(GOINSTALL_FLAGS) getstatus

build-tree:
	mkdir -p build-tree/src
	mkdir -p build-tree/bin
	mkdir -p build-tree/pkg

clean:
	-$(RM) -r build-tree/pkg
	-$(RM) -r build-tree/bin

distclean:
	-$(RM) -r build-tree

### NOTE:  Make sure the checkouts have the correct tags in the lines below!
deps:	distclean build-tree
	mkdir -p build-tree/src/github.com/kuroneko && cd build-tree/src/github.com/kuroneko && git clone http://github.com/kuroneko/configureit.git && cd configureit && git checkout v0.1
	mkdir -p build-tree/src/goprotobuf.googlecode.com && cd build-tree/src/goprotobuf.googlecode.com && hg clone -r go.r60 http://goprotobuf.googlecode.com/hg

archive.deps:	deps
	tar czf ../orchestra-$(VERSION).build-tree.tar.gz -C build-tree --exclude .git --exclude .hg . 

archive.release:	archive.deps
	git archive --format=tar --prefix=orchestra-$(VERSION)/ v$(VERSION) | gzip -9c > ../orchestra-$(VERSION).tgz

.PHONY : debian debian.orig debian.debian debian.build-tree archive.deps archive.release archive.head

archive.head:
	git archive --format=tar --prefix=orchestra/ HEAD | gzip -9c > ../orchestra-HEAD.tgz

DEBIAN_VERSION=$(shell dpkg-parsechangelog | grep -e 'Version:' | awk '{ print $$2 }')
DEBIAN_SRC_VERSION=$(shell echo $(DEBIAN_VERSION) | cut -d- -f 1)

debian:	debian.orig debian.debian debian.build-tree clean
	cd .. && dpkg-source -b $(PWD)

debian.orig:
	git archive --format=tar --prefix=orchestra-$(DEBIAN_SRC_VERSION)/ v$(DEBIAN_SRC_VERSION) | gzip -9c > ../orchestra_$(DEBIAN_SRC_VERSION).orig.tar.gz

debian.debian:
	tar zcf ../orchestra_$(DEBIAN_VERSION).debian.tar.gz -C debian .

debian.build-tree:	deps
	tar zcf ../orchestra_$(DEBIAN_VERSION).orig-build-tree.tar.gz -C build-tree --exclude .git --exclude .hg .
