package gee

import "strings"

type node struct {
	pattern  string  // 待匹配路由 例如: /get/:id
	part     string  // 路由中的一部分 例如: :id
	children []*node //子节点 例如: [:id, name, *]
	isWild   bool    // 是否精准匹配 part含有: 或 *时为true
}

// 第一个匹配成功的节点, 用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}

	return nil
}

// 所有匹配成功的节点, 用于查询
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 获取路由part节点名称
	// 递归查找每一层的节点，如果没有匹配到当前part的节点，则新建一个
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		// 将查找到的child添加到当前节点的子节点中
		n.children = append(n.children, child)
	}
	// 递归插入子节点
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	// 匹配节点相同的children
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
