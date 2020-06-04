How to use our API.

Our API contains 5 functions:

1\. Registration. To register you need to go to http://127.0.0.1:8000/ and give your name as parameter.

For example:

    http://127.0.0.1:8000/?name=Bob

Output:

    Welcome, Bob! Your token is ZjvunVkYsIWMUjugGtHiqDNYF

API will give you a token.

2\. Creating a new file. To create a file you need to go to http://127.0.0.1:8000/create and give API these three parameters:

1) token - your token.
2) title - name of the file without .txt.
3) text - text of the file.

For example:

    http://127.0.0.1:8000/create?token=ZjvunVkYsIWMUjugGtHiqDNYF&title=hello&text=Hello World!

Output:

    File was created

API will create file with title hello and text of the file will be Hello World!

3\. Taking list of your files. To do this you need to go to http://127.0.0.1:8000/showfiles and give API your token.

For example:

    http://127.0.0.1:8000/showfiles?token=ZjvunVkYsIWMUjugGtHiqDNYF

Output:

    These are your files: hello bob1

4\. Taking content of the file. To do this you need to go to http://127.0.0.1:8000/showtext and give API these two parameters:

1) token - your token.
2) title - name of the file without .txt.

For example:

    http://127.0.0.1:8000/showtext?token=ZjvunVkYsIWMUjugGtHiqDNYF&title=hello

Output:

    Text of the file is: Hello World!

5\. Deleting the file. To do this you need to go to http://127.0.0.1:8000/delete and give API these two parameters:

1) token - your token
2) title - name of the file without .txt.

For example:

    http://127.0.0.1:8000/delete?token=ZjvunVkYsIWMUjugGtHiqDNYF&title=hello

Output:

    File was deleted


End of the README.md.