# challenge-go-word-filter

Challenge - Process data from file filtering it to output file.

## Description

Given an input file, the program should filter it from an dictionary (could be an list of English words) then remove it word and save to output file.


Sample of dictionary (black list):

```
paper
```

Sample of input file:

```
Rock, paper, scissors!
Rock, *paper*, scissors!
Rock! paper! scissors!
Rock+paper+scissors?
```

Sample output:

```
Rock, , scissors!
Rock, **, scissors!
Rock! ! scissors!
Rock++scissors?
```

## Resolutions

### Proposal01

* Load all dictionary and save it in memory
* Read each line of input and remove the dictionary words
* Save new line to output

