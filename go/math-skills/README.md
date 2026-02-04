
# Math skills

The purpose of this project is for you to calculate the following from a docker-created txt file:

1. Average
2. Median
3. Variance
4. Standard Deviation



## Documentation

[Documentation](https://platform.zone01.gr/intra/athens/div-01/math-skills?event=200)


## Built with

* ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
## Installation

1. Clone the repo

```bash
 git clone "https://platform.zone01.gr/git/ychaniot/math-skills.git"
```
    
2. Change directory

```bash
 cd math-skills
```

3. Git reset last version

```bash
 git reset --hard "specified-commit"
```
## Usage/Examples
Example of how to run your program:
```
go run . data.txt
```
Result:
```
Average: 35
Median: 4
Variance: 5
Standard Deviation: 65
```
## Prepare for Tests

To test the program, the auditor will need to download a zip file provided in this [link](https://github.com/01-edu/public/blob/master/subjects/math-skills/audit/README.md).

1. After download, go to:
C:\Users\ypati\Downloads\

2. Unzip the downloaded file.

3. Go to:
 \stat-bin-dockerized\stat-bin

4. Copy all the files 
5. Paste the copied files to the directory of the user program.

## Test the program

To test the program (in the user's program directory) run the following command:
```
./bin/math-skills
```
or :
```
sudo ./run.sh math-skills
```

It will run the program downloaded earlier. This program is responsible for generating a file called data.txt, which will be used as input for the second program. Also this program will print the calculations needed and based on these calculations the program will be audited. If the calculations match, the user program shall pass.

To compare the calculations, run the following command:
```
go run . data.txt
```

## Key features

## Average

Definition:

 The average, often called the mean, is a measure of central tendency that represents the sum of a set of values divided by the number of values.

## Median

Definition: 

The median is another measure of central tendency that represents the middle value of a dataset when it is ordered from least to greatest. If there is an even number of values, the median is the average of the two middle values.

## Variance 

Definition: 

Variance is a measure of how spread out the values in a dataset are around the mean. It quantifies the degree of variation or dispersion of a set of values.

## Standard Deviation

Definition: 

The standard deviation is the square root of the variance. It provides a measure of the average distance of each data point from the mean and is in the same unit as the data itself.
## Contact

- ychaniot (Ypatios Chaniotakos)

Project [link](https://platform.zone01.gr/git/ychaniot/math-skills).




