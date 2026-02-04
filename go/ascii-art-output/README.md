
# Ascii-Art Output

Ascii-art is a program which consists in receiving a string as an argument and outputting the string in a graphic representation using ASCII.

Output features include:

Specify a txt file to write the output to. If txt file is not specified, the output will be printed to terminal.

Specify a font that the output will have. If font is not specified, standard font will be used by default.



## Documentation

[Documentation](https://platform.zone01.gr/intra/athens/div-01/output?event=200)


## Built with

* ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
## Installation

1. Clone the repo

```bash
 git clone "repo URL"
```
    
2. Change directory

```bash
 cd ascii-art-output
```

3. Git reset last version

```bash
 git reset --hard "specified-commit"
```
## Usage/Examples
Please visit [documentation](https://platform.zone01.gr/intra/athens/div-01/output?event=200) in order to see all the examples.

Simple Examples:

1.
```bash
    go run . --output=banner.txt Hello shadow
```
This command will write the word "Hello" to a txt file named banner with a shadow font.

2.
```bash
    go run . Hello shadow
```
This command will write the word "Hello" to terminal with a shadow font.

3.
```bash
    go run . -output banner.txt Hello shadow
```
This command will exit the program printing the message:
```bash
    Usage: go run . [OPTION] [STRING] [BANNER]

    EX: go run . --output=<fileName.txt> something standard
```



## Running Tests

The main_test.go tests the "backbone" functions of the program in order to run as needed.

To run all the tests, run the following command:
```bash
  go test 
```

To run a test in a specific function, run the following command:
```bash
  go test -run "FunctionName"
```


## Contact

- ychaniot (Ypatios Chaniotakos)
- mkouvara (Marinos Kouvaras)

Project [Link](https://platform.zone01.gr/git/ychaniot/ascii-art-output)

