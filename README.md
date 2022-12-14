# Memuniq #

**uniq** but with memory, will only output lines that are unique to it.

It uses a bloom filter which means it will never print a line it has seen before. 

Default config is an error rate of 0.1% when 1 million items are added to the filter.  
With this configuration memuniq uses about 5megs of RAM.

## Usage ##
```
Usage of ./memuniq:
  -a	Abort process if the filter file does not exist
  -f string
    	Location of bloomfilter file (default "/home/cpuboi/.cache/bloomfilter.bin")
  -i	Show information about processed lines
  -n	Create a new filter and delete the old
  -p float
    	Approximate error rate percentage, default 0.001% (default 0.001)
  -s int
    	Size of bloomfilter before major collissions occur (default 1000000)
  -v	Show verbose information
```


### Compiling ###
```
go build -ldflags="-s -w" memuniq.go
```


### Performance testing ###
Generate a textfile:  
```
tr -dc "A-Za-z 0-9" < /dev/urandom | fold -w100|head -n 1000000 > ./1mil.txt
cat ./1mil.txt | memuniq -i -v 
```

### Shrinking the binary ###
Install UPX to compress binary even further  
This shrinks size from 1,6MB to 0,6MB   
```
upx memuniq
```

### Thanks ###
Thanks to Everythingme for the Go Bloom code  
github.com/EverythingMe/inbloom  
