{%= date|time::add("+1 mil")|time::date(time::RFC3339Nano) %}{% endl %}
{%= date|time::add("-1 millennium")|time::date(time::RFC3339Nano) %}
