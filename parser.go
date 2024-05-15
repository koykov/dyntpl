package dyntpl

import (
	"bytes"
	"fmt"
	"hash/crc64"
	"os"
	"regexp"
	"strconv"

	"github.com/koykov/bytealg"
	"github.com/koykov/byteconv"
)

const (
	// Types of targets.
	targetCond = iota
	targetLoop
	targetSwitch
)

// Parser object.
type Parser struct {
	// Keep format flag. Remove all new lines and tabulations when false.
	keepFmt bool
	// Template body to parse.
	tpl []byte

	// Counters (depths) of conditions, loops and switches.
	cc, cl, cs int
}

// Target is a storage of depths needed to provide proper out from conditions, loops and switches control structures.
type target map[int]int

var (
	// Byte constants.
	empty    []byte
	space    = []byte(" ")
	spaceCBE = []byte("} ")
	comma    = []byte(",")
	uscore   = []byte("_")
	vline    = []byte("|")
	colon    = []byte(":")
	quotes   = []byte("\"'`")
	ddquote  = []byte(`""`)
	// quote      = []byte("\"")
	noFmt      = []byte(" \t\n")
	ctlExit    = []byte("exit")
	ctlOpen    = []byte("{%")
	ctlClose   = []byte("%}")
	ctlTrim    = []byte("{}% ")
	ctlTrimAll = []byte("{}%= ")
	ctxStatic  = []byte("static")
	condElse   = []byte("else")
	condEnd    = []byte("endif")
	condLen    = []byte("len")
	condCap    = []byte("cap")
	loopEnd    = []byte("endfor")
	loopBrk    = []byte("break")
	loopLBrk   = []byte("lazybreak")
	loopCont   = []byte("continue")
	swDefault  = []byte("default")
	swEnd      = []byte("endswitch")
	jq         = []byte("jsonquote")
	jqEnd      = []byte("endjsonquote")
	he         = []byte("htmlescape")
	heEnd      = []byte("endhtmlescape")
	ue         = []byte("urlencode")
	ueEnd      = []byte("endurlencode")
	bTrue      = []byte("true")
	nl         = []byte("\n")
	cr         = []byte("\r")
	crlf       = []byte("\r\n")
	tab        = []byte("\t")
	symEndl    = []byte("endl")
	symNl      = []byte("nl")
	symLf      = []byte("lf")
	symN       = []byte(`\n`)
	symCr      = []byte("cr")
	symR       = []byte(`\r`)
	symCrLn    = []byte("crlf")
	symRN      = []byte(`\r\n`)
	symTab     = []byte("tab")
	symT       = []byte(`\t`)

	// Print prefixes and replacements.
	idJ   = []byte("jsonEscape")
	idQ   = []byte("jsonQuote")
	idH   = []byte("htmlEscape")
	idL   = []byte("linkEscape")
	idU   = []byte("urlEncode")
	idA   = []byte("attrEscape")
	idC   = []byte("cssEscape")
	idJS  = []byte("jsEscape")
	outmf = 'f'                 // float precision floor
	idf   = []byte("floorPrec") // float precision floor
	outmF = 'F'                 // float precision ceil
	idF   = []byte("ceilPrec")  // float precision ceil

	// Operation constants.
	opEq  = []byte("==")
	opNq  = []byte("!=")
	opGt  = []byte(">")
	opGtq = []byte(">=")
	opLt  = []byte("<")
	opLtq = []byte("<=")
	opInc = []byte("++")
	opDec = []byte("--")

	// Regexp to clear template.
	reCutComments = regexp.MustCompile(`{#[^#]*#}`)
	reCutFmt      = regexp.MustCompile(`\n+\t*\s*`)

	// Regexp to parse print instructions.
	reTplPS    = regexp.MustCompile(`^([jhqluacJfF.\d]*)=\s*(.*) (?:prefix|pfx) (.*) (?:suffix|sfx) (.*)`)
	reTplP     = regexp.MustCompile(`^([jhqluacJfF.\d]*)=\s*(.*) (?:prefix|pfx) (.*)`)
	reTplS     = regexp.MustCompile(`^([jhqluacJfF.\d]*)=\s*(.*) (?:suffix|sfx) (.*)`)
	reTpl      = regexp.MustCompile(`^([jhqluacJfF.\d]*)=\s*(.*)`)
	reTplCB    = regexp.MustCompile(`^([^(\s]+)\(([^)]*)\)`)
	reModPfxF  = regexp.MustCompile(`([fF]+)\.*(\d*).*`)
	reModNoVar = regexp.MustCompile(`([^(]+)\(([^)]*)\)`)
	reMod      = regexp.MustCompile(`([^(]+)\(*([^)]*)\)*`)

	// Regexp to parse context instruction.
	reCtxAs  = regexp.MustCompile(`(?:context|ctx) (\w+),*\s*(\w*)\s*=\s*([\w\s.,:|()"'\[\]]+) as ([\[\]\*\w]*)` + "")
	reCtxDot = regexp.MustCompile(`(?:context|ctx) (\w+),*\s*(\w*)\s*=\s*([\w\s.,:|()"'\[\]]+)\.\(([\[\]\*\w]*)\)` + "")
	reCtx    = regexp.MustCompile(`(?:context|ctx) (\w+),*\s*(\w*)\s*=\s*([\w\s.,:|()"'\[\]]+)`)
	reCtxS0  = regexp.MustCompile(`(?:context|ctx) (\w+),*\s*(\w*)\s*=\s*"+([^"]+)"+`)
	reCtxS1  = regexp.MustCompile(`(?:context|ctx) (\w+),*\s*(\w*)\s*=\s*'+([^']+)+'`)

	// Regexp to parse counter instructions.
	reCntr     = regexp.MustCompile(`(?:counter|cntr) (\w+)`)
	reCntrInit = regexp.MustCompile(`(?:counter|cntr) (\w+)\s*=\s*(\d+)`)
	reCntrOp0  = regexp.MustCompile(`(?:counter|cntr) (\w+)(\+\+|--)`)
	reCntrOp1  = regexp.MustCompile(`(?:counter|cntr) (\w+)(\+\d+|-\d+)`)

	// Regexp to parse condition instruction.
	reCond        = regexp.MustCompile(`if .*`)
	reCondExpr    = regexp.MustCompile(`if (.*)(==|!=|>=|<=|>|<)(.*)`)
	reCondHelper  = regexp.MustCompile(`if ([^(]+)\(*([^)]*)\)`)
	reCondComplex = regexp.MustCompile(`if .*&&|\|\||\(|\).*`)
	reCondOK      = regexp.MustCompile(`if (\w+),*\s*(\w*)\s*:*=\s*([^(]+)\(*([^)]*)\)(.*)\s*;\s*([!\w]+)`)
	reCondAsOK    = regexp.MustCompile(`if (\w+),*\s*(\w*)\s*:*=\s*([^(]+)\(*([^)]*)\) as (\w*)\s*;\s*([!\w]+)`)
	reCondDotOK   = regexp.MustCompile(`if (\w+),*\s*(\w*)\s*:*=\s*([^(]+)\(*([^)]*)\)\.\((\w*)\)\s*;\s*([!\w]+)`)
	reCondExprOK  = regexp.MustCompile(`if .*;\s*([!:\w]+)(.*)(.*)`)

	// Regexp to parse loop instruction.
	reLoop      = regexp.MustCompile(`for .*`)
	reLoopRange = regexp.MustCompile(`for ([^:]+)\s*:*=\s*range\s*([^\s]*)\s*(?:separator|sep)*\s*(.*)` + "")
	reLoopCount = regexp.MustCompile(`for (\w*)\s*:*=\s*(\w+)\s*;\s*\w+\s*(<|<=|>|>=|!=)+\s*([^;]+)\s*;\s*\w*(--|\+\+)+\s*(?:separator|sep)*\s*(.*)`)
	reLoopBrk   = regexp.MustCompile(`break (\d+)`)
	reLoopLBrk  = regexp.MustCompile(`lazybreak (\d+)`)

	// Regexp to parse switch instruction.
	reSwitch           = regexp.MustCompile(`^switch\s*(.*)`)
	reSwitchCase       = regexp.MustCompile(`case ([^<=>!]+)([<=>!]{2})*(.*)`)
	reSwitchCaseHelper = regexp.MustCompile(`case ([^(]+)\(*([^)]*)\)`)

	// Regexp to parse include instruction.
	reInc = regexp.MustCompile(`(?:include|\.) (.*)`)

	crc64Tab = crc64.MakeTable(crc64.ISO)

	// Suppress go vet warning.
	_ = ParseFile
)

