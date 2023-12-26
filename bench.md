# Benchmarks

```
BenchmarkCtx
BenchmarkCtx/get
BenchmarkCtx/get-8        	 9033891	       130.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkCtx/getWithPool
BenchmarkCtx/getWithPool-8         	 6750801	       174.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl
BenchmarkTpl/condition
BenchmarkTpl/condition-8           	 1598194	       754.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/conditionHelper
BenchmarkTpl/conditionHelper-8     	 3620493	       316.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/conditionLC
BenchmarkTpl/conditionLC-8         	 2701341	       432.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/conditionNoStatic
BenchmarkTpl/conditionNoStatic-8   	 1265862	       975.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/conditionOK
BenchmarkTpl/conditionOK-8         	  330678	      3477 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/conditionStr
BenchmarkTpl/conditionStr-8        	 4167588	       275.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/counter0
BenchmarkTpl/counter0-8            	  680997	      1663 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/counter1
BenchmarkTpl/counter1-8            	  605508	      2106 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/ctxOK
BenchmarkTpl/ctxOK-8               	 1535572	       878.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/exit
BenchmarkTpl/exit-8                	 3193695	       378.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/includeHost
BenchmarkTpl/includeHost-8         	 2308333	       483.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/includeHostJS
BenchmarkTpl/includeHostJS-8       	 1468213	       829.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCount
BenchmarkTpl/loopCount-8           	  221358	      5101 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountBreak
BenchmarkTpl/loopCountBreak-8      	  211104	      5365 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountBreakN
BenchmarkTpl/loopCountBreakN-8     	  199575	      5664 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountContinue
BenchmarkTpl/loopCountContinue-8   	  160099	      7154 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountCtx
BenchmarkTpl/loopCountCtx-8        	  194980	      6452 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountLazybreak
BenchmarkTpl/loopCountLazybreak-8  	  252459	      4593 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountLazybreakN
BenchmarkTpl/loopCountLazybreakN-8 	  196512	      6022 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountStatic
BenchmarkTpl/loopCountStatic-8     	  199984	      5977 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopRange
BenchmarkTpl/loopRange-8           	  251602	      4031 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopRangeLazybreakN
BenchmarkTpl/loopRangeLazybreakN-8 	  218239	      5218 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/raw
BenchmarkTpl/raw-8                 	 7077624	       165.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/simple
BenchmarkTpl/simple-8              	  851829	      1295 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/switch
BenchmarkTpl/switch-8              	 1448582	       805.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/switchNoCondition
BenchmarkTpl/switchNoCondition-8   	 1918423	       627.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/field404
BenchmarkTpl/field404-8            	 2101747	       575.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/strAnyMap
BenchmarkTpl/strAnyMap-8           	 1248040	      1018 ns/op	       0 B/op	       0 allocs/op
BenchmarkInternalPool
BenchmarkInternalPool/ipoolUsePool
BenchmarkInternalPool/ipoolUsePool-8         	 1772076	       615.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod
BenchmarkMod/modDefault
BenchmarkMod/modDefault-8                    	 2776508	       418.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modDefaultStatic
BenchmarkMod/modDefaultStatic-8              	 1943688	       733.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modDefault1
BenchmarkMod/modDefault1-8                   	 2922708	       383.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modJSONEscape
BenchmarkMod/modJSONEscape-8                 	 2393976	       513.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modJSONEscapeShort
BenchmarkMod/modJSONEscapeShort-8            	 2300906	       536.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modJSONEscapeDbl
BenchmarkMod/modJSONEscapeDbl-8              	 1000000	      1069 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modJSONQuoteShort
BenchmarkMod/modJSONQuoteShort-8             	 2291535	       530.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modHtmlEscape
BenchmarkMod/modHtmlEscape-8                 	  911385	      1496 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modHtmlEscapeShort
BenchmarkMod/modHtmlEscapeShort-8            	  594717	      1714 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modAttrEscape
BenchmarkMod/modAttrEscape-8                 	  392866	      2908 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modAttrEscapeMB
BenchmarkMod/modAttrEscapeMB-8               	 1344667	       887.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modCSSEscape
BenchmarkMod/modCSSEscape-8                  	 1472287	       800.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modJSEscape
BenchmarkMod/modJSEscape-8                   	 1368214	       988.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modJSEscapeMB
BenchmarkMod/modJSEscapeMB-8                 	 1000000	      1074 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modLinkEscape
BenchmarkMod/modLinkEscape-8                 	 2131496	       554.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modURLEncode
BenchmarkMod/modURLEncode-8                  	 1954170	       618.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modURLEncode2
BenchmarkMod/modURLEncode2-8                 	 1576759	       769.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modURLEncode3
BenchmarkMod/modURLEncode3-8                 	 1228908	       966.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modIfThen
BenchmarkMod/modIfThen-8                     	 4199079	       281.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modIfThenElse
BenchmarkMod/modIfThenElse-8                 	 2519228	       483.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modRound
BenchmarkMod/modRound-8                      	  456318	      2636 ns/op	       0 B/op	       0 allocs/op
```
