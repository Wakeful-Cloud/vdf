package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	vdf "github.com/wakeful-cloud/vdf"
)

func main() {
	//Read the file
	bytes, err := ioutil.ReadFile("./map-test.vdf")

	if err != nil {
		panic(err)
	}

	//Read VDF
	vdfMap, err := vdf.ReadVdf(bytes)

	if err != nil {
		panic(err)
	}

	//Covert to JSON
	rawJSON, err := json.Marshal(vdfMap)

	if err != nil {
		panic(err)
	}

	//Log
	log.Print(string(rawJSON))

	//Write VDF
	rawVdf, err := vdf.WriteVdf(vdfMap)

	if err != nil {
		panic(err)
	}

	//Write the file
	ioutil.WriteFile("./out.vdf", rawVdf, 0666)
}
