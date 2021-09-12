
[
	{
		{% switch %}
		{% case item.Index == 0 %}
			"name": {%= item.Name %},
		{% case item.Index == 1 %}
			{%= item.Slug pfx "slug": sfx , %}
		{% case item.Uid == -1 %}
			"no_data": true
		{% endswitch %}
	}
]