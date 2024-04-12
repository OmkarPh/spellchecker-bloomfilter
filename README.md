
## A Bloomfilter powered Spell checker implemented in golang

Inspired from [Coding Challenges - Build Your Own Spell Checker Using A Bloom Filter
](https://codingchallenges.fyi/challenges/challenge-bloom/)

## Prerequisites

- Go 1.22 or later

<!-- ## Try out

```bash
# Using brew (for mac or linux)
brew install omkarph/tap/spellchecker-bloomfilter
redis-server-lite

# or

# Using a release archive from https://github.com/OmkarPh/spellchecker-bloomfilter/releases/latest
cd ~/Downloads # Go to the Downloads folder of your machine
mkdir spellchecker-bloomfilter
tar -xf "release_archive_file" -C spellchecker-bloomfilter
cd spellchecker-bloomfilter
./spellchecker-bloomfilter

# Test
``` -->

## Features

- Uses `.opbf` file to store Bloomfilter data
- Case sensitivity options


## Usage

- Check words `go run . coding challeges are fun to buidl`
- Seed your own dictionary
  
  `go run . -dict ~/projects/spellchecker-bloomfilter/dict.txt`

- Rebuild bloom filter `-build`
  - Customise number of hash functions & size of bloom filter.
     Note - Size refers to number of bits used on disk for filter.
    
    `go run . -build -hashes 3 -size 5200000`

  - With own dictionary

    `go run . -build -dict data/dict.txt -hashes 3 -size 5200000` (With custom dictionary)


## Local setup

```bash
# Clone this repository
git clone https://github.com/OmkarPh/spellchecker-bloomfilter.git
cd spellchecker-bloomfilter

# Install dependencies
go get .

# Run the server
go run .

# Build release executible
go build -o build/spellchecker-bloomfilter

./build/spellchecker-bloomfilter coding challeges are fun to buidl
```

## Output
