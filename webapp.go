package main

import (
	"gopkg.in/macaron.v1"
	"gopkg.in/mgo.v2"
	"net/http"
	"fmt"
	"io/ioutil"
	"gopkg.in/mgo.v2/bson"
	"bytes"
	"io"
	"image"
)

func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Post("/", upload)
	m.Get("/search/:id", func(ctx *macaron.Context) {
		// Adapted from: https://go-macaron.com/docs/middlewares/templating
		ctx.Data["id"] = search(ctx.Params(":id"))
		ctx.HTML(200, "hello")
	})
	m.Run(8080)
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

	file_id := my_file.Id().(bson.ObjectId).Hex()
	response := file_id

	// Close the file
	err = my_file.Close()
	if err != nil {
		fmt.Println(err)
	}

	// Write a log type message
	fmt.Printf("%d bytes written to the Mongodb instance\n", n)
	fmt.Println(response)
}

func search(s string) image.Image{

	file_id := s

	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	defer session.Close()
	// Specify the Mongodb database
	my_db := session.DB("Images")

	//open file from GridFS
	file, err := my_db.GridFS("fs").OpenId(bson.ObjectIdHex(file_id))
	if err != nil {
		panic(err)
	}

	//copy buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		panic(err)
	}

	//decode buffer
	img, _, err := image.Decode(&buf)

	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}
	return img
}
