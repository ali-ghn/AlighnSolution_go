#!/bin/bash
tr -dc A-Za-z0-9_ < /dev/urandom | head -c "$1" | xargs