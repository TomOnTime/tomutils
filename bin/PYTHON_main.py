#!/usr/bin/env python

# Sample python "main" with command line flags.

import sys
import argparse

def main():
  # parse cmd line flags.
  # http://docs.python.org/2/library/argparse.html#the-add-argument-method
  parser = argparse.ArgumentParser()
  parser.add_argument('-k', '--key', required=False, metavar='string', help='Key to be used in memcached')
  parser.add_argument('-d', action='count', default=0)
  parser.add_argument('-x', action='store_true')  # Default false
  parser.add_argument('-y', action='store_false')  # Default true
  parser.add_argument('first', nargs=1)
  parser.add_argument('second', nargs=1)
  parser.add_argument('argv', nargs='*')  # Or + for "at least one"
  args = parser.parse_args()
  print args

if __name__ == '__main__':
  sys.exit(main())
