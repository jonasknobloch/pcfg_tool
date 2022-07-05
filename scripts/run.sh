#!/usr/bin/env bash

time STDIN=material/small/sentences ./pcfg_tool parse material/small/grammar.rules material/small/grammar.lexicon
time STDIN=material/large/testsentences ./pcfg_tool parse material/large/grammar.rules material/large/grammar.lexicon