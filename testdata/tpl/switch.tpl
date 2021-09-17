{% ctx exactStatus = 78 %}{
	"permission": "{% switch user.Status %}
	{% case 10 %}
		anonymous
	{% case 45 %}
		logged in
	{% case exactStatus %}
		privileged
	{% default %}
		unknown
{% endswitch %}"
}