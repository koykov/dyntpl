<Tracking event="
	{% switch track.Event %}
	{% case param.Start %}
		start
	{% case param.FirstQuartile %}
		firstQuartile
	{% case param.Midpoint %}
		midpoint
	{% case param.ThirdQuartile %}
		thirdQuartile
	{% case param.Complete %}
		complete
	{% default %}
		unknown
	{% endswitch %}
">
	<![CDATA[{%= track.Value %}]]>
</Tracking>