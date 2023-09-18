#!/bin/bash
MIN=0
MAX=12345
COUNT=$COUNT
for i in `seq $COUNT`; do 
    rnd=$(( $RANDOM % ($MAX + 1 - $MIN) + $MIN ))
    printf "http://url.url/%d %d\n" $rnd $rnd
done >> file.txt