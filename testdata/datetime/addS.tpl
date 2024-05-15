{%= date|time::add("+1 s")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("-1 sec")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("+1 second")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("-5 seconds")|time::date(time::StampNano) %}
