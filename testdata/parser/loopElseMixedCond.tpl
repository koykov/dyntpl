<select name="fmt">
	{% for i:=0; i<4; i++ %}
    {% if i<100 %}
      <option value="{%= i %}">{%= allowFmt[i] %}</option>
    {% else %}
      <option>overflow</option>
    {% endif %}
  {% else %}
    <option>N/D</option>
	{% endfor %}
</select>
