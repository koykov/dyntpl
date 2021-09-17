<h2>Status</h2><p>
{% if user.Status >= 60 %}Privileged user, your balance: {%= user.Finance.Balance %}.
{% else %}You don't have enough privileges.{% endif %}</p>