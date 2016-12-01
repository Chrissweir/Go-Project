package main

import (
	"gopkg.in/macaron.v1"
	"gopkg.in/mgo.v2"
	"net/http"
	"fmt"
	"io/ioutil"
	"gopkg.in/mgo.v2/bson"
	"encoding/base64"
	"encoding/json"
)
type Image struct {
	ImageId     string `json:"imageid" bson:"imageid"`
	FileName    string `json:"filename" bson:"filename"`
	Encoded     string `json:"encoded" bson:"encoded"`
	User     string `json:"user" bson:"user"`
}

type UserImage struct {
	ImageId     string `json:"imageid" bson:"imageid"`
	FileName    string `json:"filename" bson:"filename"`
	Encoded     string `json:"encoded" bson:"encoded"`
	User     string `json:"user" bson:"user"`
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
var LoginError string = ""
var UserDetails string = "null"

func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	m.Post("/", upload)
	m.Combo("/Home").
		Get(func(ctx *macaron.Context){
		ctx.Data["Auth"] = UserDetails
		ctx.HTML(200,"Home")}).
		Post(upload)

	m.Get("/Logout", func(ctx *macaron.Context){
		UserDetails = "null"
		ctx.Data["Auth"] = UserDetails
		ctx.HTML(200,"Logout")
	})

	m.Combo("/Login").
		Get(confirmUser, func(ctx *macaron.Context){
		ctx.Data["Error"] = LoginError
		ctx.HTML(200,"Login")}).
		Post(login)

	m.Combo("/Registration").
		Get(confirmUser,func(ctx *macaron.Context){
		ctx.HTML(200,"Registration")}).
		Post(register)

	m.Get("/MyImages", func(ctx *macaron.Context){
		ctx.Data["Auth"] = UserDetails
		ctx.Data["ImageList"] = userImages(nil,nil)
		ctx.HTML(200,"MyImages")})

	m.Get("/link", func(ctx *macaron.Context) {
		// Adapted from: https://go-macaron.com/docs/middlewares/templating
		ctx.Data["Auth"] = UserDetails
		ctx.Data["Id"] = response
		ctx.HTML(200, "FileId")
	})

	m.Get("/search/:id", func(ctx *macaron.Context, w http.ResponseWriter) {
		// Adapted from: https://go-macaron.com/docs/middlewares/templating
		ctx.Data["Id"] = search(ctx.Params(":id"))
		ctx.HTML(200, "Image")
	})
	m.Run(8080)
}

func upload(w http.ResponseWriter, req *http.Request) string{
	fmt.Println("Uploadhandler start")
	session, err := mgo.Dial("mongodb://test:test@ds113958.mlab.com:13958/heroku_t76cfn1s")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	// Specify the Mongodb database
	my_db := session.DB("heroku_t76cfn1s")

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
	image_id:=bson.NewObjectId().Hex()
	img := &Image{
		ImageId: image_id,
		FileName:  filename,
		Encoded:   encodedStr,
		User:	UserDetails,
	}
	if err != nil {
		fmt.Println(err)
	}
	c := my_db.C("images")
	c.Insert(img)
	if err != nil {
		fmt.Println(err)
	}

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

	session, err := mgo.Dial("mongodb://test:test@ds113958.mlab.com:13958/heroku_t76cfn1s")
	if err != nil {
		panic(err)
	}

	defer session.Close()
	// Specify the Mongodb database
	my_db := session.DB("heroku_t76cfn1s")
	//open file from GridFS
	c := my_db.C("images")
	encodedStr := Encoded{}
	err = c.Find(bson.M{"imageid": img_id}).One(&encodedStr)
	if err != nil {
		panic(err)
	}
	return encodedStr.EncodedStr
}

func register(w http.ResponseWriter, req *http.Request){
	fmt.Println("Uploadhandler start")
	session, err := mgo.Dial("mongodb://test:test@ds113958.mlab.com:13958/heroku_t76cfn1s")
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
	my_db := session.DB("heroku_t76cfn1s")
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

func login(w http.ResponseWriter, req *http.Request) string{

	session, err := mgo.Dial("mongodb://test:test@ds113958.mlab.com:13958/heroku_t76cfn1s")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	// Specify the Mongodb database
	//my_db := session.DB("Images")

	// Adapted from: http://stackoverflow.com/questions/22159665/store-uploaded-file-in-mongodb-gridfs-using-mgo-without-saving-to-memory
	// Retrieve the form data
	email := req.FormValue("email")
	password := req.FormValue("password")

	//Check if there is an error
	if err != nil {
		fmt.Println(err)
	}
	my_db := session.DB("heroku_t76cfn1s")
	//open file from GridFS
	c := my_db.C("users")
	auth := User{}
	err = c.Find(bson.M{"email": email}).One(&auth)

	if password == auth.Password {
		http.Redirect(w, req, "/MyImages", 303)
		LoginError = ""
		UserDetails = auth.Email
	} else {
		LoginError = "Incorrect Details"
		http.Redirect(w, req, "/Login", 303)
	}

	fmt.Println(auth.Email, auth.Password, auth.UserName)
	return LoginError
}

func confirmUser(w http.ResponseWriter, req *http.Request){
	fmt.Println(UserDetails)
	if UserDetails !="null" {
		http.Redirect(w, req, "/MyImages", 303)
	}
}

func userImages(w http.ResponseWriter, req *http.Request) string {
	session, err := mgo.Dial("mongodb://test:test@ds113958.mlab.com:13958/heroku_t76cfn1s")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	// Specify the Mongodb database
	//Check if there is an error
	if err != nil {
		fmt.Println(err)
	}
	my_db := session.DB("heroku_t76cfn1s")
	//open file from GridFS
	c := my_db.C("images")

	var listImage []UserImage
	err = c.Find(bson.M{"user": UserDetails}).All(&listImage)

	fmt.Println(&listImage)
	imagesList := &listImage
	images, err := json.Marshal(imagesList)
	return string(images)

}

