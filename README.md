How to use our API

Our API contains 5 functions:

1. Registration. To register you need to go to http://127.0.0.1:8000/ and give your name as parameter. For example:

    http://127.0.0.1:8000/?name=Bob

API will give you token.

2. Creating a new file. To create a file you need to go to http://127.0.0.1:8000/create and give API these three parameters:

1) token - your token.
2) title - name of the file without .txt.
3) text - text of the file.

For example:

    http://127.0.0.1:8000/create?token=ZjvunVkYsIWMUjugGtHiqDNYF&title=hello&text=Hello World!

3. Taking list of your files. To do this you need to go to http://127.0.0.1:8000/catdir and give API your token. For example:

    http://127.0.0.1:8000/catdir?token=ZjvunVkYsIWMUjugGtHiqDNYF

4. Taking content of the file. To do this you need to go to http://127.0.0.1:8000/cat and give API these two parameters:

1) token - your token.
2) title - name of the file without .txt.

For example:

    http://127.0.0.1:8000/cat?token=ZjvunVkYsIWMUjugGtHiqDNYF&title=hello

5. Deleting the file. To do this you need to go to http://127.0.0.1:8000/del and give API these two parameters:

1) token - your token
2) title - name of the file without .txt.

For example:

    http://127.0.0.1:8000/del?token=ZjvunVkYsIWMUjugGtHiqDNYF&title=hello

End of the README.md.