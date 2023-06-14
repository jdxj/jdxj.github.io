---
title: "Trie树"
date: 2023-06-13T11:06:55+08:00
---

Trie 树，也叫“字典树”

how，hi，her，hello，so，see组成的trie树

![](https://static001.geekbang.org/resource/image/28/32/280fbc0bfdef8380fcb632af39e84b32.jpg?wh=1142*573)

存储结构

```java
class TrieNode {
  char data;
  TrieNode children[26];
}
```

![](https://static001.geekbang.org/resource/image/f5/35/f5a4a9cb7f0fe9dcfbf29eb1e5da6d35.jpg?wh=1142*697)

```java
public class Trie {
  private TrieNode root = new TrieNode('/'); // 存储无意义字符

  // 往Trie树中插入一个字符串
  public void insert(char[] text) {
    TrieNode p = root;
    for (int i = 0; i < text.length; ++i) {
      int index = text[i] - 'a';
      if (p.children[index] == null) {
        TrieNode newNode = new TrieNode(text[i]);
        p.children[index] = newNode;
      }
      p = p.children[index];
    }
    p.isEndingChar = true;
  }

  // 在Trie树中查找一个字符串
  public boolean find(char[] pattern) {
    TrieNode p = root;
    for (int i = 0; i < pattern.length; ++i) {
      int index = pattern[i] - 'a';
      if (p.children[index] == null) {
        return false; // 不存在pattern
      }
      p = p.children[index];
    }
    if (p.isEndingChar == false) return false; // 不能完全匹配，只是前缀
    else return true; // 找到pattern
  }

  public class TrieNode {
    public char data;
    public TrieNode[] children = new TrieNode[26];
    public boolean isEndingChar = false;
    public TrieNode(char data) {
      this.data = data;
    }
  }
}
```

- 构建 Trie 树的过程，需要扫描所有的字符串，时间复杂度是 O(n)（n 表示所有字符串的长度和）。
- 在其中查找字符串的时间复杂度是 O(k)，k 表示要查找的字符串的长度。

这种用数组存储子节点的方式有点浪费内存

- 假设我们用有序数组，数组中的指针按照所指向的子节点中的字符的大小顺序排列。查询的时候，我们可以通过二分查找的方法，快速查找到某个字符应该匹配的子
  节点的指针。但是，在往 Trie 树中插入一个字符串的时候，我们为了维护数组中数据的有序性，就会稍微慢了点。
- 实际上，Trie 树的变体有很多，都可以在一定程度上解决内存消耗的问题。比如，缩点优化，就是对只有一个子节点的节点，而且此节点不是一个串的结束节点，
  可以将此节点与子节点合并。这样可以节省空间，但却增加了编码难度。

![](https://static001.geekbang.org/resource/image/87/11/874d6870e365ec78f57cd1b9d9fbed11.jpg?wh=1142*581)

# Trie的要求

- 字符串中包含的字符集不能太大。我们前面讲到，如果字符集太大，那存储空间可能就会浪费很多。即便可以优化，但也要付出牺牲查询、插入效率的代价。
- 要求字符串的前缀重合比较多，不然空间消耗会变大很多。
- 如果要用 Trie 树解决问题，那我们就要自己从零开始实现一个 Trie 树，还要保证没有 bug，这个在工程上是将简单问题复杂化，除非必须，一般不建议这样
  做。
- 通过指针串起来的数据块是不连续的，而 Trie 树中用到了指针，所以，对缓存并不友好，性能上会打个折扣。
