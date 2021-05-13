#!/bin/bash

cargo build --release
cp target/release/libwasmvm.a api
# FIXME: re-enable stripped so when we approach a production release, symbols are nice for debugging
# strip api/libwasmvm.so
