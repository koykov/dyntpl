{# don't panic, this template MUST produce bad JSON #}
{"id":"foo","name":"{% jsonquote %}{%= userName|raw %}{% endjsonquote %}"}
