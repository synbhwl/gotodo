# MENTAL MODELS TO KEEP MYSELF SANE

## how to make a post router that puts todo into json file
- have a global struct with capitalised first letter
- give all of them json tags so you can write directly into them while you're decoding using json

> it basically says, "hey, here is a json tag Title string `json:"title"`, if request's body has anyting json with that key, put it into exactly this Title field of the variable that we made of this struct type

- get the req body and parse it using json.NewDecoder(r.Body).Decode(&variable name of Task type)

> this basically says that "decode the request's json body and put that damn thing inside the variable we made. where? well, wherever is the matching tag of this json

- the newTask variable is already made, now gotta add the other fields into it
- now we gotta open the file or create a file if one doesnt exist by file, file, ferr := os.OpenFile("tasks.json", os.O_RDWR | os.O_CREATE, 0644), also dont forget to close the file

> this says "open the file of name that and make it both a reader and writer or if it doesnt exist then make one by that name. and 0644 is the permission code that you should use as a liscense"

- the "file" variable is merely a reader/writer object.
- now read from that object using filebytes, rerr := io.ReadAll(file)

> this basically says, "so im gonna read the whole file but i'll give you a slice of bytes, do whatever you'd like to do with those"

- but reading isnt enough we gotta do stuff with it
- but we cant do stuff with it until its a go variable 
- so be it, we'll make a var tasks []Task, this is a slice of type Task that will hold all the data in the the file 
- but the data are supposed to be tasks and if we just hold bytes into the new var, its useless, so we'll unmarshall the bytes into the Go objects like structs and put it inside the slice by jsonerr := json.Unmarshal(filebytes, &tasks) 

- now we can add the newTask to this slice, and we will by appending
- now thats not all, we have to write it back into the file right 
- but we can only write bytes tho
- so we'll marshall the whole tasks slice by updated, mrshlerr := json.MarshalIndent(tasks, "", " ")

- finally the updated var holds the byte we can write into the file and we will by werr := os.WriteFile("tasks.json", updated, 0644)

- and finally you can send the res to the client that the task is successfully added 

### simply speaking
- global struct for Task with tags
- make newTask var 
- parse req body and put it into newTask, the tag will handle allocation
- write rest of the fields 
- open file as reader/writer and defer the closing
- read the file into bytes 
- unmarshall the json bytes into a slice tasks var that holds all the tasks
- append the newTask to the tasks
- masrhall the variable back into json bytes
- write the bytes into file
- send the res