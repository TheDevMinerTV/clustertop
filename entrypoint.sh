#!/bin/sh

su app -c "/clustertop --listen-addr :80 $*"