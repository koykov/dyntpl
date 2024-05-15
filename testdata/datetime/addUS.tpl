{%= date|time::add("+5 us")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("+5 usec")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("+5 µs")|time::date(time::StampNano) %}
