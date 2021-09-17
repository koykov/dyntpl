<h2>History</h2>
{% for i := 0; i < 10; i++ %}
	<ul>
		{% for j := 0; j < 10; j++ %}
		<li>Amount: {%= user.Finance.History[j].Cost %}<br/>
			Description: {%= user.Finance.History[j].Comment %}<br/>
			Date: {%= user.Finance.History[j].DateUnix %}
			{% if j == 2 %}{% lazybreak 2 %}{% endif %}
		</li>
		{% endfor %}
	</ul>
{% endfor %}