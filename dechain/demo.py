#!/usr/bin/env python3

import sys
import argparse

class Chains:

  graph = {}  # key = URL. value = next in chain.
  changed = set()  # set of which URLs have changed.

  def load_from_api(self):
    """Loads the data.  Replace this with something that loads from the API."""
    self.graph = {
        # Normal links:
        'A': 'X',
        'B': 'Y',
        'C': 'Z',
        # Long chain: D->E->F->G->H->I-J
        # (due to ordering, this will require multiple passes)
        'D': 'E',
        'E': 'F',
        'F': 'G',
        'G': 'H',
        'H': 'I',
        'I': 'J',
        # Another chain: O->N->M->L->K
        # (due to ordering, this will resolve in 1 pass)
        'L': 'K',
        'M': 'L',
        'N': 'M',
        'O': 'N',
        # Infinite loop:
        'R': 'R',
    }

  def flatten_once(self):
    """The 1-pass flattening.
    For any chains that are longer than 1 hop, reduces them by
    at least 1 hop (maybe more).
    """
    for i in sorted(self.graph.keys()):  # Sorted for stable unit tests.
      n = self.graph[i]
      if n in self.graph:
        self.graph[i] = self.graph[n]
        self.changed.add(i)

  def print_all(self):
    """Print chains."""
    for k in sorted(self.graph.keys()):
      if k in self.changed:
        c = ' * '
      else:
        c = '   '
      print(k, self.graph[k], c, self.chain_length(k), self.chainToString(k))

  def total_length(self):
    """Returns the sum of the lengths of all chains."""
    total = 0
    for k in self.graph.keys():
      total += self.chain_length(k)
    return total

  def chain_length(self, k):
    """
    Returns the chain length. If the chain is infinite, the return value
    is the step number at which the infinite loop began.
    TODO(tlim): Infinite loops should result in an exception.
    """
    seen = {}
    return self._chain_length_helper(k, seen)

  def _chain_length_helper(self, k, seen):
    if k in seen:
      return 0
    seen[k] = True
    if k in self.graph:
      return 1 + self._chain_length_helper(self.graph[k], seen)
    else:
      return 0

  def chainToString(self, k):
    """
    Returns the chain as a string like "A->B->C".
    """
    seen = {}
    return self._chainToString_helper(k, seen)

  def _chainToString_helper(self, k, seen):
    if k in seen:
      return "REPEATING"
    seen[k] = True
    if k in self.graph:
      return str(k) + "->" + self._chainToString_helper(self.graph[k], seen)
    else:
      return str(k)

  def changed_keys(self):
    """Iterator that reveals the changed items."""
    for it in sorted(self.graph.items()):
      if it[0] in self.changed:
        yield it


def main():
  # parse cmd line flags.
  # http://docs.python.org/2/library/argparse.html#the-add-argument-method
  parser = argparse.ArgumentParser()
  parser.add_argument('--real', '-r', action='store_true')
  args = parser.parse_args()

  g = Chains()
  g.load_from_api()
  print("BEFORE:")
  g.print_all()

  passcount = 0
  while True:
    t_before = g.total_length()
    g.flatten_once()
    t_after = g.total_length()
    # No change? Must be at the optimal state.
    if t_before == t_after:
      break
    passcount += 1
    print("PASS:", passcount)
    g.print_all()

  print("Use the API to make this changes:")
  for chain in g.changed_keys():
    if args.real:
      print(chain, "MAKING CHANGE (not implemented)")
    else:
      print(chain, "DEBUG MODE. NOT MAKING CHANGE.")

if __name__ == '__main__':
  sys.exit(main())
