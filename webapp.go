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
//Structures-----------------------------------------------------------------------------
//Image Upload Structure
type Image struct {
	ImageId     string `json:"imageid" bson:"imageid"`
	FileName    string `json:"filename" bson:"filename"`
	Encoded     string `json:"encoded" bson:"encoded"`
	User     string `json:"user" bson:"user"`
}
//User Image Retrieval Structure
type UserImage struct {
	ImageId     string `json:"imageid" bson:"imageid"`
	FileName    string `json:"filename" bson:"filename"`
	Encoded     string `json:"encoded" bson:"encoded"`
	User     string `json:"user" bson:"user"`
}
//User Registration & Login Structure
type User struct {
	Id          bson.ObjectId `bson:"_id,omitempty"`
	UserName    string `json:"username" bson:"username"`
	Password     string `json:"password" bson:"password"`
	Email     string `json:"email" bson:"email"`
}
//Uploaded Image Link Retrieval Structure
type Encoded struct {
	EncodedStr string   `json:"encoded" bson:"encoded"`
}

//Global Variables-----------------------------------------------------------------------------
var response string = ""//response used for sending back the link for the uploaded Image
var LoginError string = ""//LoginError used for when the user enters incorrect login details
var UserDetails string = "null"//UserDetails used to check if a user is logged in - set to null as default

//Functions-------------------------------------------------------------------------------------
//Main Function-------------------------------------------------------------------------------------
func main() {
	m := macaron.Classic()
	m.Use(macaron.Renderer())
	// Some of the following functions have been adapted from: https://go-macaron.com/docs/middlewares/templating
	//Handler for /
	m.Post("/", upload)//Call upload func when Post method is activated from the root page

	//Combo Get and Post for the Home page
	m.Combo("/Home"). Get(func(ctx *macaron.Context){ //Handler for /Home
		ctx.Data["Auth"] = UserDetails//Send UserDetails to the Home template page to check if user is logged in
		ctx.HTML(200,"Home")}). //Load the Home template
		Post(upload) //Call upload func when post method is activated from the Home page

	//Get for the Logout page
	m.Get("/Logout", func(ctx *macaron.Context){ //Handler for /Logout
		UserDetails = "null" //Set UserDetails back to null when user signs out
		ctx.Data["Auth"] = UserDetails //Send UserDetails to the Logout template page to check if user is logged in
		ctx.HTML(200,"Logout")}) //Load the Logout template

	//Combo Get and Post for the Login page
	m.Combo("/Login").//Handler for /Login
		Get(confirmUser, func(ctx *macaron.Context){ //Call confirmUser func to check if a user is logged in
		ctx.Data["Error"] = LoginError //Send LoginError to the Login template page when the user enters incorrect login details
		ctx.HTML(200,"Login")}). //Load the login template
		Post(login) //Call login func when post method is activated from the Login page

	//Combo Get and Post for the Registration page
	m.Combo("/Registration").Get(confirmUser,func(ctx *macaron.Context){ //Handler for /Registration
		ctx.HTML(200,"Registration")}). //Load the Registration template
		Post(register) //Call register func when post method is activated from the Registration page

	//Get for the MyImages page
	m.Get("/MyImages", func(ctx *macaron.Context){ //Handler for /MyImages
		ctx.Data["Auth"] = UserDetails //Send UserDetails to the MyImages page to check if user is logged in
		ctx.Data["ImageList"] = userImages(nil,nil) //Send the return value from func userImages (the users images) to the MyImages template
		ctx.HTML(200,"MyImages")}) //Load the MyImages template

	//Get for the link page
	m.Get("/link", func(ctx *macaron.Context) { //Handler for /link
		ctx.Data["Auth"] = UserDetails //Send UserDetails to the Link page to check if user is logged in
		ctx.Data["Id"] = response // Send the uploaded image id to the FileId template page
		ctx.HTML(200, "fileId") //Load the FileId template
	})

	//Get for the search page
	m.Get("/search/:id", func(ctx *macaron.Context, w http.ResponseWriter) { //Handler for /search that takes in paramaters
		ctx.Data["Id"] = search(ctx.Params(":id"))// Send the return value of the search func (the requested image data) to the Image template page
		ctx.HTML(200, "Image") //Load the Image template
	})

	m.Run(8080) //Run on port 8080
}
//Upload Function-------------------------------------------------------------------------------------
func upload(w http.ResponseWriter, req *http.Request) string{
	//Connect to mLabs Mongodb database on Heroku
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
	//Encode data into base64 and assign it to string encodedStr
	encodedStr := base64.StdEncoding.EncodeToString([]byte(data))
	if err != nil {
		fmt.Println(err)
	}

	// Set the filename as the uploadfile name
	filename := handler.Filename
	if err != nil {
		fmt.Println(err)
	}
	//Set image_id as a new Bson Object Id and convert it to Hex
	image_id:=bson.NewObjectId().Hex()
	//Assign img as the structure Image and populate it with variables
	img := &Image{
		ImageId: image_id,
		FileName:  filename,
		Encoded:   encodedStr,
		User:	UserDetails,
	}
	if err != nil {
		fmt.Println(err)
	}
	//Set the database collection
	c := my_db.C("images")
	//Insert img into the database
	c.Insert(img)
	if err != nil {
		fmt.Println(err)
	}
	//Set response as the image_id
	response = image_id
	if err != nil {
		fmt.Println(err)
	}
	//Redirect the user to the /link page
	http.Redirect(w, req, "/link", 303)
	//Return response (string)
	return response
}

