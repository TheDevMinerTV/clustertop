#!/bin/sh

chown -R app:app /static

su app -c "/clustertop --listen-addr :80 $*"