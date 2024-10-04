<h2>History</h2>
<ul>
	{% for i := 0; i < 10; i++ %}
	{% if i > 2 %}{% continue if i > 2 %}{% endif %}
	<li>Amount: {%= user.Finance.History[i].Cost %}<br/>
		Description: {%= user.Finance.History[i].Comment %}<br/>
		Date: {%= user.Finance.History[i].DateUnix %}
	</li>
	{% endfor %}
</ul>
