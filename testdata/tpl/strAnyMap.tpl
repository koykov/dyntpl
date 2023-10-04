before{
{% ctx list = map_.x.y.z.([]string) %}
{% for _, x := range list separator | %}
  {%= x %}
{% endfor %}
}after
