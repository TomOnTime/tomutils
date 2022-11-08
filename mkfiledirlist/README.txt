This is the fastest way to generate all.txt, files.txt and dirs.txt:

find * -print0 | mkfiledirlist

ProTip: On macOS, install GNU Find and use it instead.  
`brew install findutils` will install it as "gfind".
GNU Find is faster in this situation because it knows
a lstat() is not required for each file.  We tested
with `fs_usage` and found that "find" does the lstat no matter what.
