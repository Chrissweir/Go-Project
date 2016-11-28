# GO Project Documentaion 
## Title: Image Bucket

  - Christopher Weir - G00309429
  - Gareth Lynskey - G00312651
  - Patrick Griffin - G00314635
  - Paul Dolan - G00297086

### Introduction
In summary, this a single Web application that a user can upload an image to be stored and can be retrieve in a later date. The user simply drags and drops the image or can "Choose file" to be uploaded.

### Technologies used

Languages | Libraries | Frameworks | Database
------------ | ------------- | ------------- | -------------
GO| Bootstrap| Macaron | MongoDB
HTML | - | - |- 
JavaScript | - | - | - 
CSS | - | - | - 


### Database MongoDB - Backend
This Project inherits the use of MongoDB as a backend for storing and retrieving images through the use of GridFS. GridFS is to store and to retrieve large files such as images and videos or audio files Its data is stored within MongoDB collections. GridFS strores files even greater than its document size of 16mb. 
**Why we used MongoDB**
MongoDB stores data in collections which is more beneficial for this project as a relationship database would be ineficient in populating information MongoDB has a built in easy solution for for sharing/partitioning your database rather than MySQL table performance will degrade when crossing the 5-10GB per table.
##### MongoDB User Guide - Set up
Download and install mongo depending on you operating system from the [Mongo Website](https://www.mongodb.com/download-center?jmp=nav#community) and follow the installation processes. When MongoDB is downloaded and installed locate the mogoDB folder. It's usually in **C:\Program Files\MongoDB** and go into the **Server** folder then click **3.2** folder then into the **Bin** folder 
>C:\Program Files\MongoDB\Server\3.2\bin 
>

**Note:** Path of bin may change if mongo updates in the future
In the bin folder there will be a bunch of mongo executables. **MongoD** located in the bin folder is the actual database which is gonna run in the background **mongo** (commandline interface) is the application where you create/insert/update/delete your database.
open up cmd prompt and navigate to the bin directory
>C:\Program Files\MongoDB\Server\3.2\bin
>
**dir** to verify all the contents

If its your first time setting up MongoDB and execute **mongod** you will  find that it will tell you that you need a data directory. So, now we need a data directory
execute:
>mkdir \data\db
>

Now run **mongod** and it will show the port 27017 and run **mongo** on another cmd prompt and it'll show the connection.

##### Creating a shortcut for Mongo
Instead of locating the bin folder to execute the mongo executables everytime you can create a shortcut.
navigate back to 
>C:\Program Files\MongoDB\Server\3.2\bin
>

right click any file and click properties copy the location:
>C:\Program Files\MongoDB\Server\3.2\bin
>
so you tell the computer everytime you type mongod look to this location
that way you dont have to type it out everythime in the cmd line

- go to Control **Panel\System and Security\System**
- go to **advanced system settings**
- go to **environmental variables**
- create new variable called **PATH**

Now paste in for variable value **C:\Program Files\MongoDB\Server\3.2\bin**
now you can type mongod and it now starts up without having to navigate to the location every time. Open up another cmd line and type mongo and now ready to type commands.

#### Why we tried GridFS
We started using GridFS to store and to retrieve large files such as images and videos or audio files. Its data is stored within MongoDB collections. GridFS strores files even greater than its document size of 16mb.

GridFS divides a file into **chunks** and **files** and stores each chunk of data in a seperate document Chunks stores the binary chunks files stores the files metadata
- fs.files
- fs.chunks

The form for chunks:
```json
{
  "_id" : <ObjectId>,
  "files_id" : <ObjectId>,
  "n" : <num>,
  "data" : <binary>
}
```
The form for files:
```json
{
  "_id" : <ObjectId>,
  "length" : <num>,
  "chunkSize" : <num>,
  "uploadDate" : <timestamp>,
  "md5" : <hash>,
  "filename" : <string>,
  "contentType" : <string>,
  "aliases" : <string array>,
  "metadata" : <dataObject>,
}
```

We had images posting to the database and we were able to retrieve the data. However we were unable to convert this data back into the image for the user to view.
```go
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
```

#### References
https://www.tutorialspoint.com/mongodb/mongodb_gridfs.htm - Tutorial on GridFS
https://docs.mongodb.com/manual/core/gridfs/ - MongoDB Website
