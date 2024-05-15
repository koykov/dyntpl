{%= date|time::add("+1 w")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("+10 week")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("-7 weeks")|time::date(time::StampNano) %}
