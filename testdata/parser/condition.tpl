
<h1>Profile</h1>
{% if user.Id == 0 %}
You should <a href="{%= loginUrl %}">log in</a>.
{% else %}
Welcome, {%= user.Name %}.
{% endif %}
Copyright {%= brand %}