// Parse initializes parser and parse the template body.
func Parse(tpl []byte, keepFmt bool) (tree *Tree, err error) {
	p := &Parser{
		tpl:     tpl,
		keepFmt: keepFmt,
	}
	p.cutComments()
	p.cutFmt()

	hsum := crc64.Checksum(p.tpl, crc64Tab)
	if tree = tplDB.getTreeByHash(hsum); tree != nil {
		return
	}

	// Prepare template tree.
	tree = &Tree{hsum: hsum}
	target := newTarget(p)
	tree.nodes, _, err = p.parseTpl(tree.nodes, 0, target)
	return
}

// ParseFile initializes parser and parse file contents.
func ParseFile(fileName string, keepFmt bool) (tree *Tree, err error) {
	_, err = os.Stat(fileName)
	if os.IsNotExist(err) {
		return
	}
	var raw []byte
	raw, err = os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file %s", fileName)
	}
	return Parse(raw, keepFmt)
}

// Remove all comments from the template body.
func (p *Parser) cutComments() {
	p.tpl = reCutComments.ReplaceAll(p.tpl, empty)
}

// Remove template formatting if needed.
func (p *Parser) cutFmt() {
	if p.keepFmt {
		return
	}
	p.tpl = reCutFmt.ReplaceAll(p.tpl, empty)
	p.tpl = bytealg.Trim(p.tpl, noFmt)
}

