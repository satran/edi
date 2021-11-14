# Dabba

Dabba is a box for all your notes, thoughts, pictures, or any digital file. Although I use it for my daily use, I don't think it is ready for public consumption. 

## Why
There is a particular freedom in using plain text files to store notes and tasks. It frees you from the whims and fancies of the companies that make software for such. I have been using plain text files for a lot of things. Creating a list of groceries, keep track of the food I ate, or time I spent on some task. It is close to a paper notebook that doesn't have any format. You create the format. This was supplemented with Shell script like awk, sed, and grep. Looking at the tasks across the files were easy. But there was no cross platform way for me to use shell scripts. And thus Dabba was born. Dabba, in Hindi, means a box. A box of all your things. 

Dabba does a few things. It is in the very early stage of development so don't expect every thing to work. There is no inherent system to organise, rather you are expected to use some system that works for you to organise them. 

## Install
You need a [Go](https://golang.org) compiler. With the compiler ready you just need to run
```
go install github.com/satran/dabba
```

If you have the executable in the PATH directory you can start the server by
```
dabba -addr "localhost:8080" -dir ~/docs 
```

## Template engine
It uses a template engine derived from Go's template package. Unlike the default parser it uses (("(( ))")). There are a few builtin functions.

- link | l - creates a link
```
((link "https://satran.in" "My Blog"))
generates
<a href="https://satran.in">My Blog</a>
```

- image | i - embed an image
```
((image "hello.png" "Caption"))
generates
<img src="hello.png" alt="Caption" />
```

- sh | runs a shell script in the dabba directory 
```
((sh `grep -ri '- \[ \]' .`))
Finds all the "todos" in your files.
```
