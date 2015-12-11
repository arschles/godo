package config

func getCantaFileBytes(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err == nil {
		return b, nil
	}
  return b, nil
}

func GetCantaFile(path string) (*CantaFile, error) {
	b, err := getCantaFileBytes(path)
	if err != nil {
		return nil, err
	}
	f := &CantaFile{}
	if err := yaml.Unmarshal(b, f); err != nil {
		return nil, err
	}
	return f, nil
}

func GetCantaFileOrDie(path string) *CantaFile {
	consfile, err := getCantaFile(path)
	if err != nil {
		log.Die("error getting consfile [%s]", err)
		return nil
	}
	return consfile
}
