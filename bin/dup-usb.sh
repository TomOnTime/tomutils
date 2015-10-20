#!/bin/bash -xe

# Copy a USB stick on a MacOS X host.

# Directory to copy onto each USB stick:
SOURCE=~/USB

DEVNUM=$1
DISK=/dev/disk$DEVNUM
PART=${DISK}s1
DEST=/tmp/mnt$DEVNUM

while true ; do

  # Wait for a USB stick to be inserted:
  while ! mount|fgrep -s $DISK >/dev/null ; do
    echo Waiting for USB stick in $DISK
    sleep 5
  done

  # Erase disk, unmount it, remount it under a predictable name:
  #diskutil erasedisk MS-DOS PICC2012 $DISK
  #diskutil erasedisk FAT32 LOPSAEAST13 MBR $DISK
  diskutil erasedisk FAT32 LOPSAEAST14 MBR $DISK
  diskutil umountDisk $DISK
  mkdir -p "$DEST" && diskutil mount -mountPoint "$DEST" "$PART"
  
  # Copy:
  cd "$SOURCE" && tar cf - . | (cd "$DEST" && tar xpf - )
  
  # Remove cruft and eject.
  sudo mdutil -i off -E  "$DEST"
  ls "$DEST"
  ( cd "$DEST" && find . -name ._\* -delete )
  ( cd "$DEST" && rm -rf .fseventsd .Spot* ._* .Trash* )
  ls -a "$DEST"
  diskutil umountdisk $DISK
  diskutil eject $DISK

done
