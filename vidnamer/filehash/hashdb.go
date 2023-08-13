package filehash

type DB struct {
	data      []Info
	fileToSig map[string]string
	sigToFile map[string]string
}

func Initialize(hashlistFile string) (*DB, error) {
	sha256db, err := FromFile(hashlistFile)
	if err != nil {
		return nil, err
	}

	return FromInfoList(sha256db), nil
}

func FromInfoList(items []Info) *DB {

	hdb := DB{}
	hdb.data = items
	hdb.fileToSig = map[string]string{}
	hdb.sigToFile = map[string]string{}

	for _, item := range items {
		//fmt.Printf("DEBUG: FromInfo: hdb.fileToSig[%q] = %q\n", item.Filename, item.Signature)
		hdb.fileToSig[item.Filename] = item.Signature
		hdb.sigToFile[item.Signature] = item.Filename
	}

	return &hdb
}

func (hdb *DB) GetSignature(file string) string {
	//fmt.Printf("DEBUG: GetSignature(%q)\n", file)
	return hdb.fileToSig[file]
}


func (hdb *DB) SigToFilename(sig string) string {
	//fmt.Printf("DEBUG: SigToFilename(%q)\n", file)
	return hdb.sigToFile[sig]
}
