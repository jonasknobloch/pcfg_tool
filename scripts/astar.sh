#!/usr/bin/env bash

./pcfg_tool outside material/1000/grammar.rules material/1000/grammar.lexicon 1000

echo "pcfg_tool parse"

time STDIN=material/1000/sentences-100 ./pcfg_tool parse material/1000/grammar.rules material/1000/grammar.lexicon 1> /dev/null

echo -e "\npcfg_tool parse --astar=1000.outside"

time STDIN=material/1000/sentences-100 ./pcfg_tool parse --astar=1000.outside material/1000/grammar.rules material/1000/grammar.lexicon 1> /dev/null