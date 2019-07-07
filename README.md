## Interpretation
There was some, in my opinion, contradiction in the description; it has a few references to REST API, but at the same time it looked clear to me that most of the other requirements were for a command line program (which is, of course, not REST).

It would also made little sense to me to write an REST API interface to query another API interface in this case (although I could think of scenarios were such a thing may be sensible).

So I decided to ignore the REST references and work on a command line utility.

## Build
It's a golang program. 
For simplicity I avoided external dependencies, so it should enough to run the following command from the directory with the source code:
 `go build -o [binary]`
 
 The usual [go build options](https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies) apply.
 
## Usage
The program works with MAC addresses passed as options or from a pipe.

An API key to connect to the `macaddress.io` is required; it can be set through the `--key` argument or by a `API_KEY` environment variable.

By default the program would print (when a matching is found) the MAC address and the corresponding company name separated by a `=`, like:
`44:38:39:ff:ef:57=Cumulus Networks, Inc`.

It should be easy to split afterwards, but it's also possible to print only the company name by using the `--name-only` argument.
It's also possible to define a different separator with the `--separator` argument (the empty separator would result in no MAC address printed).

Some examples (using the `API_KEY` environment variable):

    ./gmac 44:38:39:ff:ef:57
    44:38:39:ff:ef:57=Cumulus Networks, Inc
    
    ./gmac --separator "->" 44:38:39:ff:ef:57
    44:38:39:ff:ef:57->Cumulus Networks, Inc
    
    ./gmac --separator "" 44:38:39:ff:ef:57 
    Cumulus Networks, Inc
    
    ./gmac  44:38:39:ff:ef:57 02:42:d6:36:c1:72 88:d7:f6:c7:bd:25  
    44:38:39:ff:ef:57=Cumulus Networks, Inc
    02:42:d6:36:c1:72=
    88:d7:f6:c7:bd:25=ASUSTek Computer Inc
    
    ./gmac  -name-only 44:38:39:ff:ef:57 02:42:d6:36:c1:72 88:d7:f6:c7:bd:25
    Cumulus Networks, Inc
    ASUSTek Computer Inc
    
    ./gmac  -name-only < ./macs.txt
    Cumulus Networks, Inc
    ASUSTek Computer Inc
    
    ./gmac  --separator "->" < ./macs.txt 
    44:38:39:ff:ef:57->Cumulus Networks, Inc
    02:42:d6:36:c1:72->
    88:d7:f6:c7:bd:25->ASUSTek Computer Inc

Errors are sent to `stderr`, while output to `stdout`.

Each given MAC address generates a call to the API; I did not implement some sort of caching (although I thought about it), as you can always use a `sort -u` or the like in the pipe (where duplications are more likely to appear).

## Docker 
There is, as requested, a Docker file that can be use for running the program in a container. 
It's a very simple one and it works, but I did not spend too much time on it, as, it does not seem an optimal solution to me compared to just using compiling the source code directly.

Some examples (using the `--env-file` option for the API key)
    
    docker run --env-file ./api_key_env gmac 44:38:39:ff:ef:57
    44:38:39:ff:ef:57=Cumulus Networks, Inc
    
    docker run --env-file ./api_key_env gmac --name-only 44:38:39:ff:ef:57
    Cumulus Networks, Inc
    
    docker run -i --env-file ./api_key_env gmac --separator "==>" < ./macs.txt
    44:38:39:ff:ef:57==>Cumulus Networks, Inc
    02:42:d6:36:c1:72==>
    88:d7:f6:c7:bd:25==>ASUSTek Computer Inc

## Security
From the program point of view the only case I could think of is the leaking of the API key. 
Initially I had the possibility to read it from a file (thus relying on the file-system for access), but I removed it as setting the env variable through sourcing a file may be enough (YMMV):

    source api_key_env && ./gmac 44:38:39:ff:ef:57
    44:38:39:ff:ef:57=Cumulus Networks, Inc
    
