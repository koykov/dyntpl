
{% if request.Secure == 1 %}
	{% if user.AllowBuy %}
		{%= user.Name %}, you may buy our items safely.
	{% else %}
		{%= user.Name %}, you should confirm your account first.
	{% endif %}
{% endif %}