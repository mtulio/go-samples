

WORDS_ALPHA_URL := https://github.com/dwyl/english-words/raw/master/words_alpha.txt
WORDS_ALPHA_PATH := ./hack/inputDict.txt
DOWNLOAD_FORCE ?= 0


# ##################
# Makefile functions
define do_download
	wget $(WORDS_ALPHA_URL) -O $(WORDS_ALPHA_PATH);
endef

.PHONY: download clean build
download:
	@if [[ $(DOWNLOAD_FORCE) -eq 1  ]]; then \
		echo "#> INPUT file forced to be updated."; \
		$(call do_download) \
		exit 0; \
	fi
	@if [[ ! -f $(WORDS_ALPHA_PATH) ]]; then \
		echo "#> INPUT file not found, downloading."; \
		$(call do_download) \
		exit 0; \
	fi

clean:
	@rm -f $(WORDS_ALPHA_PATH)


build:
	@go build

run:
	@./word-filter