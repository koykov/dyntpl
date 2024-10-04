[
  {% for i:=0; i<10; i++ sep , %}
    {"{%=i%}":{%=i%}}
    {% break if i>5 %}
  {% endfor %}
  ,
  {% for k, v := range list separator , %}
    {% lazybreak if empty(k) %}
    {"{%= k|default(v) %}":{%= v %}}
  {% endfor %}
  ,
  {% for j:=0; j<10; j++ sep , %}
    {% for k:=0; k<10; k++ sep , %}
      {% continue if k<j %}
      {"{%= j %}":{%= k %}}
    {% endfor %}
  {% endfor %}
  ,
  {% for j:=0; j<10; j++ sep , %}
    {% for k, v := range list sep , %}
      {% break 2 if k == -1 %}
      {"{%= k %}":{%= v %}}
    {% endfor %}
  {% endfor %}
  ,
  {% for j:=0; j<10; j++ sep , %}
    {% for k:=0; k<10; k++ sep , %}
      {% lazybreak 2 if k == j %}
      {"{%= j %}":{%= k %}}
    {% endfor %}
  {% endfor %}
]
