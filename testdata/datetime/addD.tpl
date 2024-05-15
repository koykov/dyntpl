{%= date|time::add("-5 d")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("+1 day")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("+10 days")|time::date(time::StampNano) %}
