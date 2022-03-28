# Makefile to build the command lines and tests in Seele project.
# This Makefile doesn't consider Windows Environment. If you use it in Windows, please be careful.

SHELL := /bin/bash

#BASEDIR = $(shell pwd)
BASEDIR = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))

$(info Current DIR: '${BASEDIR}')

.PHONY: default discover register

default: register

all: discover register

register:
	cd ${BASEDIR}/register && go run .

discover:
	cd ${BASEDIR}/discover && go run .
