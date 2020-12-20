#!/bin/bash

grep '"name"' | awk '{ print $2 }' | tr -d '",'|sort | uniq -c |sort - | awk '$1 != 1 { print $2 }'
