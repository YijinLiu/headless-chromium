.PHONY: dev clean

VERSION?=57.0.2987.110

SHELL=/bin/bash
WORK_DIR:=$(shell readlink -f $(dir $(lastword $(MAKEFILE_LIST))))
PROJECT_ROOT:=$(shell readlink -f $(WORK_DIR)/../..)

dev: clean
	mkdir -p $(WORK_DIR)/docker_tmp_data
	for i in $(PROJECT_ROOT)/scripts $(PROJECT_ROOT)/cc ; do \
		cp -apf $$i $(WORK_DIR)/docker_tmp_data/ ; \
	done
	docker build -t hc/dev:$(VERSION) --build-arg=version=$(VERSION) .

clean:
	rm -rf docker_tmp_data
