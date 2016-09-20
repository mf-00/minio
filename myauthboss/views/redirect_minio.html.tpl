
{{$loggedin := .loggedin}}
{{$minioToken := .minioToken}}

{{if and $loggedin $minioToken}}
<script type="text/javascript">
	localStorage.token = {{$minioToken}}
	window.location.replace("/")
</script>
{{else}}
<script type="text/javascript">
	window.location.replace("/auth/login")
</script>
{{end}}

<div class="row">
</div>
