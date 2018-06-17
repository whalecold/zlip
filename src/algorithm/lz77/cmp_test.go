package lz77

import (
	"testing"
	"fmt"
)

func TestLz77Cmp(t *testing.T) {
	//sre := []byte("112312311123")
	//sre := []byte("mnoabczxyabcuvabczxydefghuingabczx4abcabc")
	sre := []byte("Also known as advanced fee fraud (AFF), 4-1-9 scams`" +
		" are named after the section of the Nigerian penal code that deals with `" +
			"fraud. Although originally originating in Nigeria these scams can originate `" +
				"from anywhere. If you fall for one of these at best you will lose thousands of dollars; `" +
					"at worst you will lose your life. These usually start with an email from a bank official or `" +
						"the relative of a recently deceased African president or a government minister informing you that `" +
							"they have access to millions of dollars but need your help to get the money out of the country. The `" +
								"end result is that when the deal is threatened you will be asked for money to secure the release of `" +
									"the funds. Do not under any circumstances reply to these letters, people have been murdered while `" +
										"following up with these scams.")
	//sre := []byte("这棵树也称为码树，其实就是码表的一种形式化描述，每个节点（除了叶子节点）都会分出两个分支，左分支代表比特0，右边分支代表1，`" +
	//	"	从根节点到叶子节点的一个比特序列就是码字。因为只有叶子节点可以是码字，所以这样也符合一个码字不能是另一个码字的前缀这一原则。`" +
	//	"说到这里，可以说一下另一个话题，就是一个映射表map在内存中怎么存储，没有相关经验的可以跳过，map实现的是key-->value这样的一个表，`" +
	//	"map的存储一般有哈希表和树形存储两类，树形存储就可以采用上面这棵树，树的中间节点并没有什么含义，叶子节点的值表示value，从根节点到叶子节`" +
	//	"点上分支的值就是key，这样比较适合存储那些key由多个不等长字符组成的场合，比如key如果是字符串，那么把二叉树的分支扩展很多，成为多叉树，`" +
	//	"每个分支就是a,b,c,d这种字符，这棵树也就是Trie树，是一种很好使的数据结构。利用树的遍历算法，就实现了一个有序Map。好了`" +
	//	"，我们理解了Huffman编码的思想，我们来看看distance的实际情况。ZIP中滑动窗口大小固定为32KB，也就是说，distance的值范围是1-32768。那么，通过上面`" +
	//	"的方式，统计频率后，就得到32768个码字，按照上面这种方式可以构建出来。于是我们会遇到一个最大的问题，那就是这棵树太大了，怎么记录呢？好了，个人认为到`" +
	//	"了ZIP的核心了，那就是码树应该怎么缩小，以及码树怎么记录的问题。")
	//sre := []byte("As mentioned above,there are many kinds of
	// wireless systems other than cellular.")
	//sre := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	//fmt.Printf("result : %v\n", string(sre))
	fmt.Printf("result : %v\n", len(sre))
	re := Lz77Compress(sre, uint64(len(sre)))
	//fmt.Printf("result : %v\n", string(re))
	fmt.Printf("result : %v\n", len(re))

	un := UnLz77Compress(re)
	fmt.Printf("--- %v\n", string(un))
	//lz77Cmp([]byte("1322232111231"), int64(len([]byte("1322232111231"))))
}
