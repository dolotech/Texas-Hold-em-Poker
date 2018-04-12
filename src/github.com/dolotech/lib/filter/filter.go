package filter

import (
	"bytes"
	_ "fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"
	"github.com/golang/glog"
)

var dicFileTrie map[string]*Trie = make(map[string]*Trie) //map for dictionary file to trie
var dicFileTime map[string]int64 = make(map[string]int64) //map for dictionary file to it's built time

func LoadDicFiles(dicFiles []string) { //load serveral dictionary file to build tries
	for _, dicFile := range dicFiles {
		dicFileTrie[dicFile] = nil
		dicFileTime[dicFile] = 0
	}
	go buildDicFileTrie() //asyn build dictionary file trie,when completed the old trie will be replaced
}

func buildDicFileTrie() { //when server start-up the trie will be built automatically
	for {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP)
		for dicFile, _ := range dicFileTrie {
			stat, e := os.Stat(dicFile)
			if e != nil || stat.ModTime().Unix() > dicFileTime[dicFile] { //maybe deleted file or updated file
				data, e := ioutil.ReadFile(dicFile)
				if e != nil || len(data) <= 0 { //file not exist or empty
					dicFileTrie[dicFile] = nil //delete the old trie,maybe concurrency problem
				}
				dictionary := bytes.Split(data, []byte("\n"))
				var tree Trie
				tree.InitRootNode()
				tree.BuildTrie(dictionary)
				dicFileTrie[dicFile] = &tree             //replace the old trie,maybe concurrency problem
				dicFileTime[dicFile] = time.Now().Unix() //save the replace time
			}
		}
		<-c //after we refresh the all dictionary trie we will block to the next SIGHUP signal and refresh the all dictionary trie again
	}
}

//separator charactors,if we want to match "abcde" in "abc112312de" ,separator charactors can be set to "123"
type Seps []rune

func (s Seps) Len() int {
	return len(s)
}
func (s Seps) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Seps) Less(i, j int) bool {
	return s[i] < s[j]
}

func FindSepC(seps Seps, charactor rune) bool {
	i := sort.Search(len(seps), func(i int) bool { return seps[i] >= charactor })
	return i < len(seps) && seps[i] == charactor
}

// we use a dictionary file to filter text,all separator in text will be bypassed,and matched words will be replace by rep
func FilterText(dicFile string, text []rune, seps Seps, rep rune) {
	T := dicFileTrie[dicFile] //save trie to temporary variable to eliminate concurrency problem

	if T == nil { //no matched trie for dictionary file
		return
	}
	sort.Sort(seps)        //sort seps to speed up search process
	walkNode := T.RootNode //walk through the trie from root

	var nextNode *Node //point to walkNode's next node
	for i, charactor := range text { //travel the text,i is used as an index to present charactor
		for {
			if FindSepC(seps, charactor) { //omit charactor in seps
				break
			}
			nextNode = walkNode.BinGetChildNodeByVal(charactor)
			if nextNode == nil { //not found next node whose val is charactor
				if nil != walkNode.SuffixNode { //point to suffix node, redo the find operation,maybe its suffix node is root
					walkNode = walkNode.SuffixNode
					continue
				} else { //only root node's suffix node is null
					walkNode = T.RootNode //restart from root
					break                 //break to handle next charactor
				}
			} else { // find node
				if nextNode.EOW { //if a word end up with next node
					depth := nextNode.Depth //get the word length
					j := i                  //search back from i
					for depth > 0 { //until to root
						for j >= 0 && FindSepC(seps, text[j]) { //omit
							j--
						}
						if j >= 0 {
							text[j] = rep //replace with rep
							j--
						}
						depth--
					}
					walkNode = T.RootNode //restart from root
				} else { //not EOW
					walkNode = nextNode //step to next node
				}
				break //break to handle next charactor
			}
		}
	}

}




// we use a dictionary file to filter text,all separator in text will be bypassed,and matched words will be replace by rep
func IsInValid(dicFile string, text []rune, seps Seps, rep rune)(valid bool) {
	T := dicFileTrie[dicFile] //save trie to temporary variable to eliminate concurrency problem

	if T == nil { //no matched trie for dictionary file
		glog.Errorf("empty is %+v",T)
		return
	}
	sort.Sort(seps)        //sort seps to speed up search process
	walkNode := T.RootNode //walk through the trie from root

	var nextNode *Node               //point to walkNode's next node
	for i, charactor := range text { //travel the text,i is used as an index to present charactor
		for {
			if FindSepC(seps, charactor) { //omit charactor in seps
				break
			}
			nextNode = walkNode.BinGetChildNodeByVal(charactor)
			if nextNode == nil { //not found next node whose val is charactor
				if nil != walkNode.SuffixNode { //point to suffix node, redo the find operation,maybe its suffix node is root
					walkNode = walkNode.SuffixNode
					continue
				} else { //only root node's suffix node is null
					walkNode = T.RootNode //restart from root
					break                 //break to handle next charactor
				}
			} else { // find node
				if nextNode.EOW { //if a word end up with next node
					depth := nextNode.Depth //get the word length
					j := i                  //search back from i
					for depth > 0 {         //until to root
						for j >= 0 && FindSepC(seps, text[j]) { //omit
							j--
						}
						if j >= 0 {
							valid=true
							//text[j] = rep //replace with rep
							//j--
						}
						depth--
					}
					walkNode = T.RootNode //restart from root
				} else { //not EOW
					walkNode = nextNode //step to next node
				}
				break //break to handle next charactor
			}
		}
	}
	return
}
