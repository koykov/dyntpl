{
	"permission": "{% switch %}
	{% case user.Status <= 10 %}
		anonymous
	{% case user.Status <= 45 %}
		logged in
	{% case user.Status >= 60 %}
		privileged
	{% default %}
		unknown
{% endswitch %}"
}