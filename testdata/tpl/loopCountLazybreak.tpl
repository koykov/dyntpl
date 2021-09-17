<ul>{% for i := 0; i < 10; i++ %}
	<li>
		{%= i %}: {%= user.Finance.History[i].Cost|default(0) %}
		{% if i == 2 %}{% lazybreak %}{% endif %}
	</li>
{% endfor %}</ul>