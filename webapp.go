package main

import (
	"gopkg.in/macaron.v1"
	"gopkg.in/mgo.v2"
	"net/http"
	"fmt"
	"io/ioutil"
)

func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Post("/", upload)
	m.Run()
}

func upload(req *http.Request) {
	fmt.Println("Uploadhandler start")
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	defer session.Close()
// Adapted from: http://stackoverflow.com/questions/22159665/store-uploaded-file-in-mongodb-gridfs-using-mgo-without-saving-to-memory
	// Retrieve the form file 
	file, handler, err := req.FormFile("uploadfile")
	//Check if there is an error
	if err != nil {
		fmt.Println(err)
	}
	// Read the file into memory
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	// Specify the Mongodb database
	my_db := session.DB("Images")
	// Set the filename as the uploadfile name
	filename := handler.Filename

	// Create the file in the Mongodb Gridfs instance
	my_file, err := my_db.GridFS("fs").Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	// Write the file to the Mongodb Gridfs instance
	n, err := my_file.Write(data)
	if err != nil {
		fmt.Println(err)
	}

	// Close the file
	err = my_file.Close()
	if err != nil {
		fmt.Println(err)
	}

	// Write a log type message
	fmt.Printf("%d bytes written to the Mongodb instance\n", n)
}
