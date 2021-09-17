<h2>Status</h2><p>
{% ctx permissionLimit = 60 %}
{% if user.Status >= permissionLimit %}Privileged user, your balance: {%= user.Finance.Balance %}.
{% else %}You don't have enough privileges.{% endif %}</p>