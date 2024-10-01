<select name="type">
	{% for k, v := range user.historyTags %}
    <option value="{%= k %}">{%= v %}</option>
  {% else %}
    <option>N/D</option>
	{% endfor %}
</select>

