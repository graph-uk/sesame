package sesameServer

import (
	"net/http"
)

func (t *SesameServer) pageLetMeInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte(`
<html>
<body>
<form action="/" method="post">
    <div>
        <label for="name">Token:</label>
        <input type="text" id="Token" name="Token" />
		<input type="submit" value="Submit">
    </div>
</form>
</body>
</html>
		`))
	}
}
