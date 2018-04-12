Aho-Corasick is a multiple string matching algorithm
I implement the algorithm in trie.go
In filter.go, I use the built trie to search sensitive words,and filter them out
In test1.go, I test the filter function
In test2.go, I implement a simple http server, and the dictionary can be asynchronously hot-updated and hot-reloading using command: 'kill -1 pid'

go run test1.go
go run test2.go

I tested for a dictionary file which contains 140W lines of diffrent sensitive phrases and 2100W charactors and totally 35M in size. It takes 50s to build the according trie, and it takes 1.5s to filter the all phrases out from a given text file which is the same to the dictionary file ( the worst case ) in a routine. And the total memory usage of the process is 2.4G