//Search Function-------------------------------------------------------------------------------------
func search(s string) string{
	//Search func takes in a string s
	//Assigns img_id as s value
	img_id := s

	//Connect to mLabs Mongodb database on Heroku
	session, err := mgo.Dial("mongodb://test:test@ds113958.mlab.com:13958/heroku_t76cfn1s")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Specify the Mongodb database
	my_db := session.DB("heroku_t76cfn1s")

	//Set the database collection
	c := my_db.C("images")
	//encoded is a new Encoded{} struct
	encodedStr := Encoded{}
	//Query database for the image matching the img_id and return the first record
	//Assign the query results to encodedStr
	err = c.Find(bson.M{"imageid": img_id}).One(&encodedStr)
	if err != nil {
		panic(err)
	}
	//Return the encodedStr struct value EncodedStr (string)
	return encodedStr.EncodedStr
}

//Register Function-------------------------------------------------------------------------------------
func register(w http.ResponseWriter, req *http.Request){
	//Connect to mLabs Mongodb database on Heroku
	session, err := mgo.Dial("mongodb://test:test@ds113958.mlab.com:13958/heroku_t76cfn1s")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Specify the Mongodb database
	my_db := session.DB("heroku_t76cfn1s")

	//Set the database collection
	c := my_db.C("users")

	// Adapted from: http://stackoverflow.com/questions/22159665/store-uploaded-file-in-mongodb-gridfs-using-mgo-without-saving-to-memory
	// Retrieve the form data
	username := req.FormValue("username")
	password := req.FormValue("password")
	email := req.FormValue("email")
	//Check if there is an error
	if err != nil {
		fmt.Println(err)
	}
	//Assign user as the structure User and populate it with variables
	user := &User{
		Id: bson.NewObjectId(),
		UserName:  username,
		Password:   password,
		Email:	email,
	}
	if err != nil {
		fmt.Println(err)
	}
	//Insert user into the database
	c.Insert(user)
	if err != nil {
		fmt.Println(err)
	}
	//Redirect the user to the Login Page
	http.Redirect(w, req, "/Login", 303)
}

//Login Function-------------------------------------------------------------------------------------
func login(w http.ResponseWriter, req *http.Request) string{
	//Connect to mLabs Mongodb database on Heroku
	session, err := mgo.Dial("mongodb://test:test@ds113958.mlab.com:13958/heroku_t76cfn1s")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Specify the Mongodb database
	my_db := session.DB("heroku_t76cfn1s")

	// Adapted from: http://stackoverflow.com/questions/22159665/store-uploaded-file-in-mongodb-gridfs-using-mgo-without-saving-to-memory
	// Retrieve the form data
	email := req.FormValue("email")
	password := req.FormValue("password")

	//Check if there is an error
	if err != nil {
		fmt.Println(err)
	}
	//Specify the database collection
	c := my_db.C("users")
	//Assign auth as the structure User
	auth := User{}
	//Query database for the email matching the email entered and return the first record
	//Assign the query results to auth
	err = c.Find(bson.M{"email": email}).One(&auth)
	//Check if the inserted password matches the password that retrieved from the database
	if password == auth.Password {
		//If the passwords match then redirect to /MyImages
		http.Redirect(w, req, "/MyImages", 303)
		//Set the LoginError variable to default
		LoginError = ""
		//Set the UserDetails to the email entered
		UserDetails = auth.Email
	} else {
		//If the passwords did not match then set LoginError to "Incorrect Details"
		LoginError = "Incorrect Details"
		//Redirect the user back to the login page to try again
		http.Redirect(w, req, "/Login", 303)
	}
	//Return the LoginError
	return LoginError
}

//ConfirmUser Function-------------------------------------------------------------------------------------
func confirmUser(w http.ResponseWriter, req *http.Request){
	//Check if a user is logged in and if they are, then redirect to the MyImages page
	if UserDetails !="null" {
		http.Redirect(w, req, "/MyImages", 303)
	}
}

//UserImages Function-------------------------------------------------------------------------------------
func userImages(w http.ResponseWriter, req *http.Request) string {
	//Connect to mLabs Mongodb database on Heroku
	session, err := mgo.Dial("mongodb://test:test@ds113958.mlab.com:13958/heroku_t76cfn1s")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Specify the Mongodb database
	my_db := session.DB("heroku_t76cfn1s")
	//Check if there is an error
	if err != nil {
		fmt.Println(err)
	}
	//Set the database collection
	c := my_db.C("images")
	//Set listImage as a UserImage struct array
	var listImage []UserImage
	//Query the database to find all the records that are associated with the users email (UserDetails)
	//Set query results to listImage
	err = c.Find(bson.M{"user": UserDetails}).All(&listImage)
	imagesList := &listImage
	//Json marshel imageList and assign its value to images
	images, err := json.Marshal(imagesList)
	//Return the string value of images
	return string(images)

}

