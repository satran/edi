# Dabba

Dabba is a box for all your notes, thoughts, pictures, or any digital file. Although I use it for my daily use, I don't think it is ready for public consumption. 

You can see a [demo video](https://youtu.be/MdNUClZ_ppw) of it.

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

## Markdown
Dabba parses markdown with a few special tricks. 

- All tasks are rendered a bit more nicely. 
  So a `- [ ]` shows up as a block. A `- [s]` is also rendered. This is useful if you would like to make stages of your tasks.

- Code blocks can be evaluated.
  When you define an inline code using `` `!echo "hello world"` `` it runs the command `echo "hello world"` in a shell.
  
  Simarly a code fence 

 ````
    ```!
    var="a variable"
    echo $var
    ```
 ````

  will also be evaluated in a shell. Both use `!` as a signifier to evaluate the code. The output of both will not be parsed further and thus you can have these generating HTML.

- Internal links can be type using `(())`
  So `((Dabba Start))` will generate `<a href="Dabba Start">Dabba Start</a>`.
