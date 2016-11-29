package main

import (
	"gopkg.in/macaron.v1"
	"gopkg.in/mgo.v2"
	"net/http"
	"fmt"
	"io/ioutil"
	"gopkg.in/mgo.v2/bson"
	"encoding/base64"
)
type Image struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	FileName    string `json:"filename" bson:"filename"`
	Encoded     string `json:"encoded" bson:"encoded"`
}

type User struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	UserName    string `json:"username" bson:"username"`
	Password     string `json:"password" bson:"password"`
	Email     string `json:"email" bson:"email"`
}

type Encoded struct {
	EncodedStr string   `json:"encoded" bson:"encoded"`
}

var response string = ""
func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Post("/", upload)
	m.Get("/Login", func(ctx *macaron.Context){
		ctx.HTML(200,"Login")
	})
	m.Combo("/Registration").
		Get(func(ctx *macaron.Context){
		ctx.HTML(200,"Registration")}).
		Post(register)
	m.Post("/Registration", register)
	m.Get("/link", func(ctx *macaron.Context) {
		// Adapted from: https://go-macaron.com/docs/middlewares/templating
		ctx.Data["Id"] = response
		ctx.HTML(200, "fileId")
	})
	m.Get("/search/:id", func(ctx *macaron.Context, w http.ResponseWriter) {
		// Adapted from: https://go-macaron.com/docs/middlewares/templating
		ctx.Data["Id"] = search(ctx.Params(":id"))
		ctx.HTML(200, "hello")
	})
	m.Run(8080)
}

func upload(w http.ResponseWriter, req *http.Request) string{
	fmt.Println("Uploadhandler start")
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	// Specify the Mongodb database
	my_db := session.DB("Images")

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
	encodedStr := base64.StdEncoding.EncodeToString([]byte(data))
	if err != nil {
		fmt.Println(err)
	}

	// Set the filename as the uploadfile name
	filename := handler.Filename
	if err != nil {
		fmt.Println(err)
	}
	img := &Image{
		Id: bson.NewObjectId(),
		FileName:  filename,
		Encoded:   encodedStr,
	}
	if err != nil {
		fmt.Println(err)
	}
	c := my_db.C("images")
	c.Insert(img)
	if err != nil {
		fmt.Println(err)
	}
	image_id := img.Id.Hex()
	response = image_id
	fmt.Println(response)
	if err != nil {
		fmt.Println(err)
	}
	/*
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
		response = file_id

		// Close the file
		err = my_file.Close()
		if err != nil {
			fmt.Println(err)
		}
	*/
	http.Redirect(w, req, "/link", 303)
	return response
}

func search(s string) string{

	img_id := s

	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	defer session.Close()
	// Specify the Mongodb database
	my_db := session.DB("Images")
	//open file from GridFS
	c := my_db.C("images")
	id:= bson.ObjectIdHex(img_id)
	encodedStr := Encoded{}
	err = c.Find(bson.M{"_id": id}).One(&encodedStr)
	if err != nil {
		panic(err)
	}
	return encodedStr.EncodedStr
}

func register(w http.ResponseWriter, req *http.Request){
	fmt.Println("Uploadhandler start")
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	// Specify the Mongodb database
	//my_db := session.DB("Images")

	// Adapted from: http://stackoverflow.com/questions/22159665/store-uploaded-file-in-mongodb-gridfs-using-mgo-without-saving-to-memory
	// Retrieve the form data
	username := req.FormValue("username")
	password := req.FormValue("password")
	email := req.FormValue("email")
	//Check if there is an error
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(username,email,password)
	my_db := session.DB("Images")
	//open file from GridFS
	c := my_db.C("users")

	user := &User{
		Id: bson.NewObjectId(),
		UserName:  username,
		Password:   password,
		Email:	email,
	}
	if err != nil {
		fmt.Println(err)
	}

	c.Insert(user)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, req, "/Login", 303)
}