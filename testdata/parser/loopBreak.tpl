
{
	{% for i:=0; i<10; i++ %}
		foo
		{% if i == 8 %}{% break %}{% endif %}
		{% if i == 7 %}{% lazybreak %}{% endif %}
		{%= i %}
	{% endfor %}
}