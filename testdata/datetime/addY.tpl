{%= date|time::add("+1 Y")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("+1 year")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("-2 years")|time::date(time::StampNano) %}
