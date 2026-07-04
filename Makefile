# canonic.c — build, prove, and gate the kernel.
#   make test   build + run the e2e selftest (VALID · tamper caught · 255 · 237)
#   make gate   the full hygiene gate: projection purity, then the selftest
UNAME := $(shell uname -s)
CFLAGS ?= -Os -Wall -Wextra
ifneq ($(UNAME),Darwin)
CFLAGS += -Icompat -Wno-deprecated-declarations
LDLIBS += -lcrypto
endif

all: selftest

canonic.o: src/canonic.c
	$(CC) $(CFLAGS) -c $< -o $@

selftest: test/selftest.c canonic.o
	$(CC) $(CFLAGS) -Isrc $^ -o $@ $(LDLIBS)

test: selftest
	./selftest

gate:
	grep -v '^#' PROJECTION | shasum -a 256 -c -
	$(MAKE) test

clean:
	rm -f canonic.o selftest

.PHONY: all test gate clean
