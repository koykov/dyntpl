{% counter c = 0 %}
[
	{% for i := 0; i < 5; i++ separator , %}
		{% counter c++ %}
		{%= c %}
	{% endfor %}
]