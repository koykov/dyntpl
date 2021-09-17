<ul>{% for i:=0; i<5; i++ %}
	{% if h, ok := __testUserNextHistory999(user.Finance).(TestHistory); ok %}
		<li>{%= h.Cost %}</li>
	{% endif %}
{%endfor%}</ul>