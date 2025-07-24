# Gotodo 
Gotodo is a backend-only public JSON whiteboard essentially in form of a Todo list where anybody can perform CRUD operations and view the results.
There is no auth, RBAC separation or restriction.

## tech stack
- made entirely in Golang v1.20
- routing : gorilla mux 
- data storage : JSON file
- deployment : railway 

## features with endpoints
- user can send a GET request to "/hello" to recieve a greeting 
- user can send a GET request to "/tasks/all" to view all the current tasks added by all the users
- user can send a GET request to "/tasks/search" along with a query parameter 'title=your%20title' to find that specific task with the title and view it
- user can send a POST request to "/tasks/new" along with a JSON object in the body as '{"title":"write your title"}' to add a new task 
- user can send a DELETE request to "/tasks/delete" along with a query parameter 'title=your%20title' to delete the specific task
- user can send a PATCH request to "/tasks/edith" along with one mendatory query parameter 'old_title=your%old%title' and two optional params, 'new_title=your%new%title' (to change the title) and 'completed=bool' (to change the completed status of the task) to edit the two fields to the task
- user can send a GET request to "/tasks/filter" along with a query parameter 'completed=bool' to filter the tasks by completed status

## 3 ways to use 
### clients like curl, postman etc (on desktop)
- i had used curl so i will be listing all the commands for curl, other api clients can be used similarly

**curl commands**
*to recieve greeting*
curl -k -X GET "https://gotodo-production-58b4.up.railway.app/hello" -H "Content-Type: application/json"

*to view all tasks*
curl -k -X GET "https://gotodo-production-58b4.up.railway.app/tasks/all" -H "Content-Type: application/json"

*to search a specific task by title*
curl -k -X GET "https://gotodo-production-58b4.up.railway.app/tasks/search?title=learn%20Go%20fast" -H "Content-Type: application/json"

*to add a task*
curl -k -X POST "https://gotodo-production-58b4.up.railway.app/tasks/new" -H "Content-Type: application/json" -d '{"title": "made to be deleted"}'

*to delete a specific task by title*
curl -k -X DELETE "https://gotodo-production-58b4.up.railway.app/tasks/delete?title=task%20to%20be%20deleted" -H "Content-Type: application/json"

*to edit either title or completed status or both*
curl -k -X PATCH "https://gotodo-production-58b4.up.railway.app/tasks/edit?old_title=learn%20Go%20fast&new_title=learn%20Go%20fast&completed=true" -H "Content-Type: application/json"

*to get a filtered list of tasks on the basis of completed status*
curl -k -X GET "https://gotodo-production-58b4.up.railway.app/tasks/filter?completed=true" -H "Content-Type: application/json"

note - make changes to the query parameters/body as per your discretion. 

### using browser
I will be refering chrome. 
For get requests, you may write the url in the address bar along with the query parameters as follows (change the params as per your needs)

```
// to receive greeting
https://gotodo-production-58b4.up.railway.app/hello

// to view all tasks
https://gotodo-production-58b4.up.railway.app/tasks/all

// to search a specific task by title
https://gotodo-production-58b4.up.railway.app/tasks/search?title=learn%20Go%20fast

// to get a filtered list of tasks on the basis of completed status
https://gotodo-production-58b4.up.railway.app/tasks/filter?completed=true

```
However, since the other methods cannot be implemented similarly, i will list all the request scripts for all the operations. 

To use these, open the browser app, click F12 to open the devTools and go to the "console" tab and paste the following scripts for each operation.

```
// to receive greeting
fetch("https://gotodo-production-58b4.up.railway.app/hello", {
  method: "GET",
  headers: {
    "Content-Type": "application/json"
  }
})
  .then(res => res.json())
  .then(console.log)
  .catch(console.error);

// to view all tasks
fetch("https://gotodo-production-58b4.up.railway.app/tasks/all", {
  method: "GET",
  headers: {
    "Content-Type": "application/json"
  }
})
  .then(res => res.json())
  .then(console.log)
  .catch(console.error);

// to search a specific task by title
fetch("https://gotodo-production-58b4.up.railway.app/tasks/search?title=learn%20Go%20fast", {
  method: "GET",
  headers: {
    "Content-Type": "application/json"
  }
})
  .then(res => res.json())
  .then(console.log)
  .catch(console.error);

// to add a task
fetch("https://gotodo-production-58b4.up.railway.app/tasks/new", {
  method: "POST",
  headers: {
    "Content-Type": "application/json"
  },
  body: JSON.stringify({ title: "made to be deleted" })
})
  .then(res => res.json())
  .then(console.log)
  .catch(console.error);

// to delete a specific task by title
fetch("https://gotodo-production-58b4.up.railway.app/tasks/delete?title=task%20to%20be%20deleted", {
  method: "DELETE",
  headers: {
    "Content-Type": "application/json"
  }
})
  .then(res => res.json())
  .then(console.log)
  .catch(console.error);

// to edit either title or completed status or both
fetch("https://gotodo-production-58b4.up.railway.app/tasks/edit?old_title=learn%20Go%20fast&new_title=learn%20Go%20faster&completed=true", {
  method: "PATCH",
  headers: {
    "Content-Type": "application/json"
  }
})
  .then(res => res.json())
  .then(console.log)
  .catch(console.error);

// to get a filtered list of tasks on the basis of completed status
fetch("https://gotodo-production-58b4.up.railway.app/tasks/filter?completed=true", {
  method: "GET",
  headers: {
    "Content-Type": "application/json"
  }
})
  .then(res => res.json())
  .then(console.log)
  .catch(console.error);

```
### using mobile 
download any api client from playstore/app store, i will be using ApiClient :REST API client by Abcoderz Software from playstore.

use the method button at the top left to choose the method and the url bar to write the url.

except the POST request to add a new task, all the URL can be written in the url bar with the required query parameters

```
// *to receive greeting* - select method GET
https://gotodo-production-58b4.up.railway.app/hello

// *to view all tasks* - select method GET
https://gotodo-production-58b4.up.railway.app/tasks/all

// *to search a specific task by title* - select method GET
https://gotodo-production-58b4.up.railway.app/tasks/search?title=learn%20Go%20fast

// *to delete a specific task by title* - select method DELETE
https://gotodo-production-58b4.up.railway.app/tasks/delete?title=task%20to%20be%20deleted

// *to edit either title or completed status or both* - select method PATCH
https://gotodo-production-58b4.up.railway.app/tasks/edit?old_title=learn%20Go%20fast&new_title=learn%20Go%20faster&completed=true

// *to get a filtered list of tasks on the basis of completed status* - select method GET
https://gotodo-production-58b4.up.railway.app/tasks/filter?completed=true

//and for the POST request to add a new task, select method as POST, paste the url 
https://gotodo-production-58b4.up.railway.app/tasks/new
// and below it in the textbox write '{"title": "task title you want to add"}'
```

## vulnerabilities as of now 
- although lack of auth and rbac makes it a public playground with 0 restrictions, it also makes the project vulnerable to spam and abuse. 
- anybody can operate DELETE and POST spam 
- the data storage is a simple JSON file and may reset when the project is redeployed