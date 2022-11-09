# memuniq
**uniq** but with memory, will only output lines that are unique to it. Handy when looping through same directories over and over.

It uses a bloom filter which means it will never print a line it has seen before. 
