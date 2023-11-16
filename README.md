# goProject
An API having endpoints to signup, login, and perform CRUD operations on notes for a specific user <br>
In this API, I have used a data structure instead of a database.
I have used Postman for testing this API.
These are the ways to use Postman for testing it:
Signup:

Set the request type to POST.
Enter the URL: http://localhost:8080/signup.
Set the request body to raw and select JSON (application/json).
Enter the JSON data for signup:
JSON like this:
{
  "name": "Deepak Sharma",
  "email": "deepakatunique@gmail.com",
  "password": "abc_123"
}

Login:

Set the request type to POST.
Enter the URL: http://localhost:8080/login.
Set the request body to raw and select JSON (application/json).
Enter the JSON data for login:
json like this : 
{
  "email": "deepakatunique@gmail.com",
  "password": "abc_123"
}
This will return a session ID (sid).


List Notes:

Set the request type to GET.
Enter the URL: http://localhost:8080/notes.
In the Headers section, add a new header with key Content-Type and value application/json.
In the Params section, add a parameter with key sid and value as the session ID copied from the login response.


Create Note:

Set the request type to POST.
Enter the URL: http://localhost:8080/notes.
Set the request body to raw and select JSON (application/json).
Enter the JSON data for creating a note:
json like this : 
{
  "sid": "<your session id which was generated when you login>",
  "note": "This is a new note"
}
This will return the ID of the newly created note.



Delete Note:

Set the request type to DELETE.
Enter the URL: http://localhost:8080/notes.
Set the request body to raw and select JSON (application/json).
Enter the JSON data for deleting a note:
json like this : 
{
  "sid": "<your-session-id>",
  "id": <note-id>
}
