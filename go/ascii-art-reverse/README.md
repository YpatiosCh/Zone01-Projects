
## Ascii-art Reverse 

Ascii-reverse consists on reversing the process from other ascii-art projects, converting the graphic representation into a text. 

- Example:
```
$ cat file.txt
 _              _   _          $
| |            | | | |         $
| |__     ___  | | | |   ___   $
|  _ \   / _ \ | | | |  / _ \  $
| | | | |  __/ | | | | | (_) | $
|_| |_|  \___| |_| |_|  \___/  $
                               $
                               $
$
$ go run . --reverse=file.txt
hello
$
```

## Installation

1. Clone the repo
```
 git clone "repo URL"
 ```
2. Change directory
```
 cd ascii-art-output
 ```
3. Git reset last version
```
 git reset --hard "specified-commit"
```

## Testing 

- To test our program as an auditor please visit [audit Documentation](https://github.com/01-edu/public/tree/master/subjects/ascii-art/reverse/audit)

- There are also unit tests for core functions of the program that someone could test running the belon command in root directory of the program :
```
go test ./file
```

## Authors 
- ychaniot (Ypatios Chaniotakos)
- mkouvara (Marinos Kouvaras)
