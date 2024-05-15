{%= date|time::add("+1 M")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("+1 month")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("-2 months")|time::date(time::StampNano) %}
