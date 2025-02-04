package contentfilter

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

// SensitiveFilter 敏感词过滤器
type SensitiveFilter struct {
	root *TrieNode
}

func NewSensitiveFilter() *SensitiveFilter {
	return &SensitiveFilter{
		root: newTrieNode(),
	}
}

func (sf *SensitiveFilter) LoadFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("读取文件出错: %w", err)
		}
		word := strings.TrimSpace(line)
		if word != "" {
			sf.AddKeyword(word)
		}
	}
	return nil
}

func (sf *SensitiveFilter) AddKeyword(keyword string) {
	node := sf.root
	for _, char := range keyword {
		subNode := node.getSubNode(char)
		if subNode == nil {
			subNode = newTrieNode()
			node.addSubNode(char, subNode)
		}
		node = subNode
	}
	node.isKeywordEnd = true
}
func (sf *SensitiveFilter) Filter(text string) string {
	if text == "" {
		return text
	}

	var (
		result   strings.Builder
		start    int
		position int
		node     = sf.root
	)

	for position < len(text) {
		char, charWidth := utf8.DecodeRuneInString(text[position:])
		if char == utf8.RuneError {
			// 如果字符解码失败，直接写入原始字节
			result.WriteByte(text[start])
			position++
			start = position
			node = sf.root
			continue
		}

		// 跳过符号
		if isSymbol(char) {
			// 如果当前节点是根节点，将符号加入结果，移动起始位置
			if node == sf.root {
				result.WriteRune(char)
				start += charWidth
			}
			// 移动当前检查位置
			position += charWidth
			continue
		}

		// 检查下级节点
		node = node.getSubNode(char)
		if node == nil {
			// 当前字符不在敏感词中，将起始字符加入结果
			result.WriteString(text[start : start+charWidth])
			// 移动起始位置和当前检查位置
			position = start + charWidth
			start = position
			// 重置节点到根节点
			node = sf.root
		} else if node.isKeywordEnd {
			// 发现敏感词，替换为 ***
			result.WriteString("***")
			// 移动起始位置和当前检查位置
			position += charWidth
			start = position
			// 重置节点到根节点
			node = sf.root
		} else {
			// 继续检查下一个字符
			position += charWidth
		}
	}

	// 将最后一批字符加入结果
	result.WriteString(text[start:])
	return result.String()
}
func isSymbol(char rune) bool {
	return !unicode.IsLetter(char) && !unicode.IsNumber(char) && (char < 0x2E80 || char > 0x9FFF)
}

type TrieNode struct {
	isKeywordEnd bool
	subNodes     map[rune]*TrieNode
}

func newTrieNode() *TrieNode {
	return &TrieNode{
		subNodes: make(map[rune]*TrieNode),
	}
}

func (tn *TrieNode) getSubNode(char rune) *TrieNode {
	return tn.subNodes[char]
}

func (tn *TrieNode) addSubNode(char rune, node *TrieNode) {
	tn.subNodes[char] = node
}
func SensitiveFilterFun(content string) string {
	// 创建敏感词过滤器
	filter := NewSensitiveFilter()
	// 加载敏感词文件
	err := filter.LoadFromFile("utils/contentfilter/sensitive-words.txt")
	if err != nil {
		log.Println("加载敏感词文件失败:", err)
		return content
	}
	filteredText := filter.Filter(content)

	return filteredText
}

//func TestA(t *testing.T) {
//	// 创建敏感词过滤器
//	filter := NewSensitiveFilter()
//
//	// 加载敏感词文件
//	err := filter.LoadFromFile("sensitive-words.txt")
//	if err != nil {
//		fmt.Println("加载敏感词文件失败:", err)
//		return
//	}
//
//	// 测试过滤
//	text := "开票hshvse赌博"
//	filteredText := filter.Filter(text)
//	fmt.Println("过滤前:", text)
//	fmt.Println("过滤后:", filteredText)
//}
