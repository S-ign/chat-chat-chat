$('#message_form').focus();
$(document).ready(function() {
	if (window.EventSource) {
		var source = new EventSource('/stream/1');
		source.addEventListener('message', function(e) {
			$('#messages').append(e.data + "</br>");
			$('#message_form').val('');
		}, true);
	} else {
		alert("NOT SUPPORTED");
	}
});
