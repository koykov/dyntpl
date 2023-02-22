# Benchmarks

```
BenchmarkTpl/condition-8   	 1352844	       860.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/conditionHelper-8         	 3497058	       341.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/conditionNoStatic-8       	 1000000	      1075 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/conditionOK-8             	  304682	      3998 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/conditionStr-8            	 3801456	       318.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/counter0-8                	  584485	      1902 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/counter1-8                	  502466	      2557 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/ctxOK-8                   	 1425633	       786.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/exit-8                    	 3650494	       325.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/includeHost-8             	 1966533	       595.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/includeHostJS-8           	 1000000	      1004 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCount-8               	  210258	      5622 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountBreak-8          	  190034	      5999 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountBreakN-8         	  184791	      6502 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountContinue-8       	  149070	      7929 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountCtx-8            	  195589	      5889 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountLazybreak-8      	  266505	      4541 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountLazybreakN-8     	  183282	      6523 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopCountStatic-8         	  211678	      5435 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopRange-8               	  270698	      4422 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/loopRangeLazybreakN-8     	  218054	      5484 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/raw-8                     	 8120418	       146.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/simple-8                  	  992167	      1188 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/switch-8                  	 1322035	       910.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkTpl/switchNoCondition-8       	 1732689	       694.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkI18n/i18n-8                   	 1391349	       852.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkI18n/i18nPlural-8             	 2007732	       601.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkI18n/i18nPluralExt-8          	 1922486	       648.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkI18n/i18nSetLocale-8          	 1309168	       919.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod
BenchmarkMod/modDefault-8              	 2511078	       471.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modDefaultStatic-8        	 1810006	       655.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modDefault1-8             	 3282094	       363.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modJSONEscape-8           	 2121403	       562.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modJSONEscapeShort-8      	 1896626	       618.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modJSONEscapeDbl-8        	  933484	      1266 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modJSONQuoteShort-8       	 1960867	       615.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modHtmlEscape-8           	  762778	      1578 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modHtmlEscapeShort-8      	  667920	      1655 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modAttrEscape-8           	  347557	      3448 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modAttrEscapeMB-8         	 1000000	      1063 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modCSSEscape-8            	 1230736	       957.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modJSEscape-8             	 1260008	       958.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modJSEscapeMB-8           	  981488	      1098 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modLinkEscape-8           	 1769179	       655.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modURLEncode-8            	 1747382	       694.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modURLEncode2-8           	 1368250	       876.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modURLEncode3-8           	 1000000	      1109 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modIfThen-8               	 3956162	       299.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modIfThenElse-8           	 2253478	       525.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkMod/modRound-8                	  400100	      2980 ns/op	       0 B/op	       0 allocs/op
```
