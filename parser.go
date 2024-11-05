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

// parser object.
type parser struct {
	target
	// Keep format flag. Remove all new lines and tabulations when false.
	keepFmt bool
	// Template body to parse.
	tpl []byte
}

var (
	// Byte constants.
	empty_   []byte
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
	opEq_  = []byte("==")
	opNq_  = []byte("!=")
	opGt_  = []byte(">")
	opGtq_ = []byte(">=")
	opLt_  = []byte("<")
	opLtq_ = []byte("<=")
	opInc_ = []byte("++")
	opDec_ = []byte("--")

	// Regexp to clear template.
	reCutComments = regexp.MustCompile(`{#[^#]*#}`)
	reCutFmt      = regexp.MustCompile(`\n+\t*\s*`)

	// Regexp to parse print instructions.
	reTplPS              = regexp.MustCompile(`^([jhqluacJfF.\d]*)=\s*(.*) (?:prefix|pfx) (.*) (?:suffix|sfx) (.*)`)
	reTplP               = regexp.MustCompile(`^([jhqluacJfF.\d]*)=\s*(.*) (?:prefix|pfx) (.*)`)
	reTplS               = regexp.MustCompile(`^([jhqluacJfF.\d]*)=\s*(.*) (?:suffix|sfx) (.*)`)
	reTpl                = regexp.MustCompile(`^([jhqluacJfF.\d]*)=\s*(.*)`)
	reTplCB              = regexp.MustCompile(`^([^(\s]+)\(([^)]*)\)`)
	reTplTernary         = regexp.MustCompile(`^([jhqluacJfF.\d]*)=\s*(.*)(==|!=|>=|<=|>|<)(.*)\s*\?\s*([^:]+):(.*)`)
	reTplTernaryHelper   = regexp.MustCompile(`^([jhqluacJfF.\d]*)=\s*([^(]+)\(*([^)]*)\)\s*\?\s*([^:]+):(.*)`)
	reTplTernaryCondExpr = regexp.MustCompile(`[jhqluacJfF.\d]*=\s*(.*)(==|!=|>=|<=|>|<)([^?]+)`)
	reModPfxF            = regexp.MustCompile(`([fF]+)\.*(\d*).*`)
	reModNoVar           = regexp.MustCompile(`([^(]+)\(([^)]*)\)`)
	reMod                = regexp.MustCompile(`([^(]+)\(*([^)]*)\)*`)

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
	// Regexp to parse break/lazybreak instructions.
	reLoopBrkN  = regexp.MustCompile(`break (\d+)`)
	reLoopLBrkN = regexp.MustCompile(`lazybreak (\d+)`)
	// Regexp to parse break-if/lazybreak-if instructions.
	reLoopBrkIf   = regexp.MustCompile(`break (if .*)`)
	reLoopBrkNIf  = regexp.MustCompile(`break (\d+) (if .*)`)
	reLoopLBrkIf  = regexp.MustCompile(`lazybreak (if .*)`)
	reLoopLBrkNIf = regexp.MustCompile(`lazybreak (\d+) (if .*)`)
	// Regexp to parse continue if-instructions.
	reLoopContIf = regexp.MustCompile(`continue (if .*)`)

	// Regexp to parse switch instruction.
	reSwitch           = regexp.MustCompile(`^switch\s*(.*)`)
	reSwitchCase       = regexp.MustCompile(`case ([^<=>!]+)([<=>!]{2})*(.*)`)
	reSwitchCaseHelper = regexp.MustCompile(`case ([^(]+)\(*([^)]*)\)`)

	// Regexp to parse include instruction.
	reInc = regexp.MustCompile(`(?:include|\.) (.*)`)

	crc64Tab = crc64.MakeTable(crc64.ISO)

	// List of lazybreak/break/continue checkers.
	xif = []struct {
		re *regexp.Regexp
		t  rtype
		i  int
		n  bool
	}{
		{reLoopLBrkNIf, typeLBreak, 2, true},   // eg: {% lazybreak 2 if v > 5 %}
		{reLoopLBrkIf, typeLBreak, 1, false},   // eg: {% lazybreak if len(x) == 4 %}
		{reLoopBrkNIf, typeBreak, 2, true},     // eg: {% break 4 if x!=3.14 %}
		{reLoopBrkIf, typeBreak, 1, false},     // eg: {% break if conditionHelperFn(x, y, z, true) %}
		{reLoopContIf, typeContinue, 1, false}, // eg: {% continue if x1 != x2 %}
	}

	// Suppress go vet warning.
	_ = ParseFile
)

