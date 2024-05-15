{%= date|time::add("-1 h")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("+1 hr")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("-21 hour")|time::date(time::StampNano) %}{% endl %}
{%= date|time::add("+1 hours")|time::date(time::StampNano) %}