// Initial parsing method.
func (p *Parser) parseTpl(nodes []Node, offset int, target *target) ([]Node, int, error) {
	var (
		up  bool
		err error
	)

	// Walk over template body and find control structures.
	o, i := offset, offset
	inCtl := false
	for !target.reached(p) || target.eqZero() {
		i = bytealg.IndexAt(p.tpl, ctlOpen, i)
		if i < 0 {
			if inCtl {
				return nodes, o, ErrUnexpectedEOF
			}
			nodes = addRaw(nodes, p.tpl[o:])
			o = len(p.tpl)
			break
		}
		if inCtl {
			// We are inside control structure.
			e := bytealg.IndexAt(p.tpl, ctlClose, i)
			if e < 0 {
				return nodes, o, ErrUnexpectedEOF
			}
			e += 2
			node := Node{}
			nodes, e, up, err = p.processCtl(nodes, &node, p.tpl[o:e], o)
			if err != nil {
				return nodes, o, err
			}
			o, i = e, e
			inCtl = false
			if up {
				break
			}
		} else {
			// Start of control structure caught.
			nodes = addRaw(nodes, p.tpl[o:i])
			o = i
			inCtl = true
		}
	}
	return nodes, o, nil
}

// General parsing method.
func (p *Parser) processCtl(nodes []Node, root *Node, ctl []byte, pos int) ([]Node, int, bool, error) {
	var (
		offset int
		up     bool
		err    error
	)

	up = false
	t := bytealg.Trim(ctl, ctlTrim)
	// Check tpl (print) structure.
	if reTplPS.Match(t) || reTplP.Match(t) || reTplS.Match(t) || reTpl.Match(t) || reTplCB.Match(t) {
		// Sequentially check print structure from the complex to the simplest.
		root.typ = TypeTpl
		if m := reTplPS.FindSubmatch(t); m != nil {
			// Tpl with prefix and suffix found.
			root.raw, root.mod = p.extractMods(m[2], m[1])
			root.prefix = m[3]
			root.suffix = m[4]
		} else if m := reTplP.FindSubmatch(t); m != nil {
			// Tpl with prefix found.
			root.raw, root.mod = p.extractMods(m[2], m[1])
			root.prefix = m[3]
		} else if m := reTplS.FindSubmatch(t); m != nil {
			// Tpl with suffix found.
			root.raw, root.mod = p.extractMods(m[2], m[1])
			root.suffix = m[3]
		} else if m := reTpl.FindSubmatch(t); m != nil {
			// Simple tpl found.
			root.raw, root.mod = p.extractMods(bytealg.Trim(m[2], ctlTrimAll), m[1])
		} else if m := reTplCB.FindSubmatch(t); m != nil {
			root.raw, root.mod = p.extractMods(bytealg.Trim(m[0], ctlTrimAll), m[1])
		} else {
			root.raw, root.mod = p.extractMods(bytealg.Trim(t, ctlTrimAll), nil)
		}
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check context structure.
	if reCtx.Match(t) {
		root.typ = TypeCtx
		var (
			m           [][]byte
			forceStatic bool
		)
		m = reCtxAs.FindSubmatch(t)
		if m == nil {
			m = reCtxDot.FindSubmatch(t)
		}
		if m == nil {
			m = reCtxS0.FindSubmatch(t)
			forceStatic = m != nil
		}
		if m == nil {
			m = reCtxS1.FindSubmatch(t)
			forceStatic = m != nil
		}
		if m == nil {
			m = reCtx.FindSubmatch(t)
		}
		root.ctxVar, root.ctxOK = m[1], m[2]
		root.ctxSrc, root.mod = p.extractMods(m[3], nil)
		root.ctxSrcStatic = isStatic(root.ctxSrc) || forceStatic
		if len(m) > 4 && len(m[4]) > 0 {
			root.ctxIns = m[4]
		} else {
			if _, ok := GetInsByVarName(byteconv.B2S(root.ctxVar)); !ok {
				root.ctxIns = ctxStatic
			}
		}
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check counter structure.
	if reCntr.Match(t) {
		root.typ = TypeCounter
		root.cntrInitF = false
		if m := reCntrInit.FindSubmatch(t); m != nil {
			root.cntrVar = m[1]
			root.cntrInitF = true
			i, err := strconv.Atoi(string(m[2]))
			if err != nil {
				return nodes, offset, up, err
			}
			root.cntrInit = i
		} else if m := reCntrOp0.FindSubmatch(t); m != nil {
			root.cntrVar = m[1]
			if bytes.Equal(m[2], opDec) {
				root.cntrOp = OpDec
			} else {
				root.cntrOp = OpInc
			}
			root.cntrOpArg = 1
		} else if m := reCntrOp1.FindSubmatch(t); m != nil {
			root.cntrVar = m[1]
			if m[2][0] == '-' {
				root.cntrOp = OpDec
			} else {
				root.cntrOp = OpInc
			}
			a, err := strconv.Atoi(string(m[2][1:]))
			if err != nil {
				return nodes, offset, up, err
			}
			root.cntrOpArg = a
		}
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check if-ok instruction.
	if reCondOK.Match(t) {
		root.typ = TypeCondOK
		var m [][]byte
		m = reCondAsOK.FindSubmatch(t)
		if m == nil {
			m = reCondDotOK.FindSubmatch(t)
		}
		if m == nil {
			m = reCondOK.FindSubmatch(t)
		}
		root.condOKL, root.condOKR = m[1], m[2]
		root.condHlp, root.condHlpArg = m[3], p.extractArgs(m[4])
		if len(m[5]) > 0 {
			root.condIns = m[5]
		}
		root.condL, root.condR, root.condStaticL, root.condStaticR, root.condOp = p.parseCondExpr(reCondExprOK, t)

		target := newTarget(p)
		p.cc++

		subNodes := make([]Node, 0)
		subNodes, offset, err = p.parseTpl(subNodes, pos+len(ctl), target)
		split := splitNodes(subNodes)
		if len(split) > 0 {
			nodeTrue := Node{typ: TypeCondTrue, child: split[0]}
			root.child = append(root.child, nodeTrue)
		}
		if len(split) > 1 {
			nodeFalse := Node{typ: TypeCondFalse, child: split[1]}
			root.child = append(root.child, nodeFalse)
		}

		nodes = addNode(nodes, *root)
		return nodes, offset, up, err
	}

	// Check condition structure.
	if reCond.Match(t) {
		// Check complexity of the condition first.
		if reCondComplex.Match(t) {
			// Check if condition may be handled by the condition helper.
			if m := reCondHelper.FindSubmatch(t); m != nil {
				target := newTarget(p)
				p.cc++

				subNodes := make([]Node, 0)
				subNodes, offset, err = p.parseTpl(subNodes, pos+len(ctl), target)
				split := splitNodes(subNodes)

				root.typ = TypeCond
				root.condHlp = m[1]
				root.condHlpArg = p.extractArgs(m[2])
				root.condL, root.condR, root.condStaticL, root.condStaticR, root.condOp = p.parseCondExpr(reCondExpr, t)
				switch {
				case bytes.Equal(root.condHlp, condLen):
					root.condLC = lcLen
				case bytes.Equal(root.condHlp, condCap):
					root.condLC = lcCap
				}
				if len(split) > 0 {
					nodeTrue := Node{typ: TypeCondTrue, child: split[0]}
					root.child = append(root.child, nodeTrue)
				}
				if len(split) > 1 {
					nodeFalse := Node{typ: TypeCondFalse, child: split[1]}
					root.child = append(root.child, nodeFalse)
				}

				nodes = addNode(nodes, *root)
				return nodes, offset, up, err
			}
			return nodes, pos, up, fmt.Errorf("too complex condition '%s' at offset %d", t, pos)
		}
		// Create new target, increase condition counter and dive deeper.
		target := newTarget(p)
		p.cc++

		subNodes := make([]Node, 0)
		subNodes, offset, err = p.parseTpl(subNodes, pos+len(ctl), target)
		split := splitNodes(subNodes)

		root.typ = TypeCond
		root.condL, root.condR, root.condStaticL, root.condStaticR, root.condOp = p.parseCondExpr(reCondExpr, t)
		if len(split) > 0 {
			nodeTrue := Node{typ: TypeCondTrue, child: split[0]}
			root.child = append(root.child, nodeTrue)
		}
		if len(split) > 1 {
			nodeFalse := Node{typ: TypeCondFalse, child: split[1]}
			root.child = append(root.child, nodeFalse)
		}

		nodes = addNode(nodes, *root)
		return nodes, offset, up, err
	}
	// Check condition divider.
	if bytes.Equal(t, condElse) {
		root.typ = TypeDiv
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check condition end.
	if bytes.Equal(t, condEnd) {
		// End of condition caught. Decrease the counter and exit.
		p.cc--
		offset = pos + len(ctl)
		up = true
		return nodes, offset, up, err
	}

	// Check loop structure.
	if reLoop.Match(t) {
		if m := reLoopRange.FindSubmatch(t); m != nil {
			// Range loop found.
			root.typ = TypeLoopRange
			if bytes.Contains(m[1], comma) {
				kv := bytes.Split(m[1], comma)
				root.loopKey = bytealg.Trim(kv[0], space)
				if bytes.Equal(root.loopKey, uscore) {
					root.loopKey = nil
				}
				root.loopVal = bytealg.Trim(kv[1], space)
			} else {
				root.loopKey = bytealg.Trim(m[1], space)
			}
			root.loopSrc = m[2]
			if len(m) > 2 {
				root.loopSep = m[3]
			}
		} else if m := reLoopCount.FindSubmatch(t); m != nil {
			// Counter loop found.
			root.typ = TypeLoopCount
			root.loopCnt = m[1]
			root.loopCntInit = m[2]
			root.loopCntStatic = isStatic(m[2])
			root.loopCondOp = p.parseOp(m[3])
			root.loopLim = m[4]
			root.loopLimStatic = isStatic(m[4])
			root.loopCntOp = p.parseOp(m[5])
			if len(m) > 5 {
				root.loopSep = m[6]
			}
		} else {
			return nodes, 0, up, fmt.Errorf("couldn't parse loop control structure '%s' at offset %d", t, pos)
		}

		// Create new target, increase loop counter and dive deeper.
		target := newTarget(p)
		p.cl++

		root.child = make([]Node, 0)
		root.child, offset, err = p.parseTpl(root.child, pos+len(ctl), target)

		nodes = addNode(nodes, *root)
		return nodes, offset, up, err
	}
	// Check loop end.
	if bytes.Equal(t, loopEnd) {
		// End of loop caught. Decrease the counter and exit.
		p.cl--
		offset = pos + len(ctl)
		up = true
		return nodes, offset, up, err
	}
	// Check loop lazy break (including lazybreak N).
	if m := reLoopLBrk.FindSubmatch(t); m != nil {
		root.typ = TypeLBreak
		if i, _ := strconv.ParseInt(byteconv.B2S(m[1]), 10, 64); i > 0 {
			root.loopBrkD = int(i)
		}
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	} else if bytes.Equal(t, loopLBrk) {
		root.typ = TypeLBreak
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check loop break (including break N).
	if m := reLoopBrk.FindSubmatch(t); m != nil {
		root.typ = TypeBreak
		if i, _ := strconv.ParseInt(byteconv.B2S(m[1]), 10, 64); i > 0 {
			root.loopBrkD = int(i)
		}
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	} else if bytes.Equal(t, loopBrk) {
		root.typ = TypeBreak
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check loop continue.
	if bytes.Equal(t, loopCont) {
		root.typ = TypeContinue
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check switch structure.
	if m := reSwitch.FindSubmatch(t); m != nil {
		// Create new target, increase switch counter and dive deeper.
		target := newTarget(p)
		p.cs++

		root.typ = TypeSwitch
		if len(m) > 0 {
			root.switchArg = m[1]
		}
		root.child = make([]Node, 0)
		root.child, offset, err = p.parseTpl(root.child, pos+len(ctl), target)
		root.child = rollupSwitchNodes(root.child)

		nodes = addNode(nodes, *root)
		return nodes, offset, up, err
	}
	// Check switch's case with condition helper.
	if m := reSwitchCaseHelper.FindSubmatch(t); m != nil {
		root.typ = TypeCase
		root.caseHlp = m[1]
		root.caseHlpArg = p.extractArgs(m[2])
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check switch's case with simple condition.
	if reSwitchCase.Match(t) {
		root.typ = TypeCase
		root.caseL, root.caseR, root.caseStaticL, root.caseStaticR, root.caseOp = p.parseCaseExpr(t)
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check switch's default.
	if bytes.Equal(t, swDefault) {
		root.typ = TypeDefault
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check switch end.
	if bytes.Equal(t, swEnd) {
		// End of switch caught. Decrease the counter and exit.
		p.cs--
		offset = pos + len(ctl)
		up = true
		return nodes, offset, up, err
	}

	// Check tpl interrupt.
	if bytes.Equal(t, ctlExit) {
		root.typ = TypeExit
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check control symbols.
	switch {
	case bytes.Equal(t, symEndl) || bytes.Equal(t, symNl) || bytes.Equal(t, symN) || bytes.Equal(t, symLf):
		root.typ = TypeRaw
		root.raw = nl
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	case bytes.Equal(t, symCr) || bytes.Equal(t, symR):
		root.typ = TypeRaw
		root.raw = cr
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	case bytes.Equal(t, symCrLn) || bytes.Equal(t, symRN):
		root.typ = TypeRaw
		root.raw = crlf
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	case bytes.Equal(t, symTab) || bytes.Equal(t, symT):
		root.typ = TypeRaw
		root.raw = tab
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check json quote.
	if bytes.Equal(t, jq) {
		root.typ = TypeJsonQ
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	if bytes.Equal(t, jqEnd) {
		root.typ = TypeEndJsonQ
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check HTML escape.
	if bytes.Equal(t, he) {
		root.typ = TypeHtmlE
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	if bytes.Equal(t, heEnd) {
		root.typ = TypeEndHtmlE
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check URL encode.
	if bytes.Equal(t, ue) {
		root.typ = TypeUrlEnc
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	if bytes.Equal(t, ueEnd) {
		root.typ = TypeEndUrlEnc
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check include.
	if m := reInc.FindSubmatch(t); m != nil {
		root.typ = TypeInclude
		root.tpl = bytes.Split(m[1], space)
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	return nodes, 0, up, fmt.Errorf("unknown control structure '%s' at offset %d", t, pos)
}

// Parse condition to left/right parts and condition operator.
func (p *Parser) parseCondExpr(re *regexp.Regexp, expr []byte) (l, r []byte, sl, sr bool, op Op) {
	if m := re.FindSubmatch(expr); m != nil {
		l = bytealg.Trim(m[1], space)
		if len(l) > 0 && l[0] == '!' {
			l = l[1:]
			r = bTrue
			sl = false
			sr = true
			op = OpNq
		} else {
			r = bytealg.Trim(m[3], space)
			sl = isStatic(l)
			sr = isStatic(r)
			op = p.parseOp(m[2])
		}
		if len(l) > 0 {
			l = bytealg.Trim(l, quotes)
		}
		if len(r) > 0 {
			r = bytealg.Trim(r, quotes)
		}
	}
	return
}

// Parse case condition similar to condition parsing.
func (p *Parser) parseCaseExpr(expr []byte) (l, r []byte, sl, sr bool, op Op) {
	if m := reSwitchCase.FindSubmatch(expr); m != nil {
		l = bytealg.Trim(m[1], space)
		sl = isStatic(l)
		if len(m) > 1 {
			op = p.parseOp(m[2])
			r = bytealg.Trim(m[3], space)
			sr = isStatic(r)
		}
	}
	return
}

// Convert operation from string to Op type.
func (p *Parser) parseOp(src []byte) Op {
	var op Op
	switch {
	case bytes.Equal(src, opEq):
		op = OpEq
	case bytes.Equal(src, opNq):
		op = OpNq
	case bytes.Equal(src, opGt):
		op = OpGt
	case bytes.Equal(src, opGtq):
		op = OpGtq
	case bytes.Equal(src, opLt):
		op = OpLt
	case bytes.Equal(src, opLtq):
		op = OpLtq
	case bytes.Equal(src, opInc):
		op = OpInc
	case bytes.Equal(src, opDec):
		op = OpDec
	default:
		op = OpUnk
	}
	return op
}

// Split print structure to value and mods list.
func (p *Parser) extractMods(t, outm []byte) ([]byte, []mod) {
	hasVline := bytes.Contains(t, vline)
	modNoVar := reModNoVar.Match(t) && !hasVline
	if hasVline || modNoVar || len(outm) > 0 {
		// First try to extract suffix mods, like ...|default(0).
		mods := make([]mod, 0)
		chunks := bytes.Split(t, vline)
		var idx = 1
		if modNoVar {
			idx = 0
		}
		for i := idx; i < len(chunks); i++ {
			if m := reMod.FindSubmatch(chunks[i]); m != nil {
				fn := GetModFn(byteconv.B2S(m[1]))
				if fn == nil {
					continue
				}
				args := p.extractArgs(m[2])
				mods = append(mods, mod{
					id:  m[1],
					fn:  fn,
					arg: args,
				})
			}
		}

		// Second check prefix mods, like {%q= ... %}, {%u= ... %}, ...
		if len(outm) > 0 {
			getc := func(p []byte, c byte, off int) (n int) {
				for i := off; i < len(p); i++ {
					if p[i] != c {
						break
					}
					n++
				}
				return
			}
			for off := 0; off < len(outm); {
				if outm[off] == 'j' {
					// - {%j= ... %} - JSON escape.
					fn := GetModFn("jsonEscape")
					c := getc(outm, 'j', off)
					off += c
					a := arg{val: []byte(strconv.Itoa(c)), static: true}
					mods = append(mods, mod{
						id:  idJ,
						fn:  fn,
						arg: []*arg{&a},
					})
				} else if outm[off] == 'q' {
					// - {%q= ... %} - JSON quote.
					fn := GetModFn("jsonQuote")
					c := getc(outm, 'q', off)
					off += c
					a := arg{val: []byte(strconv.Itoa(c)), static: true}
					mods = append(mods, mod{
						id:  idQ,
						fn:  fn,
						arg: []*arg{&a},
					})
				} else if outm[off] == 'h' {
					// - {%h= ... %} - HTML escape.
					fn := GetModFn("htmlEscape")
					c := getc(outm, 'h', off)
					off += c
					a := arg{val: []byte(strconv.Itoa(c)), static: true}
					mods = append(mods, mod{
						id:  idH,
						fn:  fn,
						arg: []*arg{&a},
					})
				} else if outm[off] == 'l' {
					// - {%l= ... %} - link escape.
					fn := GetModFn("linkEscape")
					c := getc(outm, 'l', off)
					off += c
					a := arg{val: []byte(strconv.Itoa(c)), static: true}
					mods = append(mods, mod{
						id:  idL,
						fn:  fn,
						arg: []*arg{&a},
					})
				} else if outm[off] == 'u' {
					// - {%u= ... %} - URL encode.
					fn := GetModFn("urlEncode")
					c := getc(outm, 'u', off)
					off += c
					a := arg{val: []byte(strconv.Itoa(c)), static: true}
					mods = append(mods, mod{
						id:  idU,
						fn:  fn,
						arg: []*arg{&a},
					})
				} else if outm[off] == 'a' {
					// - {%a= ... %} - attribute escape.
					fn := GetModFn("attrEscape")
					c := getc(outm, 'a', off)
					off += c
					a := arg{val: []byte(strconv.Itoa(c)), static: true}
					mods = append(mods, mod{
						id:  idA,
						fn:  fn,
						arg: []*arg{&a},
					})
				} else if outm[off] == 'c' {
					// - {%c= ... %} - attribute escape.
					fn := GetModFn("cssEscape")
					c := getc(outm, 'c', off)
					off += c
					a := arg{val: []byte(strconv.Itoa(c)), static: true}
					mods = append(mods, mod{
						id:  idC,
						fn:  fn,
						arg: []*arg{&a},
					})
				} else if outm[off] == 'J' {
					// - {%J= ... %} - attribute escape.
					fn := GetModFn("jsEscape")
					c := getc(outm, 'J', off)
					off += c
					a := arg{val: []byte(strconv.Itoa(c)), static: true}
					mods = append(mods, mod{
						id:  idJS,
						fn:  fn,
						arg: []*arg{&a},
					})
				} else if m := reModPfxF.FindSubmatch(outm[off:]); m != nil {
					switch m[1][0] {
					case byte(outmf):
						// - {%f.<prec>= ... %} - Float with precision.
						fn := GetModFn("floorPrec")
						mods = append(mods, mod{
							id:  idf,
							fn:  fn,
							arg: []*arg{{nil, m[2], true, false}},
						})
					case byte(outmF):
						// - {%F.<prec>= ... %} - Ceil rounded to precision float.
						fn := GetModFn("ceilPrec")
						mods = append(mods, mod{
							id:  idF,
							fn:  fn,
							arg: []*arg{{nil, m[2], true, false}},
						})
					}
					off += len(m[2]) + 2
				} else {
					// Unknown print modifier. Ignore it.
					// Perhaps need report error here.
					off++
				}
			}
		}

		if modNoVar {
			return nil, mods
		}
		return chunks[0], mods
	}
	return t, nil
}

// Get list of arguments of modifier or helper, ex:
// {%= variable|mod(arg0, ..., argN) %}
//
//	^             ^
//
// {% if condHelper(arg0, ..., argN) %}...{% endif %}
//
//	^             ^
func (p *Parser) extractArgs(raw []byte) []*arg {
	r := make([]*arg, 0)
	if len(raw) == 0 {
		return r
	}
	var (
		off, pos int
		nested   bool
	)
	for {
		if pos = bytes.IndexByte(raw[off:], ','); pos == -1 {
			pos = len(raw) - off
		}
		a := raw[off : off+pos]
		if a = bytealg.Trim(a, space); len(a) > 0 {
			if a[0] == '{' {
				a = a[1:]
				nested = true
			}
			if nested {
				kv := bytes.Split(a, colon)
				if len(kv) == 2 {
					kv[0] = bytealg.Trim(kv[0], space)
					kv[1] = bytealg.Trim(kv[1], spaceCBE)
					r = append(r, &arg{
						name:   bytealg.Trim(kv[0], quotes),
						val:    bytealg.Trim(kv[1], quotes),
						static: isStatic(kv[1]),
					})
				}
			} else {
				a = bytealg.Trim(a, space)
				arg := arg{
					val:    bytealg.Trim(a, quotes),
					static: isStatic(a),
				}
				if bytes.Equal(arg.val, ddquote) {
					arg.val = arg.val[:0]
				}
				arg.global = GetGlobal(byteconv.B2S(arg.val)) != nil
				r = append(r, &arg)
			}
			if a[len(a)-1] == '}' {
				nested = false
			}
		}

		if off+pos >= len(raw) {
			break
		}
		off += pos + 1
	}
	return r
}

// Create new target based on current parser state.
func newTarget(p *Parser) *target {
	return &target{
		targetCond:   p.cc,
		targetLoop:   p.cl,
		targetSwitch: p.cs,
	}
}

// Check if parser reached the target.
func (t *target) reached(p *Parser) bool {
	return (*t)[targetCond] == p.cc &&
		(*t)[targetLoop] == p.cl &&
		(*t)[targetSwitch] == p.cs
}

// Check if target is a root.
func (t *target) eqZero() bool {
	return (*t)[targetCond] == 0 &&
		(*t)[targetLoop] == 0 &&
		(*t)[targetSwitch] == 0
}
