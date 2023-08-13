# Vidnamer


## Given a directory of filenames, generate a CSV for importing into Airtable

(cd ~/walt/pool && find * -type f -print0 | xargs -0 shasum -a 256 > ../sha256.list )

cd ~/gitthings/tomutils/vidnamer/cmd/vidfiletocsv
go install && (cd ~/walt/pool && vidfiletocsv  . ) | tee /tmp/list.csv

Import to Airtable

## Given an Airtable export, rename files

cd ~/gitthings/tomutils/vidnamer/cmd/vidrename
go install && ( cd ~/walt/pool && vidrename FILENAME.csv )

