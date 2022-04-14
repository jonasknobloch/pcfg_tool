#!/usr/bin/env bash

validate() {
  diff <(sort < "$1" | tr -s "\t" " ") <(sort < "$2" | tr -s "\t" " ");
}

validate grammar.lexicon material/large/grammar.lexicon
validate grammar.rules material/large/grammar.rules
validate grammar.words material/large/grammar.words
