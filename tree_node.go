package dyntpl

type Node struct {
	typ       Type
	raw       []byte
	rawStatic bool
	rawSsc    ssCache
	prefix    []byte
	suffix    []byte

	ctxVar       []byte
	ctxSrc       []byte
	ctxSrcStatic bool
	ctxSrcSsc    ssCache
	ctxIns       []byte

	condL       []byte
	condLStatic bool
	condLSsc    ssCache
	condOp      Op
	condR       []byte
	condRStatic bool
	condRSsc    ssCache

	loopKey       []byte
	loopVal       []byte
	loopSrc       []byte
	loopSrcSsc    ssCache
	loopCnt       []byte
	loopCntInit   []byte
	loopCntStatic bool
	loopCntSsc    ssCache
	loopCntOp     Op
	loopCondOp    Op
	loopLim       []byte
	loopLimStatic bool
	loopLimSsc    ssCache
	loopSep       []byte

	switchArg    []byte
	switchArgSsc ssCache

	caseL       []byte
	caseLStatic bool
	caseLSsc    ssCache
	caseOp      Op
	caseR       []byte
	caseRStatic bool
	caseRSsc    ssCache

	mod []mod

	child []Node
}

type ssCache []string

func addNode(nodes []Node, node Node) []Node {
	nodes = append(nodes, node)
	return nodes
}

func addRaw(nodes []Node, raw []byte) []Node {
	if len(raw) == 0 {
		return nodes
	}
	node := Node{typ: TypeRaw}
	node.raw, node.rawStatic, node.rawSsc, node.mod = parseRaw(raw, nil)
	nodes = append(nodes, node)
	return nodes
}

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
