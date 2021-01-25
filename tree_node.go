package dyntpl

// Node is a description of template part.
// Every piece of the template, beginning from static text and finishing of complex structures (switch, loop, ...)
// Represents by this type.
type Node struct {
	typ    Type
	raw    []byte
	prefix []byte
	suffix []byte

	ctxVar       []byte
	ctxSrc       []byte
	ctxSrcStatic bool
	ctxIns       []byte

	cntrVar   []byte
	cntrInit  int
	cntrInitF bool
	cntrOp    Op
	cntrOpArg int

	condL       []byte
	condR       []byte
	condStaticL bool
	condStaticR bool
	condOp      Op
	condHlp     []byte
	condHlpArg  []*arg

	loopKey       []byte
	loopVal       []byte
	loopSrc       []byte
	loopCnt       []byte
	loopCntInit   []byte
	loopCntStatic bool
	loopCntOp     Op
	loopCondOp    Op
	loopLim       []byte
	loopLimStatic bool
	loopSep       []byte

	switchArg []byte

	caseL       []byte
	caseR       []byte
	caseStaticL bool
	caseStaticR bool
	caseOp      Op
	caseHlp     []byte
	caseHlpArg  []*arg

	tpl [][]byte

	mod []mod

	child []Node
}

// Add new node to the destination list.
func addNode(nodes []Node, node Node) []Node {
	nodes = append(nodes, node)
	return nodes
}

// Add raw node (static text) to the list.
func addRaw(nodes []Node, raw []byte) []Node {
	if len(raw) == 0 {
		return nodes
	}
	nodes = append(nodes, Node{typ: TypeRaw, raw: raw})
	return nodes
}

// Split nodes by divider node.
func splitNodes(nodes []Node) [][]Node {
	if len(nodes) == 0 {
		return nil
	}
	split := make([][]Node, 0)
	var o int
	for i, node := range nodes {
		if node.typ == TypeDiv {
			split = append(split, nodes[o:i])
			o = i + 1
		}
	}
	if o < len(nodes) {
		split = append(split, nodes[o:])
	}
	return split
}

// Walk over the nodes list and group them by the type, need to make tree of switch structure.
func rollupSwitchNodes(nodes []Node) []Node {
	if len(nodes) == 0 {
		return nil
	}
	var (
		r     = make([]Node, 0)
		group = Node{typ: -1}
	)
	for _, node := range nodes {
		if node.typ != TypeCase && node.typ != TypeDefault && group.typ == -1 {
			continue
		}
		if node.typ == TypeCase || node.typ == TypeDefault {
			if group.typ != -1 {
				r = append(r, group)
			}
			group = node
			continue
		}
		group.child = append(group.child, node)
	}
	if len(group.child) > 0 {
		r = append(r, group)
	}
	return r
}
