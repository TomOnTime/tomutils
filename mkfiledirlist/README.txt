This is the fastest way to generate all.txt, files.txt and dirs.txt:

find * -print0 | mkfiledirlist

In theory, "find * -print0" will not "stat" each file.  However
fs_usage seems to indicate otherwise.
