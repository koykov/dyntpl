[
	{
		{% switch %}
		{% case firstItem(item) %}
			"name": {%= item.Name %},
		{% case secondItem(item, false) %}
			{%= item.Slug pfx "slug": sfx , %}
		{% case anonItem(item, 1) %}
			"no_data": true
		{% endswitch %}
	}
]
