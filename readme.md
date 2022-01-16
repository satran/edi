# Edi

Edi is a box for all your notes, thoughts, pictures, or any digital file. Although I use it for my daily use, I don't think it is ready for public consumption. 

You can see a [demo video](https://youtu.be/MdNUClZ_ppw) of it. It was initially called dabba.

## Why
There is a particular freedom in using plain text files to store notes and tasks. It frees you from the whims and fancies of the companies that make software for such. I have been using plain text files for a lot of things. Creating a list of groceries, keep track of the food I ate, or time I spent on some task. It is close to a paper notebook that doesn't have any format. You create the format. This was supplemented with Shell script like awk, sed, and grep. Looking at the tasks across the files were easy. But there was no cross platform way for me to use shell scripts. And thus Edi was born.

Edi does a few things. It is in the very early stage of development so don't expect every thing to work. There is no inherent system to organise, rather you are expected to use some system that works for you to organise them. 

## Install
You need a [Go](https://golang.org) compiler. With the compiler ready you just need to run
```
go install github.com/satran/edi/cmd/edi-http
```

If you have the executable in the PATH directory you can start the server by
```
edi-http -addr "localhost:8080" -dir ~/docs 
```

## Syntax
Edi has a custom syntax, some inspirations from markdown and creole.

- Heading: A `# ` can signify heading much like Markdown heading but it doesn't create a h1 tag. Rather it makes the line bold.

- Code Block: The main feature of edi is that you can evaluate a block of text as a shell script. There are two ways to do it. One is to use a inline block. For example ``date`` will evaluate date and replace text with the current date. A multiline code block can be used when you want to write a long script. This is passed as stdin to the shell environment to be evaluated. A multiline code is passed on as a code fence:
````
   ```
   echo "hello world"
   echo "new line"
   ```
````
This would add "hello world" and "new line" to the text.

- Links: Links can be created using `[[url|optional description]]`

- Images: Images follow links with a `!`: `[[!image url|optional alt text]]`

