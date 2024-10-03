[
  {% for i:=0; i<10; i++ sep , %}
    {"{%=i%}":{%=i%}}
    {% break if i>5 %}
  {% endif %}
  ,
  {% for k, v := range list separator , %}
    {% lazybreak if empty(k) %}
    {"{%= k|default(v) %}":{%= v %}}
  {% endif %}
  ,
  {% for j:=0; j<10; j++ sep , %}
    {% for k:=0; k<10; k++ sep , %}
      {% continue if k<j %}
      {"{%= j %}":{%= k %}}
    {% endif %}
  {% endif %}
  ,
  {% for j:=0; j<10; j++ sep , %}
    {% k, v := range list sep , %}
      {% break 2 if k == -1 %}
      {"{%= k %}":{%= v %}}
    {% endif %}
  {% endif %}
  ,
  {% for j:=0; j<10; j++ sep , %}
    {% for k:=0; k<10; k++ sep , %}
      {% lazybreak 2 if k == j %}
      {"{%= j %}":{%= k %}}
    {% endif %}
  {% endif %}
]
