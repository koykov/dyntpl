<h2>History</h2>
<ul>
{% for i := 0; i < user.Ustate; i++ %}
  <li>{%= user.Id %}</li>
{% else %}
  <li>N/D</li>
{% endfor %}
</ul>