// Parse initializes parser and parse the template body.
func Parse(tpl []byte, keepFmt bool) (tree *Tree, err error) {
	p := &parser{
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
	t := p.targetSnapshot()
	tree.nodes, _, err = p.parseTpl(tree.nodes, 0, t)
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
func (p *parser) cutComments() {
	p.tpl = reCutComments.ReplaceAll(p.tpl, empty_)
}

// Remove template formatting if needed.
func (p *parser) cutFmt() {
	if p.keepFmt {
		return
	}
	p.tpl = reCutFmt.ReplaceAll(p.tpl, empty_)
	p.tpl = bytealg.Trim(p.tpl, noFmt)
}

// Initial parsing method.
func (p *parser) parseTpl(nodes []node, offset int, t *target) ([]node, int, error) {
	var (
		up  bool
		err error
	)

	// Walk over template body and find control structures.
	o, i := offset, offset
	inCtl := false
	for !t.reached(p) || t.eqZero() {
		i = bytealg.IndexAt(p.tpl, ctlOpen, i)
		if i < 0 {
			if inCtl {
				err = ErrUnexpectedEOF
				break
			}
			nodes = addRaw(nodes, p.tpl[o:])
			o = len(p.tpl)
			break
		}
		if inCtl {
			// We are inside control structure.
			e := bytealg.IndexAt(p.tpl, ctlClose, i)
			if e < 0 {
				err = ErrUnexpectedEOF
				break
			}
			e += 2
			n := node{}
			nodes, e, up, err = p.processCtl(nodes, &n, p.tpl[o:e], o)
			if err != nil {
				break
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
	if !t.reached(p) {
		err = ErrUnbalancedCtl
	}
	return nodes, o, err
}

// General parsing method.
func (p *parser) processCtl(nodes []node, root *node, ctl []byte, pos int) ([]node, int, bool, error) {
	var (
		offset int
		up     bool
		err    error
	)

	up = false
	ct := bytealg.Trim(ctl, ctlTrim)
	// Check tpl (print) structure.
	if reTplPS.Match(ct) || reTplP.Match(ct) || reTplS.Match(ct) || reTpl.Match(ct) || reTplCB.Match(ct) || reTplTernary.Match(ct) || reTplTernaryHelper.Match(ct) {
		// Sequentially check print structure from the complex to the simplest.
		root.typ = typeTpl
		var m [][]byte
		if m = reTplTernary.FindSubmatch(ct); m != nil {
			root.typ = typeCond
			root.condL, root.condR, root.condStaticL, root.condStaticR, root.condOp = p.parseCondExpr(reTplTernaryCondExpr, ct)

			raw, mod_, noesc := p.extractMods(bytealg.Trim(m[5], space), m[1])
			nodeTrue := node{typ: typeCondTrue, child: []node{{typ: typeTpl, raw: raw, mod: mod_, noesc: noesc}}}
			root.child = append(root.child, nodeTrue)

			raw, mod_, noesc = p.extractMods(bytealg.Trim(m[6], space), m[1])
			nodeFalse := node{typ: typeCondFalse, child: []node{{typ: typeTpl, raw: raw, mod: mod_, noesc: noesc}}}
			root.child = append(root.child, nodeFalse)
		} else if m = reTplTernaryHelper.FindSubmatch(ct); m != nil {
			root.typ = typeCond
			// todo implement me
		} else if m = reTplPS.FindSubmatch(ct); m != nil {
			// Tpl with prefix and suffix found.
			root.raw, root.mod, root.noesc = p.extractMods(m[2], m[1])
			root.prefix = m[3]
			root.suffix = m[4]
		} else if m = reTplP.FindSubmatch(ct); m != nil {
			// Tpl with prefix found.
			root.raw, root.mod, root.noesc = p.extractMods(m[2], m[1])
			root.prefix = m[3]
		} else if m = reTplS.FindSubmatch(ct); m != nil {
			// Tpl with suffix found.
			root.raw, root.mod, root.noesc = p.extractMods(m[2], m[1])
			root.suffix = m[3]
		} else if m = reTpl.FindSubmatch(ct); m != nil {
			// Simple tpl found.
			root.raw, root.mod, root.noesc = p.extractMods(bytealg.Trim(m[2], ctlTrimAll), m[1])
		} else if m = reTplCB.FindSubmatch(ct); m != nil {
			root.raw, root.mod, root.noesc = p.extractMods(bytealg.Trim(m[0], ctlTrimAll), m[1])
		} else {
			root.raw, root.mod, root.noesc = p.extractMods(bytealg.Trim(ct, ctlTrimAll), nil)
		}
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check context structure.
	if reCtx.Match(ct) {
		root.typ = typeCtx
		var (
			m           [][]byte
			forceStatic bool
		)
		m = reCtxAs.FindSubmatch(ct)
		if m == nil {
			m = reCtxDot.FindSubmatch(ct)
		}
		if m == nil {
			m = reCtxS0.FindSubmatch(ct)
			forceStatic = m != nil
		}
		if m == nil {
			m = reCtxS1.FindSubmatch(ct)
			forceStatic = m != nil
		}
		if m == nil {
			m = reCtx.FindSubmatch(ct)
		}
		root.ctxVar, root.ctxOK = m[1], m[2]
		root.ctxSrc, root.mod, root.noesc = p.extractMods(m[3], nil)
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
	if reCntr.Match(ct) {
		root.typ = typeCounter
		root.cntrInitF = false
		if m := reCntrInit.FindSubmatch(ct); m != nil {
			root.cntrVar = m[1]
			root.cntrInitF = true
			i, err := strconv.Atoi(string(m[2]))
			if err != nil {
				return nodes, offset, up, err
			}
			root.cntrInit = i
		} else if m := reCntrOp0.FindSubmatch(ct); m != nil {
			root.cntrVar = m[1]
			if bytes.Equal(m[2], opDec_) {
				root.cntrOp = opDec
			} else {
				root.cntrOp = opInc
			}
			root.cntrOpArg = 1
		} else if m := reCntrOp1.FindSubmatch(ct); m != nil {
			root.cntrVar = m[1]
			if m[2][0] == '-' {
				root.cntrOp = opDec
			} else {
				root.cntrOp = opInc
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
	if reCondOK.Match(ct) {
		root.typ = typeCondOK
		var m [][]byte
		m = reCondAsOK.FindSubmatch(ct)
		if m == nil {
			m = reCondDotOK.FindSubmatch(ct)
		}
		if m == nil {
			m = reCondOK.FindSubmatch(ct)
		}
		root.condOKL, root.condOKR = m[1], m[2]
		root.condHlp, root.condHlpArg = m[3], p.extractArgs(m[4])
		if len(m[5]) > 0 {
			root.condIns = m[5]
		}
		root.condL, root.condR, root.condStaticL, root.condStaticR, root.condOp = p.parseCondExpr(reCondExprOK, ct)

		t := p.targetSnapshot()
		p.cc++

		subNodes := make([]node, 0)
		subNodes, offset, err = p.parseTpl(subNodes, pos+len(ctl), t)
		split := splitNodes(subNodes)
		if len(split) > 0 {
			nodeTrue := node{typ: typeCondTrue, child: split[0]}
			root.child = append(root.child, nodeTrue)
		}
		if len(split) > 1 {
			nodeFalse := node{typ: typeCondFalse, child: split[1]}
			root.child = append(root.child, nodeFalse)
		}

		nodes = addNode(nodes, *root)
		return nodes, offset, up, err
	}

	// Check condition divider.
	if bytes.Equal(ct, condElse) {
		root.typ = typeDiv
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check condition end.
	if bytes.Equal(ct, condEnd) {
		// End of condition caught. Decrease the counter and exit.
		p.cc--
		offset = pos + len(ctl)
		up = true
		return nodes, offset, up, err
	}

	// Check loop structure.
	if reLoop.Match(ct) {
		if m := reLoopRange.FindSubmatch(ct); m != nil {
			// Range loop found.
			root.typ = typeLoopRange
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
		} else if m := reLoopCount.FindSubmatch(ct); m != nil {
			// Counter loop found.
			root.typ = typeLoopCount
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
			return nodes, 0, up, fmt.Errorf("couldn't parse loop control structure '%s' at offset %d", ct, pos)
		}

		// Create new target, increase loop counter and dive deeper.
		t := p.targetSnapshot()
		p.cl++

		var subNodes []node
		subNodes, offset, err = p.parseTpl(subNodes, pos+len(ctl), t)
		split := splitNodes(subNodes)
		if len(split) > 1 {
			nodeTrue := node{typ: typeCondTrue, child: split[0]}
			root.child = append(root.child, nodeTrue)
			nodeFalse := node{typ: typeCondFalse, child: split[1]}
			root.child = append(root.child, nodeFalse)
		} else {
			root.child = subNodes
		}

		nodes = addNode(nodes, *root)
		return nodes, offset, up, err
	}
	// Check loop end.
	if bytes.Equal(ct, loopEnd) {
		// End of loop caught. Decrease the counter and exit.
		p.cl--
		offset = pos + len(ctl)
		up = true
		return nodes, offset, up, err
	}
	// Check loop lazybreak/break/continue-if (including N).
	for i := 0; i < len(xif); i++ {
		x := &xif[i]
		if m := x.re.FindSubmatch(ct); m != nil {
			if nodes, offset, up, err = p.processCond(nodes, root, ctl, pos, m[x.i], offset, up, false); err != nil {
				return nodes, offset, up, err
			}
			ch := node{typ: x.t}
			if len(m) > 2 {
				if i, _ := strconv.ParseInt(byteconv.B2S(m[1]), 10, 64); i > 0 {
					ch.loopBrkD = int(i)
				}
			}
			root.child = append(root.child, ch)
			nodes = addNode(nodes, *root)
			offset = pos + len(ctl)
			return nodes, offset, up, err
		}
	}
	// Check loop lazy break (including lazybreak N).
	if m := reLoopLBrkN.FindSubmatch(ct); m != nil {
		root.typ = typeLBreak
		if i, _ := strconv.ParseInt(byteconv.B2S(m[1]), 10, 64); i > 0 {
			root.loopBrkD = int(i)
		}
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	} else if bytes.Equal(ct, loopLBrk) {
		root.typ = typeLBreak
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check loop break (including break N).
	if m := reLoopBrkN.FindSubmatch(ct); m != nil {
		root.typ = typeBreak
		if i, _ := strconv.ParseInt(byteconv.B2S(m[1]), 10, 64); i > 0 {
			root.loopBrkD = int(i)
		}
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	} else if bytes.Equal(ct, loopBrk) {
		root.typ = typeBreak
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check loop continue.
	if bytes.Equal(ct, loopCont) {
		root.typ = typeContinue
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check condition structure.
	if reCond.Match(ct) {
		return p.processCond(nodes, root, ctl, pos, ct, offset, up, true)
	}

	// Check switch structure.
	if m := reSwitch.FindSubmatch(ct); m != nil {
		// Create new target, increase switch counter and dive deeper.
		t := p.targetSnapshot()
		p.cs++

		root.typ = typeSwitch
		if len(m) > 0 {
			root.switchArg = m[1]
		}
		root.child = make([]node, 0)
		root.child, offset, err = p.parseTpl(root.child, pos+len(ctl), t)
		root.child = rollupSwitchNodes(root.child)

		nodes = addNode(nodes, *root)
		return nodes, offset, up, err
	}
	// Check switch's case with condition helper.
	if m := reSwitchCaseHelper.FindSubmatch(ct); m != nil {
		root.typ = typeCase
		root.caseHlp = m[1]
		root.caseHlpArg = p.extractArgs(m[2])
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check switch's case with simple condition.
	if reSwitchCase.Match(ct) {
		root.typ = typeCase
		root.caseL, root.caseR, root.caseStaticL, root.caseStaticR, root.caseOp = p.parseCaseExpr(ct)
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check switch's default.
	if bytes.Equal(ct, swDefault) {
		root.typ = typeDefault
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check switch end.
	if bytes.Equal(ct, swEnd) {
		// End of switch caught. Decrease the counter and exit.
		p.cs--
		offset = pos + len(ctl)
		up = true
		return nodes, offset, up, err
	}

	// Check tpl interrupt.
	if bytes.Equal(ct, ctlExit) {
		root.typ = typeExit
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check control symbols.
	switch {
	case bytes.Equal(ct, symEndl) || bytes.Equal(ct, symNl) || bytes.Equal(ct, symN) || bytes.Equal(ct, symLf):
		root.typ = typeRaw
		root.raw = nl
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	case bytes.Equal(ct, symCr) || bytes.Equal(ct, symR):
		root.typ = typeRaw
		root.raw = cr
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	case bytes.Equal(ct, symCrLn) || bytes.Equal(ct, symRN):
		root.typ = typeRaw
		root.raw = crlf
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	case bytes.Equal(ct, symTab) || bytes.Equal(ct, symT):
		root.typ = typeRaw
		root.raw = tab
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check json quote.
	if bytes.Equal(ct, jq) {
		root.typ = typeJsonQ
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	if bytes.Equal(ct, jqEnd) {
		root.typ = typeEndJsonQ
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check HTML escape.
	if bytes.Equal(ct, he) {
		root.typ = typeHtmlE
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	if bytes.Equal(ct, heEnd) {
		root.typ = typeEndHtmlE
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check URL encode.
	if bytes.Equal(ct, ue) {
		root.typ = typeUrlEnc
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	if bytes.Equal(ct, ueEnd) {
		root.typ = typeEndUrlEnc
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check include.
	if m := reInc.FindSubmatch(ct); m != nil {
		root.typ = typeInclude
		root.tpl = bytes.Split(m[1], space)
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	return nodes, 0, up, fmt.Errorf("unknown control structure '%s' at offset %d", ct, pos)
}

func (p *parser) processCond(nodes []node, root *node, ctl []byte, pos int, ct []byte, offset int, up, dive bool) ([]node, int, bool, error) {
	var (
		subNodes []node
		split    [][]node
		err      error
	)
	// Check complexity of the condition first.
	if reCondComplex.Match(ct) {
		// Check if condition may be handled by the condition helper.
		if m := reCondHelper.FindSubmatch(ct); m != nil {
			root.typ = typeCond
			root.condHlp = m[1]
			root.condHlpArg = p.extractArgs(m[2])
			root.condL, root.condR, root.condStaticL, root.condStaticR, root.condOp = p.parseCondExpr(reCondExpr, ct)
			switch {
			case bytes.Equal(root.condHlp, condLen):
				root.condLC = lcLen
			case bytes.Equal(root.condHlp, condCap):
				root.condLC = lcCap
			}

			if dive {
				t := p.targetSnapshot()
				p.cc++
				subNodes, offset, err = p.parseTpl(subNodes, pos+len(ctl), t)
				split = splitNodes(subNodes)

				if len(split) > 0 {
					nodeTrue := node{typ: typeCondTrue, child: split[0]}
					root.child = append(root.child, nodeTrue)
				}
				if len(split) > 1 {
					nodeFalse := node{typ: typeCondFalse, child: split[1]}
					root.child = append(root.child, nodeFalse)
				}
				nodes = addNode(nodes, *root)
			}
			return nodes, offset, up, err
		}
		return nodes, pos, up, fmt.Errorf("too complex condition '%s' at offset %d", ct, pos)
	}
	root.typ = typeCond
	root.condL, root.condR, root.condStaticL, root.condStaticR, root.condOp = p.parseCondExpr(reCondExpr, ct)

	if dive {
		// Create new target, increase condition counter and dive deeper.
		t := p.targetSnapshot()
		p.cc++

		subNodes = make([]node, 0)
		subNodes, offset, err = p.parseTpl(subNodes, pos+len(ctl), t)
		split = splitNodes(subNodes)

		if len(split) > 0 {
			nodeTrue := node{typ: typeCondTrue, child: split[0]}
			root.child = append(root.child, nodeTrue)
		}
		if len(split) > 1 {
			nodeFalse := node{typ: typeCondFalse, child: split[1]}
			root.child = append(root.child, nodeFalse)
		}
		nodes = addNode(nodes, *root)
	}
	return nodes, offset, up, err
}

// Parse condition to left/right parts and condition operator.
func (p *parser) parseCondExpr(re *regexp.Regexp, expr []byte) (l, r []byte, sl, sr bool, op op) {
	if m := re.FindSubmatch(expr); m != nil {
		l = bytealg.Trim(m[1], space)
		if len(l) > 0 && l[0] == '!' {
			l = l[1:]
			r = bTrue
			sl = false
			sr = true
			op = opNq
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
func (p *parser) parseCaseExpr(expr []byte) (l, r []byte, sl, sr bool, op op) {
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

// Convert operation from string to op type.
func (p *parser) parseOp(src []byte) op {
	var op_ op
	switch {
	case bytes.Equal(src, opEq_):
		op_ = opEq
	case bytes.Equal(src, opNq_):
		op_ = opNq
	case bytes.Equal(src, opGt_):
		op_ = opGt
	case bytes.Equal(src, opGtq_):
		op_ = opGtq
	case bytes.Equal(src, opLt_):
		op_ = opLt
	case bytes.Equal(src, opLtq_):
		op_ = opLtq
	case bytes.Equal(src, opInc_):
		op_ = opInc
	case bytes.Equal(src, opDec_):
		op_ = opDec
	default:
		op_ = opUnk
	}
	return op_
}

func (p *parser) targetSnapshot() *target {
	cpy := p.target
	return &cpy
}

// Split print structure to value and mods list.
func (p *parser) extractMods(t, outm []byte) ([]byte, []mod, bool) {
	var noesc bool
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
				fnName := byteconv.B2S(m[1])
				fn := GetModFn(byteconv.B2S(m[1]))
				if fn == nil {
					continue
				}
				if noesc = fnName == "raw" || fnName == "noesc"; noesc {
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
			return nil, mods, noesc
		}
		return chunks[0], mods, noesc
	}
	return t, nil, noesc
}

// Get list of arguments of modifier or helper, ex:
// {%= variable|mod(arg0, ..., argN) %}
//
//	^             ^
//
// {% if condHelper(arg0, ..., argN) %}...{% endif %}
//
//	^             ^
func (p *parser) extractArgs(raw []byte) []*arg {
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
				a := arg{
					val:    bytealg.Trim(a, quotes),
					static: isStatic(a),
				}
				if bytes.Equal(a.val, ddquote) {
					a.val = a.val[:0]
				}
				a.global = GetGlobal(byteconv.B2S(a.val)) != nil
				r = append(r, &a)
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
