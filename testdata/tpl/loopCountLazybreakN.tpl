<h2>History</h2>
{% for i := 0; i < 10; i++ %}
	<ul>
		{% for j := 0; j < 10; j++ %}
		<li>Amount: {%= user.Finance.History[j].Cost %}<br/>
			Description: {%= user.Finance.History[j].Comment %}<br/>
			Date: {%= user.Finance.History[j].DateUnix %}
			{% lazybreak 2 if j == 2 %}
		</li>
		{% endfor %}
	</ul>
{% endfor %}
