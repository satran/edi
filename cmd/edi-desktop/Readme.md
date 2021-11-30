# Desktop Application

This is an experimental application. It uses [webview](https://github.com/webview/webview) to render the http page. I have tested it both on Linux and Mac and they seem to work. 

# Linux

For Linux it is quite simple, just build it and run it. If you would like to create a .desktop application it should be possible as well. 
```
go build
./edi-desktop 
```

# MacOS

```
mkdir -p EDI.app/Contents/MacOS/
go build -o EDI.app/Contents/MacOS/edi
open ./EDI.app
```

