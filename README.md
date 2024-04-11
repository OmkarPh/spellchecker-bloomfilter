
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


<!-- ## Supported commands

Detailed documentation - https://redis.io/commands/

| Command  | Syntax                                   | Example                                                   | Description                                     |
|----------|------------------------------------------|-----------------------------------------------------------|-------------------------------------------------|
| SET      | **SET key value** [NX / XX] [GET]<br/>[EX seconds / PX milliseconds<br/> / EXAT unix-time-seconds / PXAT unix-time-milliseconds / KEEPTTL]                        | redis-cli SET name omkar<br/>redis-cli SET name omkar GET KEEPTTL              | Set the string value of a key                   |
| GET      | **GET key**                                  | redis-cli GET name                                        | Get the value of a key                          |
| DEL      | **DEL key** [key ...]                        | redis-cli DEL name<br/>redis-cli DEL name age             | Delete one or more keys                         |
| INCR     | **INCR key**                                 | redis-cli INCR age                                        | Increment the integer value of a key            |
| DECR     | **DECR key**                                 | redis-cli DECR age                                        | Decrement the integer value of a key            |
| EXISTS   | **EXISTS key** [key ...]                     | redis-cli EXISTS name<br/>redis-cli EXISTS name age       | Check if a key exists                           |
| EXPIRE   | **EXPIRE key seconds** [NX / XX / GT / LT]   | redis-cli EXPIRE name 20<br/>redis-cli EXPIRE name 20 NX  | Set a key's time to live in seconds             |
| PERSIST  | **PERSIST key **                             | redis-cli PERSIST name                                    | Remove the expiration from a key                |
| TTL      | **TTL key**                                  | redis-cli TTL key                                         | Get the time to live for a key (in seconds)     |
| TYPE     | **TYPE key**                                 | redis-cli TYPE name                                       | Determine the type stored at a key              |
| PING     | **PING**                                     | redis-cli PING                                            | Ping the server                                 |
| ECHO     | **ECHO message**                           | redis-cli ECHO "Hello world"                              | Echo the given string                           | -->




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
go build -o build/spellchecker-bloomfilter -v

./build/spellchecker-bloomfilter
```

## Output
