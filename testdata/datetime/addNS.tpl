{%= date|time::add("+1 ns")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("+1 nsec")|time::date(time::StampNano) %}
