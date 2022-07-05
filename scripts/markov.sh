#!/usr/bin/env bash

validate() {
  diff <(sort < "$1" | tr -s "\t" " ") <(sort < "$2" | tr -s "\t" " ");
}

STDIN=material/small/gold.mrg STDOUT=small_gold_b.mrg ./pcfg_tool binarize
STDIN=material/small/gold_b.mrg STDOUT=small_gold.mrg ./pcfg_tool debinarize

validate small_gold_b.mrg material/small/gold_b.mrg
validate small_gold.mrg material/small/gold.mrg