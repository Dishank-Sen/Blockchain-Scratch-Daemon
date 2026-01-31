package utils

import (
	"encoding/json"
	"os"
	"path"

	"github.com/Dishank-Sen/Blockchain-Scratch-Daemon/types"
)

func SaveID(id string) error{
	dir := "/tmp"
	path := path.Join(dir, "blocd.json")
	data := types.Metadata{
		ID: id,
	}

	byteData, err := json.Marshal(data)
	if err != nil{
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil{
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil{
		return err
	}
	defer f.Close()
	_, err = f.Write(byteData)
	if err != nil{
		return err
	}
	return nil
}

func GetID() (string, error){
	byteData, err := os.ReadFile("/tmp/blocd.json")
	if err != nil{
		return "", err
	}
	var mdata types.Metadata
	if err := json.Unmarshal(byteData, &mdata); err != nil{
		return "", err
	}
	return mdata.ID, nil
}