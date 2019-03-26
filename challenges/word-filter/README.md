# challenge-go-word-filter

Challenge - Process data from file filtering it to output file.

## Description

Given an input file, the program should filter it from an dictionary (could be an [list of English words](https://github.com/dwyl/english-words/raw/master/words_alpha.txt)) then remove it word and save to output file.


Sample of dictionary ("black list"):

```
marco
zurich
vacations
paper
rock
[..]
```

Sample of input file:

```
zurich > marco > vacations; paper, rock, Rock
hdiashdiuahsh ihi ada shiasduhi aaaaaaa
marco tulio braga mtulio
github@circleci!travisci>>Miolo<<<<
Rica@*px!afinal
Rock, paper, scissors!
Rock, *paper*, scissors!
Rock! paper! scissors!
Rock+paper+scissors?
```

Sample output (using the [list of English words](https://github.com/dwyl/english-words/raw/master/words_alpha.txt) ):

> The order bellow could show different, 'cause the output file handler is writing on file as it is processed.

```
hdiashdiuahsh   shiasduhi aaaaaaa
 >  > ; , , 
 tulio braga mtulio
github@circleci!travisci>>Miolo<<<<
++?
Rica@*px!afinal
, , !
! ! !
, **, !
```

## Resolutions

### Proposal01

* Load all dictionary and save it in memory
* Read each line of input and remove the dictionary words
* Save new line to output

