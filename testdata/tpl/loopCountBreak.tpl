<h2>History</h2>
<ul>
	{% for i := 0; i < 10; i++ %}
	<li>Amount: {%= user.Finance.History[i].Cost %}<br/>
		Description: {%= user.Finance.History[i].Comment %}<br/>
		Date: {%= user.Finance.History[i].DateUnix %}
	</li>
	{% break if i == 2 %}
	{% endfor %}
</ul>
