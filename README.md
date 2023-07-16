# go-sheet

Some assets are created / distributed in multiple files instead of a single spritesheet usually easier to work in engines like Godot or Unity, this go script will merge multiple files in a folder and create single file called `spritesheet.png`

```
# this will merge all files in the folder into ../input/folder/spritesheet.png
go run main.go ../input/folder/
```

![Input folder](images/input.png?raw=true "Input folder")
![Output file](images/output.png?raw=true "Output file")
