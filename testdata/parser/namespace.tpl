{% ctx v0 = x.y.z|testns::pack(a0, a1, "foobar") %}
{%= x.y.z|testns::extract("foobar", v1)|testns::marshal() %}
{% if testns::allow(v0) %}allowed!{% endif %}
{% if v, ok := testns::filterVar(v1); ok %}filter ok{% endif %}
{% switch %}
{% case testns::firstCase(v0)%}
  first!
{% case testns::secondCase(v0) %}
  second!
{% default %}
  all!
{% endswitch %}
