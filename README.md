# tomutils
Various scripts that I find useful.



#### `center`

Filter: Center text.


#### `csv2tabs`

Filter: Turn CSV into TSV.


#### `dup-usb.sh`

Used for mass production of USB sticks.  If I wasn't cheap, I'd
buy a duplicator.  But, since I am cheap, I use a 4-port USB hub
on a Mac and run this script in 4 different windows.  Use at your
own risk!


#### `GETOPT_SAMPLE`

Example of how to process command line flags in a BASH script.


#### `joinlike`

joinlike combines all lines that have the same 'first part' where
the 'first part' is defined to be everything to the first whitespace.
The behaviour is somewhat akin to join, but it doesn't assume the
input is sorted. It works with files or pipes.

The Unix "join" command really should have this as a feature built-in.
It doesn't, so here is an implementation.


#### `mac2unix`

Filter: Turn Mac line endings into Unix line endings.


#### `md5tree`

Generate a TAB-separated summary of all the files in a tree.  Makes comparison easy.

```
# md5tree /mnt/RAIDVault        >/var/tmp/list.orig
# md5tree /mnt/RAIDVault-BACKUP >/var/tmp/list.backup
   # NOTE: For these next 2 lines TAB means press the TAB key.
# sort  -t'TAB' -k6 </var/tmp/list.backup >/var/tmp/list.backup.sorted
# sort  -t'TAB' -k6 </var/tmp/list.orig >/var/tmp/list.orig.sorted
# diff /var/tmp/list.orig.sorted /var/tmp/list.backup.sorted
```

The lines are tab-separated and the filename is the last thing on
the line.  That makes it easier to parse the data.


#### `PERL-DEPARSE-EXAMPLE`

I use deparse rarely enough that when I do, I need a reminder of
how to use it. This is my reminder.


#### `PYTHON_main.py`

When I start a new Python script, I start with this base.


#### `PYTHON_main_short.py`

When I start a new Python script, I start with this base if there
are no command line flags.  I always regret this because later I
end up needing command line flags.


#### `tablize`

Filter: Input TSV and output an HTML table.


#### `tabs2csv`

Filter: Turn TSV into CSV.


#### `uniq-unsorted`

Filter: Like Unix "uniq" but works on unsorted data.


#### `unsort`

Filter: LIke Unix "sort" but randomizes the order of the lines.
