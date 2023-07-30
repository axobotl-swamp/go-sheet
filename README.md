# go-sheet

## Action mode

Some assets are created / distributed in multiple files instead of a single spritesheet usually easier to work in engines like Godot or Unity, this go script will merge multiple files in a folder and create single file called `spritesheet.png`

```
# this will merge all files in the folder into ../input/folder/spritesheet.png
go run main.go action ../input/folder/
```

![Input folder](images/input.png?raw=true "Input folder")
![Output file](images/output.png?raw=true "Output file")

## Individual
If you have a single file for each sprite, create a folder for each action and run

```
# each folder will be a row and each file will be a column
go run main.go individual ../input/folder/
```
![Input folder](images/individual.png?raw=true "Individual mode structure folder")
