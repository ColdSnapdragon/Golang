package main

import "fmt"

func main() {
	// prereqs记录了每个课程的前置课程
	var prereqs = map[string][]string{
		"algorithms": {"data structures"},
		"calculus":   {"linear algebra"},
		"compilers": {
			"data structures",
			"formal languages",
			"computer organization",
		},
		"data structures":       {"discrete math"},
		"databases":             {"data structures"},
		"discrete math":         {"intro to programming"},
		"formal languages":      {"discrete math"},
		"networks":              {"operating systems"},
		"operating systems":     {"data structures", "computer organization"},
		"programming languages": {"data structures", "computer organization"},
	}

	var res = topo(prereqs)
	for i, name := range res {
		fmt.Printf("%d:\t%s\n",i,name)
	}
}

func topo(m map[string][]string) []string {
	var order = []string{}

	var vis = map[string]bool{}

	var dfs func([]string) // 只有先声明，匿名函数才能带自我递归
	dfs = func(vexs []string) {
		for _, v := range vexs {
			if !vis[v] {
				vis[v] = true
				dfs(m[v])
				order = append(order, v)
			}
		}
	}

	var st = []string{}
	for k := range m {
		if len(m[k]) < 1 {
			st = append(st, k)
		}
	}
	dfs(st)

	return order
}