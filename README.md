# GO Project Instructions
## Title: Image Bucket

  - Christopher Weir - G00309429
  - Gareth Lynskey - G00312651
  - Patrick Griffin - G00314635
  - Paul Dolan - G00297086
  
**We recomend you either use IntelliJ IDEA or cmder to run ImageBucket**

**Using cmd prompt or cmder:**
set the GOPATH to your workspace. For me its
```cmd
set GOPATH=C:\Users\paddy\Documents\GoCode
```
**When the GOPATH is set you need packages**
```cmd
go get gopkg.in/macaron.v1
go get gopkg.in/mgo.v2
go get gopkg.in/mgo.v2/bson
```

**git clone this repository or download this project**

For ImageBucket to work you need Mongodb. We have decided to host our mongodb database on Heroku using mLabs for the users convienience of the user.

**Using cmd prompt, cmder or IntelliJ IDEA**

Navigate into the project directory and build and run:

For cmd prompt:
```
cd Go-Project
go build webapp.go
go run webapp.go
```
For cmder:

```
cd Go-Project
go build webapp.go
.\webapp.exe
```
For IntelliJ IDEA:

Open the project through the file browser then build and run webapp.go.
Browse and upload an image, register and sign in!

To return to the image upload page simply click the ImageBucket Logo on the top left of the navigation bar.


#### References
https://www.tutorialspoint.com/mongodb/mongodb_gridfs.htm - Tutorial on GridFS
https://docs.mongodb.com/manual/core/gridfs/ - MongoDB Website
