# Memuniq #

**uniq** but with memory, will only output lines that are unique to it. Handy when looping through same directories over and over.

It uses a bloom filter which means it will never print a line it has seen before. 

Default config is an error rate of 0.1% when 1 million items are added to the filter.  
With this configuration memuniq uses about 5megs of RAM.

### Compiling ###
```
go build -o memuniq -ldflags="-s -w" main.go
```
### Shrinking the binary ###
Install UPX to compress binary even further  
This shrinks size from 1,6MB to 0,6MB   
```
upx memuniq
```

### Performance testing ###
Generate a textfile:  
```
tr -dc "A-Za-z 0-9" < /dev/urandom | fold -w100|head -n 1000000 > ./1mil.txt
cat ./1mil.txt | memuniq 
```

