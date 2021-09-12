
<h2>Export history</h2>
<label>Type</label>
<select name="type">
	{% for k, v := range user.historyTags %}
	<option value="{%= k %}">{%= v %}</option>
	{% endfor %}
</select>
<label>Format</label>
<select name="fmt">
	{% for i:=0; i<4; i++ %}
	<option value="{%= i %}">{%= allowFmt[i] %}</option>
	{% endfor %}
</select>