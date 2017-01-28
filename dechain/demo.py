#!/usr/bin/env python3

import sys
import argparse
import pprint

global args

def load_chains():
    return {
            # Normal links:
            'A': 'X',
            'B': 'Y',
            'C': 'Z',
            # Long chain: D->E->F->G->H->I-J
            'D': 'E',
            'E': 'F',
            'F': 'G',
            'G': 'H',
            'H': 'I',
            'I': 'J',
            # Another chain: O->N->M->L->K
            'L': 'K',
            'M': 'L',
            'N': 'M',
            'O': 'N',
            # Infinite loop:
            'O': 'O',
            }

def optimize_one_pass(chains):
    """The 1-pass optimizer.
    For any chains that are longer than 1 hop, reduces them by
    at least 1 hop (maybe more).
    """
    for i in sorted(chains.keys()):
        n =  chains[i]
        if n in chains:
            chains[i] = chains[n]

def print_chains(chains):
    """Print chains."""
    for j in sorted(chains.keys()):
        print(j, chains[j], chain_length(chains, j))

def count_chains(chains):
    """Returns the total length of all chains."""
    total = 0
    for j in chains.keys():
        total += chain_length(chains, j)
    return total

def chain_length(chains, i):
    """
    Returns the chain length. If the chain is infinite, the return value
    is 1000000000 plus the step at which the infinite loop began.
    """
    seen = {}
    return chain_length_helper(chains, i, seen)

def chain_length_helper(chains, i, seen):
    if i in seen:
        return 1000000000 - 1
    seen[i] = True
    if i in chains:
        return 1 + chain_length_helper(chains, chains[i], seen)
    else:
        return 0

def main():
  # parse cmd line flags.
  # http://docs.python.org/2/library/argparse.html#the-add-argument-method
  parser = argparse.ArgumentParser()
  parser.add_argument('-n', '--dry-run', action='store_true')
  args = parser.parse_args()

  chains = load_chains()
  #pprint.pprint(chains)
  print("BEFORE:")
  print_chains(chains)

  passcount = 0
  while True:
      t_before = count_chains(chains)
      optimize_one_pass(chains)
      t_after = count_chains(chains)
      if t_before == t_after:
          break

      passcount += 1
      print("PASS:", passcount)
      print_chains(chains)


if __name__ == '__main__':
  sys.exit(main())
