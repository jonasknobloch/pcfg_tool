#!/usr/bin/env bash

make -f ./Makefile -C . build
source ./third_party/disco-dop/venv/bin/activate
./pcfg_tool parse ./material/1000/grammar.rules ./material/1000/grammar.lexicon < ./material/1000/sentences-100 > pcfg_tool.mrg
discodop eval --fmt=bracket ./material/1000/training_b-100.mrg ./pcfg_tool.mrg | grep labeled.*: