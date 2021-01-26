# Developers Guide

This repository contains testing and development utilities.

## Setup

clone into $GOPATH/src/github.com/paulcarlton-ww/testutils:

    mkdir -p $GOPATH/src/github.com/paulcarlton-ww
    cd $GOPATH/src/github.com/paulcarlton-ww
    git clone git@github.com:paulcarlton-ww/testutils.git
    cd testutils

This project requires the following software:

    golangci-lint --version = 1.30.0
    golang version >= 1.13.1

You can install these in the project bin directory using the 'setup.sh' script:

    bin/setup.sh

The setup.sh script can safely be run at any time. It installs the required software.

## Development

The Makefile in the project's top level directory will compile, build and test all components.

    make

To run the build and test in a docker container, type:

    make check

If changes are made to go source imports you may need to perform a go mod vendor, type:

    make gomod-update


