

vi ~/Dropbox/private/302belleville/hosts.txt 

# Setup

echo export FIOS_SESSION='SESSION'
echo export FIOS_XSRF='XSRF'

./dhcp-list.sh | ./filterdups.sh
./dhcp-list.sh | ./grepdyn.sh
./dhcp-list.sh | ./grepname.sh SonosZ
