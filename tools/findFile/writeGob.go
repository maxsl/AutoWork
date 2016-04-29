package findFile

import (
	"encoding/gob"
	"os"
)

func WriteGob(id string, data interface{}) error {
	File, err := os.Create("tmp/" + id + ".gob")
	if err != nil {
		return err
	}
	defer File.Close()
	En := gob.NewEncoder(File)
	return En.Encode(data)
}

func ReadGob(id string, data interface{}) error {
	File, err := os.Open("tmp/" + id + ".gob")
	if err != nil {
		return err
	}
	defer File.Close()
	De := gob.NewDecoder(File)
	return De.Decode(data)
}
