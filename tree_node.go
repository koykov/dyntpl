package dyntpl

// node is a description of template part.
// Every piece of the template, beginning from static text and finishing of complex structures (switch, loop, ...)
// Represents by this type.
type node struct {
	typ    rtype
	raw    []byte
	prefix []byte
	suffix []byte
	noesc  bool

	ctxVar       []byte
	ctxSrc       []byte
	ctxOK        []byte
	ctxSrcStatic bool
	ctxIns       []byte

	cntrVar   []byte
	cntrInit  int
	cntrInitF bool
	cntrOp    op
	cntrOpArg int

	condL, condOKL []byte
	condR, condOKR []byte
	condStaticL    bool
	condStaticR    bool
	condOp         op
	condHlp        []byte
	condHlpArg     []*arg
	condIns        []byte
	condLC         lc

	loopKey       []byte
	loopVal       []byte
	loopSrc       []byte
	loopCnt       []byte
	loopCntInit   []byte
	loopCntStatic bool
	loopCntOp     op
	loopCondOp    op
	loopLim       []byte
	loopLimStatic bool
	loopSep       []byte
	loopBrkD      int

	switchArg []byte

	caseL       []byte
	caseR       []byte
	caseStaticL bool
	caseStaticR bool
	caseOp      op
	caseHlp     []byte
	caseHlpArg  []*arg

	tpl [][]byte
	loc []byte

	mod   []mod
	child []node
}

// Add new node to the destination list.
func addNode(nodes []node, node node) []node {
	nodes = append(nodes, node)
	return nodes
}

// Add raw node (static text) to the list.
func addRaw(nodes []node, raw []byte) []node {
	if len(raw) == 0 {
		return nodes
	}
	nodes = append(nodes, node{typ: typeRaw, raw: raw})
	return nodes
}

// Split nodes by divider node.
func splitNodes(nodes []node) [][]node {
	if len(nodes) == 0 {
		return nil
	}
	split := make([][]node, 0)
	var o int
	for i, n := range nodes {
		if n.typ == typeDiv {
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
func rollupSwitchNodes(nodes []node) []node {
	if len(nodes) == 0 {
		return nil
	}
	var (
		r     = make([]node, 0)
		group = node{typ: -1}
	)
	for _, n := range nodes {
		if n.typ != typeCase && n.typ != typeDefault && group.typ == -1 {
			continue
		}
		if n.typ == typeCase || n.typ == typeDefault {
			if group.typ != -1 {
				r = append(r, group)
			}
			group = n
			continue
		}
		group.child = append(group.child, n)
	}
	if len(group.child) > 0 {
		r = append(r, group)
	}
	return r
}
