# GO Project Instructions
## Title: Image Bucket

  - Christopher Weir - G00309429
  - Gareth Lynskey - G00312651
  - Patrick Griffin - G00314635
  - Paul Dolan - G00297086

**In cmd prompt:**
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

While that is downloading MongoDB needs to be running in the background for ImageBucket to work if you dont have mongo set up
please follow the instructions on the __***wiki***__ page above

After mongo is installed and you followed the shortcut instructions provided in the wiki or you already have it installed run the command to get mongo running
```
mongod
```
Now that mongo is running its time to move back to the project.

navigate into the project and build and run:
```
cd Go-Project
go build webapp.go
.\webapp.exe
```
Browse and upload an image

If you would like to display the image data go back to the command promt and run:
```
mongo
```
This connects to mongo
```
show dbs
```
You will find that a database has been made called Images
```
use Images
```
Changed the database

```
show collections
```
A collection called images is displayed

db.images.find()
Displays the image data that has been created

You can do the same with Usera.




#### References
https://www.tutorialspoint.com/mongodb/mongodb_gridfs.htm - Tutorial on GridFS
https://docs.mongodb.com/manual/core/gridfs/ - MongoDB Website
