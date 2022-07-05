#!/usr/bin/env bash

time STDIN=testsentences10 ./pcfg_tool -a=grammar.outside parse material/large/grammar.rules material/large/grammar.lexicon
time STDIN=testsentences10 ./pcfg_tool -r=10000 -a=grammar.outside parse material/large/grammar.rules material/large/grammar.lexicon
time STDIN=testsentences10 ./pcfg_tool -t=0.000001 -a=grammar.outside parse material/large/grammar.rules material/large/grammar.lexicon