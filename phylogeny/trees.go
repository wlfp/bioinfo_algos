package phylogeny

import (
	"fmt"
	"slices"
	"strings"
)

func (distanceMatrix distanceMatrix) String() string {
	var builder strings.Builder
	for index, cluster := range distanceMatrix.clusters {
		if slices.Contains(distanceMatrix.deadClusterIndices, index) {
			continue
		}
		if cluster.node == nil {
			continue
		}
		builder.WriteByte('\n')
		builder.WriteString(cluster.node.String())
	}
	return builder.String()
}

type treeNode struct {
	name         string
	left, right  *treeNode
	branchLength float64
	height       float64
}

func (node *treeNode) String() string {
	if node == nil {
		return ""
	}
	var builder strings.Builder
	node.buildPrettyString(&builder, "", true, false)
	return strings.TrimSuffix(builder.String(), "\n")
}

func (node *treeNode) label() string {
	if node == nil {
		return ""
	}
	return node.name + node.metricsSuffix()
}

func (node *treeNode) children() []*treeNode {
	var result []*treeNode
	if node.left != nil {
		result = append(result, node.left)
	}
	if node.right != nil {
		result = append(result, node.right)
	}
	return result
}

func (node *treeNode) metricsSuffix() string {
	if node == nil {
		return ""
	}
	var parts []string
	if node.height != 0 {
		parts = append(parts, fmt.Sprintf("h=%.2f", node.height))
	}
	if node.branchLength != 0 {
		parts = append(parts, fmt.Sprintf("len=%.2f", node.branchLength))
	}
	if len(parts) == 0 {
		return ""
	}
	return " [" + strings.Join(parts, ", ") + "]"
}

func (node *treeNode) buildPrettyString(builder *strings.Builder, prefix string, isTail, hasParent bool) {
	builder.WriteString(prefix)
	if hasParent {
		if isTail {
			builder.WriteString("\\-- ")
		} else {
			builder.WriteString("+-- ")
		}
	}
	builder.WriteString(node.label())
	builder.WriteByte('\n')

	children := node.children()
	for index, child := range children {
		childPrefix := prefix
		if hasParent {
			if isTail {
				childPrefix += "    "
			} else {
				childPrefix += "|   "
			}
		}
		child.buildPrettyString(builder, childPrefix, index == len(children)-1, true)
	}
}
