
{
	{% for i:=0; i<10; i++ %}
		bar
		{% for j:=0; i<10; i++ %}
			foo
			{% if j == 8 %}{% break 2 %}{% endif %}
			{% if j == 7 %}{% lazybreak 2 %}{% endif %}
			{%= j %}
		{% endfor %}
		{%= i %}
	{% endfor %}
}